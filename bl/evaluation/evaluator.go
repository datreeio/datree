package evaluation

import (
	"encoding/json"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
)

const (
	SKIP_RULE_PREFIX string = "datree.io/skip/"
)

type CLIClient interface {
	SendEvaluationResult(request *cliClient.EvaluationResultRequest) (*cliClient.SendEvaluationResultsResponse, error)
}

type Evaluator struct {
	cliClient           CLIClient
	ciContext           *ciContext.CIContext
	jsonSchemaValidator *jsonSchemaValidator.JSONSchemaValidator
}

func New(c CLIClient) *Evaluator {
	return &Evaluator{
		cliClient:           c,
		ciContext:           ciContext.Extract(),
		jsonSchemaValidator: jsonSchemaValidator.New(),
	}
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
	Token              string
	ClientId           string
	CliVersion         string
	K8sVersion         string
	PolicyName         string
	CiContext          *ciContext.CIContext
	RulesData          []cliClient.RuleData
	FilesData          []cliClient.FileData
	FailedYamlFiles    []string
	FailedK8sFiles     []string
	PolicyCheckResults FailedRulesByFiles
}

var OSInfoFn = NewOSInfo

func (e *Evaluator) SendEvaluationResult(evaluationRequestData EvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error) {
	osInfo := OSInfoFn()
	sendEvaluationResultsResponse, err := e.cliClient.SendEvaluationResult(&cliClient.EvaluationResultRequest{
		K8sVersion: evaluationRequestData.K8sVersion,
		ClientId:   evaluationRequestData.ClientId,
		Token:      evaluationRequestData.Token,
		PolicyName: evaluationRequestData.PolicyName,
		Metadata: &cliClient.Metadata{
			CliVersion:      evaluationRequestData.CliVersion,
			Os:              osInfo.OS,
			PlatformVersion: osInfo.PlatformVersion,
			KernelVersion:   osInfo.KernelVersion,
			CIContext:       evaluationRequestData.CiContext,
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

func (e *Evaluator) Evaluate(policyCheckData PolicyCheckData) (PolicyCheckResultData, error) {
	rulesCount := len(policyCheckData.Policy.Rules)

	if len(policyCheckData.FilesConfigurations) == 0 {
		return PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, FailedRulesByFiles{}, rulesCount}, nil
	}

	emptyPolicyCheckResult := PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}

	var filesData []cliClient.FileData
	for _, filesConfiguration := range policyCheckData.FilesConfigurations {
		filesData = append(filesData, cliClient.FileData{FilePath: filesConfiguration.FileName, ConfigurationsCount: len(filesConfiguration.Configurations)})
	}

	rulesData := []cliClient.RuleData{}
	for _, rule := range policyCheckData.Policy.Rules {
		rulesData = append(rulesData, cliClient.RuleData{Identifier: rule.RuleIdentifier, Name: rule.RuleName})
	}

	// map of files paths to map of rules to failed rule data
	failedRulesByFiles := make(FailedRulesByFiles)
	for _, filesConfiguration := range policyCheckData.FilesConfigurations {
		for _, configuration := range filesConfiguration.Configurations {
			// add all configurations skipped rules to the skipped rules map
			err := e.evaluateConfiguration(failedRulesByFiles, policyCheckData, filesConfiguration.FileName, configuration)
			if err != nil {
				return emptyPolicyCheckResult, err
			}
		}
	}

	formattedResults := FormattedResults{}
	formattedResults.EvaluationResults = e.formatEvaluationResults(failedRulesByFiles, len(policyCheckData.FilesConfigurations), rulesCount)

	if !policyCheckData.IsInteractiveMode {
		formattedResults.NonInteractiveEvaluationResults = e.formatNonInteractiveEvaluationResults(formattedResults.EvaluationResults, failedRulesByFiles, policyCheckData.PolicyName, rulesCount)
	}

	return PolicyCheckResultData{formattedResults, rulesData, filesData, failedRulesByFiles, rulesCount}, nil
}

func (e *Evaluator) evaluateConfiguration(failedRulesByFiles FailedRulesByFiles, policyCheckData PolicyCheckData, fileName string, configuration extractor.Configuration) error {
	skipAnnotations := extractSkipAnnotations(configuration)

	configurationName, configurationKind := extractConfigurationInfo(configuration)

	configurationJson, err := json.Marshal(configuration)
	if err != nil {
		return err
	}

	for _, rule := range policyCheckData.Policy.Rules {

		failedRule, err := e.evaluateRule(rule, configurationJson, configurationName, configurationKind, skipAnnotations)
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

func (e *Evaluator) evaluateRule(rule policy_factory.RuleWithSchema, configurationJson []byte, configurationName string, configurationKind string, skipAnnotations map[string]string) (*cliClient.FailedRule, error) {
	ruleSchemaJson, err := json.Marshal(rule.Schema)
	if err != nil {
		return nil, err
	}

	validationResult, err := e.jsonSchemaValidator.ValidateYamlSchemaNew(string(ruleSchemaJson), string(configurationJson))

	if err != nil {
		return nil, err
	}

	occurrences := len(validationResult)
	skipMessage, skipRuleExists := skipAnnotations[SKIP_RULE_PREFIX+rule.RuleIdentifier]

	if occurrences < 1 && !skipRuleExists {
		return nil, nil
	}

	configuration := cliClient.Configuration{
		Name:        configurationName,
		Kind:        configurationKind,
		Occurrences: occurrences,
		IsSkipped:   false,
		SkipMessage: "",
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

// This method creates a NonInteractiveEvaluationResults structure
// from EvaluationResults.
func (e *Evaluator) formatNonInteractiveEvaluationResults(formattedEvaluationResults *EvaluationResults, evaluationResults FailedRulesByFiles, policyName string, totalRulesInPolicy int) *NonInteractiveEvaluationResults {
	fileNameRuleMapper := formattedEvaluationResults.FileNameRuleMapper
	ruleMapper := make(map[string]string)

	for filePath := range evaluationResults {
		for ruleIdentifier := range evaluationResults[filePath] {
			ruleMapper[ruleIdentifier] = ruleIdentifier
		}
	}

	nonInteractiveEvaluationResults := NonInteractiveEvaluationResults{}

	for fileName, rules := range fileNameRuleMapper {
		formattedEvaluationResults := FormattedEvaluationResults{}
		formattedEvaluationResults.FileName = fileName

		for _, rule := range rules {
			ruleResult := RuleResult{Identifier: ruleMapper[rule.Identifier], Name: rule.Name, MessageOnFailure: rule.MessageOnFailure, OccurrencesDetails: rule.OccurrencesDetails}
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
		PolicyName:         policyName,
		TotalRulesInPolicy: totalRulesInPolicy,
		TotalRulesFailed:   formattedEvaluationResults.Summary.TotalFailedRules,
		TotalSkippedRules:  formattedEvaluationResults.Summary.TotalSkippedRules,
		TotalPassedCount:   formattedEvaluationResults.Summary.TotalPassedRules,
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
						MetadataName: configuration.Name,
						Kind:         configuration.Kind,
						Occurrences:  configuration.Occurrences,
						IsSkipped:    configuration.IsSkipped,
						SkipMessage:  configuration.SkipMessage,
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

func extractConfigurationInfo(configuration extractor.Configuration) (string, string) {
	kind := ""
	name := ""

	nonStringKind := configuration["kind"]
	if nonStringKind != nil {
		kind = nonStringKind.(string)
	}

	nonObjectMetadata := configuration["metadata"]
	if nonObjectMetadata != nil {
		nonStringName := nonObjectMetadata.(map[string]interface{})["name"]
		if nonStringName != nil {
			name = nonStringName.(string)
		}
	}

	return name, kind
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

	configurationMetadata, ok := configuration["metadata"].(map[string]interface{})
	if !ok {
		return nil
	}

	annotationsMap, ok := configurationMetadata["annotations"].(map[string]interface{})
	if !ok {
		return nil
	}

	for annotationKey, annotationValue := range annotationsMap {
		if strings.Contains(annotationKey, SKIP_RULE_PREFIX) {
			skipAnnotations[annotationKey] = annotationValue.(string)
		}
	}

	return skipAnnotations
}
