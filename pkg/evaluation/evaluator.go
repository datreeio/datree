package evaluation

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/mikefarah/yq/v4/pkg/yqlib"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

const (
	SKIP_RULE_PREFIX string = "datree.skip/"
)

type CLIClient interface {
	SendEvaluationResult(request *cliClient.EvaluationResultRequest) (*cliClient.SendEvaluationResultsResponse, error)
}

type Evaluator struct {
	cliClient           CLIClient
	ciContext           *ciContext.CIContext
	jsonSchemaValidator *jsonSchemaValidator.JSONSchemaValidator
	yqlibEvaluator      yqlib.Evaluator
}

func New(c CLIClient, ciContext *ciContext.CIContext) *Evaluator {
	return &Evaluator{
		cliClient:           c,
		ciContext:           ciContext,
		jsonSchemaValidator: jsonSchemaValidator.New(),
		yqlibEvaluator:      newYqEvaluator(),
	}
}

func newYqEvaluator() yqlib.Evaluator {
	yqEvaluator := yqlib.NewAllAtOnceEvaluator()
	logger := yqlib.GetLogger()
	backendLogger := logging.NewLogBackend(os.Stderr, "", 0)
	backendLoggerLeveled := logging.AddModuleLevel(backendLogger)
	backendLoggerLeveled.SetLevel(logging.ERROR, "")
	logger.SetBackend(backendLoggerLeveled)
	return yqEvaluator
}

type FileNameRuleMapper map[string]map[string]*Rule
type FailedRulesByFiles map[string]map[string]*cliClient.FailedRule
type EvaluationResultsSummery struct {
	TotalFailedRules  int
	TotalSkippedRules int
	TotalPassedRules  int
	FilesCount        int
	FilesPassedCount  int
}

type EvaluationResults struct {
	FileNameRuleMapper FileNameRuleMapper
	Summary            EvaluationResultsSummery
}

type FormattedResults struct {
	EvaluationResults               *EvaluationResults
	NonInteractiveEvaluationResults *NonInteractiveEvaluationResults
}

type EvaluationRequestData struct {
	Token                     string
	ClientId                  string
	CliVersion                string
	K8sVersion                string
	PolicyName                string
	CiContext                 *ciContext.CIContext
	RulesData                 []cliClient.RuleData
	FilesData                 []cliClient.FileData
	FailedYamlFiles           []string
	FailedK8sFiles            []string
	PolicyCheckResults        FailedRulesByFiles
	EvaluationDurationSeconds float64
}

var OSInfoFn = utils.NewOSInfo

func (e *Evaluator) SendEvaluationResult(evaluationRequestData EvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error) {
	osInfo := OSInfoFn()
	sendEvaluationResultsResponse, err := e.cliClient.SendEvaluationResult(&cliClient.EvaluationResultRequest{
		K8sVersion: evaluationRequestData.K8sVersion,
		ClientId:   evaluationRequestData.ClientId,
		Token:      evaluationRequestData.Token,
		PolicyName: evaluationRequestData.PolicyName,
		Metadata: &cliClient.Metadata{
			CliVersion:                evaluationRequestData.CliVersion,
			Os:                        osInfo.OS,
			PlatformVersion:           osInfo.PlatformVersion,
			KernelVersion:             osInfo.KernelVersion,
			CIContext:                 evaluationRequestData.CiContext,
			EvaluationDurationSeconds: evaluationRequestData.EvaluationDurationSeconds,
		},
		FailedYamlFiles:    evaluationRequestData.FailedYamlFiles,
		FailedK8sFiles:     evaluationRequestData.FailedK8sFiles,
		AllExecutedRules:   evaluationRequestData.RulesData,
		AllEvaluatedFiles:  evaluationRequestData.FilesData,
		PolicyCheckResults: evaluationRequestData.PolicyCheckResults,
	})

	return sendEvaluationResultsResponse, err
}

type PolicyCheckData struct {
	FilesConfigurations []*extractor.FileConfigurations
	IsInteractiveMode   bool
	PolicyName          string
	Policy              policy_factory.Policy
}

type PolicyCheckResultData struct {
	FormattedResults FormattedResults
	RulesData        []cliClient.RuleData
	FilesData        []cliClient.FileData
	RawResults       FailedRulesByFiles
	RulesCount       int
}

type EvaluateData struct {
	PolicyCheckData PolicyCheckData
	Verbose         bool
}

func (e *Evaluator) Evaluate(evaluateData EvaluateData) (PolicyCheckResultData, error) {
	rulesCount := len(evaluateData.PolicyCheckData.Policy.Rules)

	if len(evaluateData.PolicyCheckData.FilesConfigurations) == 0 {
		return PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, FailedRulesByFiles{}, rulesCount}, nil
	}

	emptyPolicyCheckResult := PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}

	var filesData []cliClient.FileData
	for _, filesConfiguration := range evaluateData.PolicyCheckData.FilesConfigurations {
		filesData = append(filesData, cliClient.FileData{FilePath: filesConfiguration.FileName, ConfigurationsCount: len(filesConfiguration.Configurations)})
	}

	rulesData := []cliClient.RuleData{}
	for _, rule := range evaluateData.PolicyCheckData.Policy.Rules {
		rulesData = append(rulesData, cliClient.RuleData{Identifier: rule.RuleIdentifier, Name: rule.RuleName})
	}

	// map of files paths to map of rules to failed rule data
	failedRulesByFiles := make(FailedRulesByFiles)
	for _, filesConfiguration := range evaluateData.PolicyCheckData.FilesConfigurations {
		for _, configuration := range filesConfiguration.Configurations {
			// add all configurations skipped rules to the skipped rules map
			err := e.evaluateConfiguration(failedRulesByFiles, evaluateData.PolicyCheckData, filesConfiguration.FileName, configuration)
			if err != nil {
				return emptyPolicyCheckResult, err
			}
		}
	}

	formattedResults := FormattedResults{}
	formattedResults.EvaluationResults = e.formatEvaluationResults(failedRulesByFiles, len(evaluateData.PolicyCheckData.FilesConfigurations), rulesCount)

	nonInteractiveEvaluationData := NonInteractiveEvaluationData{
		FormattedEvaluationResults: formattedResults.EvaluationResults,
		EvaluationResults:          failedRulesByFiles,
		PolicyName:                 evaluateData.PolicyCheckData.PolicyName,
		TotalRulesInPolicy:         rulesCount,
		Verbose:                    evaluateData.Verbose,
	}

	formattedResults.NonInteractiveEvaluationResults = e.formatNonInteractiveEvaluationResults(nonInteractiveEvaluationData)

	return PolicyCheckResultData{formattedResults, rulesData, filesData, failedRulesByFiles, rulesCount}, nil
}

func (e *Evaluator) evaluateConfiguration(failedRulesByFiles FailedRulesByFiles, policyCheckData PolicyCheckData, fileName string, configuration extractor.Configuration) error {
	skipAnnotations := extractSkipAnnotations(configuration)

	for _, rule := range policyCheckData.Policy.Rules {
		failedRule, err := e.evaluateRule(rule, configuration.Payload, configuration.MetadataName, configuration.Kind, skipAnnotations, configuration.YamlNode)
		if err != nil {
			return err
		}

		if failedRule == nil {
			continue
		}
		addFailedRule(failedRulesByFiles, fileName, rule.RuleIdentifier, failedRule)
	}

	return nil
}

func (e *Evaluator) evaluateRule(rule policy_factory.RuleWithSchema, configurationJson []byte, configurationName string, configurationKind string, skipAnnotations map[string]string, yamlNode yaml.Node) (*cliClient.FailedRule, error) {
	ruleSchemaJson, err := json.Marshal(rule.Schema)
	if err != nil {
		return nil, err
	}

	validationResult, err := e.jsonSchemaValidator.ValidateYamlSchema(string(ruleSchemaJson), string(configurationJson))

	if err != nil {
		return nil, err
	}

	occurrences := len(validationResult)
	skipMessage, skipRuleExists := skipAnnotations[SKIP_RULE_PREFIX+rule.RuleIdentifier]

	if occurrences < 1 && !skipRuleExists {
		return nil, nil
	}

	configuration := cliClient.Configuration{
		Name:             configurationName,
		Kind:             configurationKind,
		Occurrences:      occurrences,
		IsSkipped:        false,
		SkipMessage:      "",
		FailureLocations: []cliClient.FailureLocation{},
	}

	for _, detailedResult := range validationResult {
		failedErrorLine, failedErrorColumn := e.getFailedRuleLineAndColumn(detailedResult.InstanceLocation, yamlNode)

		failureLocation := cliClient.FailureLocation{
			SchemaPath:        detailedResult.InstanceLocation,
			FailedErrorLine:   failedErrorLine,
			FailedErrorColumn: failedErrorColumn,
		}

		configuration.FailureLocations = append(configuration.FailureLocations, failureLocation)
	}

	if skipRuleExists {
		configuration.IsSkipped = true
		configuration.SkipMessage = skipMessage
	}

	failedRule := &cliClient.FailedRule{
		Name:             rule.RuleName,
		DocumentationUrl: rule.DocumentationUrl,
		MessageOnFailure: rule.MessageOnFailure,
		Configurations:   []cliClient.Configuration{configuration},
	}

	return failedRule, nil
}

type NonInteractiveEvaluationData struct {
	FormattedEvaluationResults *EvaluationResults
	EvaluationResults          FailedRulesByFiles
	PolicyName                 string
	TotalRulesInPolicy         int
	Verbose                    bool
}

// This method creates a NonInteractiveEvaluationResults structure
// from EvaluationResults.
func (e *Evaluator) formatNonInteractiveEvaluationResults(nonInteractiveEvaluationData NonInteractiveEvaluationData) *NonInteractiveEvaluationResults {
	fileNameRuleMapper := nonInteractiveEvaluationData.FormattedEvaluationResults.FileNameRuleMapper
	ruleMapper := make(map[string]string)

	for filePath := range nonInteractiveEvaluationData.EvaluationResults {
		for ruleIdentifier := range nonInteractiveEvaluationData.EvaluationResults[filePath] {
			ruleMapper[ruleIdentifier] = ruleIdentifier
		}
	}

	nonInteractiveEvaluationResults := NonInteractiveEvaluationResults{}

	for fileName, rules := range fileNameRuleMapper {
		formattedEvaluationResults := FormattedEvaluationResults{}
		formattedEvaluationResults.FileName = fileName

		for _, rule := range rules {
			ruleResult := RuleResult{Identifier: ruleMapper[rule.Identifier], Name: rule.Name, MessageOnFailure: rule.MessageOnFailure, OccurrencesDetails: rule.OccurrencesDetails}
			if nonInteractiveEvaluationData.Verbose {
				ruleResult.DocumentationUrl = rule.DocumentationUrl
			}

			formattedEvaluationResults.RuleResults = append(
				formattedEvaluationResults.RuleResults,
				&ruleResult,
			)
		}
		nonInteractiveEvaluationResults.FormattedEvaluationResults = append(
			nonInteractiveEvaluationResults.FormattedEvaluationResults,
			&formattedEvaluationResults,
		)
	}
	nonInteractiveEvaluationResults.PolicySummary = &PolicySummary{
		PolicyName:         nonInteractiveEvaluationData.PolicyName,
		TotalRulesInPolicy: nonInteractiveEvaluationData.TotalRulesInPolicy,
		TotalRulesFailed:   nonInteractiveEvaluationData.FormattedEvaluationResults.Summary.TotalFailedRules,
		TotalSkippedRules:  nonInteractiveEvaluationData.FormattedEvaluationResults.Summary.TotalSkippedRules,
		TotalPassedCount:   nonInteractiveEvaluationData.FormattedEvaluationResults.Summary.TotalPassedRules,
	}

	return &nonInteractiveEvaluationResults
}

func (e *Evaluator) formatEvaluationResults(evaluationResults FailedRulesByFiles, filesCount int, rulesCount int) *EvaluationResults {
	mapper := make(map[string]map[string]*Rule)

	totalFailedCount := 0
	totalSkippedCount := 0
	failedFilesCount := len(evaluationResults)

	for filePath := range evaluationResults {
		if _, exists := mapper[filePath]; !exists {
			mapper[filePath] = make(map[string]*Rule)
		}

		for ruleIdentifier, failedRule := range evaluationResults[filePath] {
			// file and rule not already exists in mapper
			if _, exists := mapper[filePath][ruleIdentifier]; !exists {
				mapper[filePath][ruleIdentifier] = &Rule{
					Identifier:         ruleIdentifier,
					Name:               failedRule.Name,
					DocumentationUrl:   failedRule.DocumentationUrl,
					MessageOnFailure:   failedRule.MessageOnFailure,
					OccurrencesDetails: []OccurrenceDetails{},
				}
			}

			for _, configuration := range failedRule.Configurations {
				mapper[filePath][ruleIdentifier].OccurrencesDetails = append(
					mapper[filePath][ruleIdentifier].OccurrencesDetails,
					OccurrenceDetails{
						MetadataName:     configuration.Name,
						Kind:             configuration.Kind,
						Occurrences:      configuration.Occurrences,
						IsSkipped:        configuration.IsSkipped,
						SkipMessage:      configuration.SkipMessage,
						FailureLocations: configuration.FailureLocations,
					},
				)
			}
		}

		allRulesAreSkipped := true

		for _, rule := range mapper[filePath] {
			skippedOccurrences := 0
			totalOccurrences := len(rule.OccurrencesDetails)

			for _, occurrence := range rule.OccurrencesDetails {
				if occurrence.IsSkipped {
					skippedOccurrences++
				} else {
					allRulesAreSkipped = false
				}
			}

			if totalOccurrences == skippedOccurrences {
				totalSkippedCount++
			} else if skippedOccurrences >= 1 {
				totalSkippedCount++
				totalFailedCount++
			} else {
				totalFailedCount++
			}
		}

		if allRulesAreSkipped {
			failedFilesCount--
		}
	}

	results := &EvaluationResults{
		FileNameRuleMapper: mapper,
		Summary: EvaluationResultsSummery{
			TotalFailedRules:  totalFailedCount,
			TotalSkippedRules: totalSkippedCount,
			TotalPassedRules:  (rulesCount * filesCount) - (totalFailedCount + totalSkippedCount),
			FilesCount:        filesCount,
			FilesPassedCount:  filesCount - failedFilesCount,
		},
	}

	return results
}

type Result = gojsonschema.Result

func addFailedRule(currentFailedRulesByFiles FailedRulesByFiles, fileName string, ruleIdentifier string, failedRule *cliClient.FailedRule) {
	fileData, ok := currentFailedRulesByFiles[fileName]

	if !ok {
		currentFailedRulesByFiles[fileName] = map[string]*cliClient.FailedRule{ruleIdentifier: failedRule}
		return
	}

	if exitingRule, ok := fileData[ruleIdentifier]; ok {
		exitingRule.Configurations = append(exitingRule.Configurations, failedRule.Configurations...)
	} else {
		currentFailedRulesByFiles[fileName][ruleIdentifier] = failedRule
	}
}

func extractSkipAnnotations(configuration extractor.Configuration) map[string]string {
	skipAnnotations := make(map[string]string)

	for annotationKey, annotationValue := range configuration.Annotations {
		if strings.Contains(annotationKey, SKIP_RULE_PREFIX) {
			skipAnnotations[annotationKey] = annotationValue.(string)
		}
	}

	return skipAnnotations
}

func (e *Evaluator) getFailedRuleLineAndColumn(schemaPath string, yamlNode yaml.Node) (failedErrorLine int, failedErrorColumn int) {

	instanceLocationYqPath := strings.Replace(schemaPath, "/", ".", -1)
	instanceLocationYqPath = regexp.MustCompile(`\d+`).ReplaceAllString(instanceLocationYqPath, `[$0]`)

	nodeList, err := e.yqlibEvaluator.EvaluateNodes(instanceLocationYqPath, &yamlNode)
	if err != nil {
		return
	}

	candidateNode := nodeList.Back().Value.(*yqlib.CandidateNode).Node

	return candidateNode.Line, candidateNode.Column
}
