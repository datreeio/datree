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
	SKIP_ALL_KEY     string = SKIP_RULE_PREFIX + "ALL"
)

type CLIClient interface {
	SendEvaluationResult(request *cliClient.EvaluationResultRequest) (*cliClient.SendEvaluationResultsResponse, error)
}

type Evaluator struct {
	cliClient CLIClient
	ciContext *ciContext.CIContext
	jsonSchemaValidator *jsonSchemaValidator.JSONSchemaValidator
}

func New(c CLIClient) *Evaluator {
	return &Evaluator{
		cliClient: c,
		ciContext: ciContext.Extract(),
		jsonSchemaValidator: jsonSchemaValidator.New(),
	}
}

type FileNameRuleMapper map[string]map[string]*Rule
type FailedRulesByFiles map[string]map[string]*cliClient.FailedRule

type EvaluationResults struct {
	FileNameRuleMapper FileNameRuleMapper
	Summary            struct {
		TotalFailedRules int
		FilesCount       int
		TotalPassedCount int
	}
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

	// map of files paths to map of rules to failed rule data
	failedRulesByFiles := make(FailedRulesByFiles)

	emptyPolicyCheckResult := PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}

	var filesData []cliClient.FileData
	for _, filesConfiguration := range policyCheckData.FilesConfigurations {
		filesData = append(filesData, cliClient.FileData{FilePath: filesConfiguration.FileName, ConfigurationsCount: len(filesConfiguration.Configurations)})
	}

	rulesData := []cliClient.RuleData{}
	for _, rule := range policyCheckData.Policy.Rules {
		rulesData = append(rulesData, cliClient.RuleData{Identifier: rule.RuleIdentifier, Name: rule.RuleName})
	}

	for _, filesConfiguration := range policyCheckData.FilesConfigurations {
		for _, configuration := range filesConfiguration.Configurations {
			skipAnnotations := extractSkipAnnotations(configuration)
			if _, ok := skipAnnotations[SKIP_ALL_KEY]; ok {
				// addSkipRule
				continue
			}

			configurationName, configurationKind := extractConfigurationInfo(configuration)

			configurationJson, err := json.Marshal(configuration)
			if err != nil {
				return emptyPolicyCheckResult, err
			}

			for _, rule := range policyCheckData.Policy.Rules {

				failedRule, err := e.evaluateRule(rule, configurationJson, configurationName, configurationKind, skipAnnotations)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				if failedRule == nil {
					continue
				}

				addFailedRule(failedRulesByFiles, filesConfiguration.FileName, rule.RuleIdentifier, failedRule)
			}
		}

	}

	formattedResults := FormattedResults{}
	formattedResults.EvaluationResults = e.formatEvaluationResults(failedRulesByFiles, len(policyCheckData.FilesConfigurations))

	if !policyCheckData.IsInteractiveMode {
		formattedResults.NonInteractiveEvaluationResults = e.formatNonInteractiveEvaluationResults(formattedResults.EvaluationResults, failedRulesByFiles, policyCheckData.PolicyName, rulesCount)
	}

	return PolicyCheckResultData{formattedResults, rulesData, filesData, failedRulesByFiles, rulesCount}, nil
}

func (e *Evaluator) evaluateRule(rule policy_factory.RuleWithSchema, configurationJson []byte, configurationName string, configurationKind string, skipAnnotations map[string]string) (*cliClient.FailedRule, error) {
	ruleSchemaJson, err := json.Marshal(rule.Schema)
	if err != nil {
		return nil, err
	}

	validationResult, err := e.jsonSchemaValidator.ValidateYamlSchema(string(ruleSchemaJson), string(configurationJson))

	if err != nil {
		return nil, err
	}

	occurrences := countOccurrences(validationResult)

	if occurrences < 1 {
		return nil, nil
	}

	configurationData := cliClient.Configuration{
		Name:        configurationName,
		Kind:        configurationKind,
		Occurrences: occurrences,
		IsSkipped:   false,
		SkipMessage: "",
	}

	if skipMessage, ok := skipAnnotations[SKIP_RULE_PREFIX+rule.RuleIdentifier]; ok {
		configurationData.IsSkipped = true
		configurationData.SkipMessage = skipMessage
	}

	failedRule := &cliClient.FailedRule{
		Name:             rule.RuleName,
		DocumentationUrl: rule.DocumentationUrl,
		MessageOnFailure: rule.MessageOnFailure,
		Configurations:   []cliClient.Configuration{configurationData},
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
		TotalPassedCount:   formattedEvaluationResults.Summary.TotalPassedCount,
	}

	return &nonInteractiveEvaluationResults
}

func (e *Evaluator) formatEvaluationResults(evaluationResults FailedRulesByFiles, filesCount int) *EvaluationResults {
	mapper := make(map[string]map[string]*Rule)

	totalFailedCount := 0
	totalPassedCount := filesCount

	for filePath := range evaluationResults {
		if _, exists := mapper[filePath]; !exists {
			mapper[filePath] = make(map[string]*Rule)
			totalPassedCount = totalPassedCount - 1
		}

		for ruleIdentifier, failedRuleData := range evaluationResults[filePath] {
			// file and rule not already exists in mapper
			// all configuration is skipped then skip
			if _, exists := mapper[filePath][ruleIdentifier]; !exists {
				totalFailedCount++
				mapper[filePath][ruleIdentifier] = &Rule{
					Identifier:         ruleIdentifier,
					Name:               failedRuleData.Name,
					DocumentationUrl:   failedRuleData.DocumentationUrl,
					MessageOnFailure:   failedRuleData.MessageOnFailure,
					OccurrencesDetails: []OccurrenceDetails{},
				}
			}

			for _, configuration := range failedRuleData.Configurations {
				mapper[filePath][ruleIdentifier].OccurrencesDetails = append(
					mapper[filePath][ruleIdentifier].OccurrencesDetails,
					OccurrenceDetails{MetadataName: configuration.Name, Kind: configuration.Kind, Occurrences: configuration.Occurrences},
				)
			}
		}
	}

	results := &EvaluationResults{
		FileNameRuleMapper: mapper,
		Summary: struct {
			TotalFailedRules int
			FilesCount       int
			TotalPassedCount int
		}{
			TotalFailedRules: totalFailedCount,
			FilesCount:       filesCount,
			TotalPassedCount: totalPassedCount,
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

func countOccurrences(validationResult *Result) int {
	count := 0
	for _, err := range validationResult.Errors() {
		if err.Type() != "condition_then" && err.Type() != "number_all_of" {
			count = count + 1
		}
	}
	return count
}

func extractSkipAnnotations(configuration extractor.Configuration) map[string]string {
	skipAnnotations := make(map[string]string)
	if configurationMetadata, ok := configuration["metadata"].(map[string]interface{}); ok {
		if annotationsMap, ok := configurationMetadata["annotations"].(map[string]interface{}); ok {
			for annotationKey, annotationValue := range annotationsMap {
				if strings.Contains(annotationKey, SKIP_RULE_PREFIX) {
					skipAnnotations[annotationKey] = annotationValue.(string)
				}
			}
			return skipAnnotations
		}
	}

	return nil
}
