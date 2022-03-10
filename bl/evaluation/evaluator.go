package evaluation

import (
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"

	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/yamlSchemaValidator"
)

type CLIClient interface {
	SendEvaluationResult(request *cliClient.EvaluationResultRequest) (*cliClient.SendEvaluationResultsResponse, error)
}

type Evaluator struct {
	cliClient CLIClient
	osInfo    *OSInfo
	ciContext *ciContext.CIContext
}

func New(c CLIClient) *Evaluator {
	return &Evaluator{
		cliClient: c,
		osInfo:    NewOSInfo(),
		ciContext: ciContext.Extract(),
	}
}

type RuleSchema struct {
	RuleName string
	Schema   map[string]interface{}
}

type FileNameRuleMapper map[string]map[string]*Rule

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
	CliId              string
	CliVersion         string
	K8sVersion         string
	PolicyName         string
	CiContext          *ciContext.CIContext
	RulesData          []cliClient.RuleData
	FilesData          []cliClient.FileData
	FailedYamlFiles    []string
	FailedK8sFiles     []string
	PolicyCheckResults map[string]map[string]cliClient.FailedRule
}

func (e *Evaluator) SendEvaluationResult(evaluationRequestData EvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error) {
	sendEvaluationResultsResponse, err := e.cliClient.SendEvaluationResult(&cliClient.EvaluationResultRequest{
		K8sVersion: evaluationRequestData.K8sVersion,
		ClientId:   evaluationRequestData.CliId,
		Token:      evaluationRequestData.CliId,
		PolicyName: evaluationRequestData.PolicyName,
		Metadata: &cliClient.Metadata{
			CliVersion:      evaluationRequestData.CliVersion,
			Os:              e.osInfo.OS,
			PlatformVersion: e.osInfo.PlatformVersion,
			KernelVersion:   e.osInfo.KernelVersion,
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
	RawResults       map[string]map[string]cliClient.FailedRule
	RulesCount       int
}

func (e *Evaluator) Evaluate(policyCheckData PolicyCheckData) (PolicyCheckResultData, error) {

	if len(policyCheckData.FilesConfigurations) == 0 {
		return PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, map[string]map[string]cliClient.FailedRule{}, 0}, nil
	}

	yamlSchemaValidatorInst := yamlSchemaValidator.New()
	rulesCount := len(policyCheckData.Policy.Rules)

	// map of files paths to map of rules to failed rule data
	failedRulesByFiles := make(map[string]map[string]cliClient.FailedRule)

	rulesData := []cliClient.RuleData{}
	var filesData []cliClient.FileData

	emptyPolicyCheckResult := PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}

	for _, filesConfiguration := range policyCheckData.FilesConfigurations {
		filesData = append(filesData, cliClient.FileData{FilePath: filesConfiguration.FileName, ConfigurationsCount: len(filesConfiguration.Configurations)})

		for _, configuration := range filesConfiguration.Configurations {
			for _, ruleSchema := range policyCheckData.Policy.Rules {
				rulesData = append(rulesData, cliClient.RuleData{Identifier: ruleSchema.RuleIdentifier, Name: ruleSchema.RuleName})

				configurationKind, configurationName := extractConfigurationInfo(configuration)

				configurationJson, err := json.Marshal(configuration)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				ruleSchemaJson, err := json.Marshal(ruleSchema.Schema)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				validationResult, err := yamlSchemaValidatorInst.Validate(string(ruleSchemaJson), string(configurationJson))

				if err != nil {
					return emptyPolicyCheckResult, err
				}

				failedRulesByFiles = calculateFailedRulesByFiles(validationResult, filesConfiguration.FileName, ruleSchema, configurationName, configurationKind)
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

// This method creates a NonInteractiveEvaluationResults structure
// from EvaluationResults.
func (e *Evaluator) formatNonInteractiveEvaluationResults(formattedEvaluationResults *EvaluationResults, evaluationResults map[string]map[string]cliClient.FailedRule, policyName string, totalRulesInPolicy int) *NonInteractiveEvaluationResults {
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

func (e *Evaluator) formatEvaluationResults(evaluationResults map[string]map[string]cliClient.FailedRule, filesCount int) *EvaluationResults {
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
			if _, exists := mapper[filePath][ruleIdentifier]; !exists {
				totalFailedCount++
				mapper[filePath][ruleIdentifier] = &Rule{
					Identifier:         ruleIdentifier,
					Name:               failedRuleData.Name,
					MessageOnFailure:   failedRuleData.MessageOnFailure,
					OccurrencesDetails: []OccurrenceDetails{},
				}
			}

			for _, configuration := range failedRuleData.Configurations {
				mapper[filePath][ruleIdentifier].OccurrencesDetails = append(
					mapper[filePath][ruleIdentifier].OccurrencesDetails,
					OccurrenceDetails{MetadataName: configuration.Name, Kind: configuration.Kind},
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
	kind := configuration["kind"].(string)
	metadata := configuration["metadata"]
	name := metadata.(map[string]interface{})["name"].(string)

	return name, kind
}

type Result = gojsonschema.Result

func calculateFailedRulesByFiles(validationResult *Result, fileName string, ruleSchema policy_factory.RuleSchema, configurationName string, configurationKind string) map[string]map[string]cliClient.FailedRule {
	failedRulesByFiles := make(map[string]map[string]cliClient.FailedRule)

	if len(validationResult.Errors()) > 0 {
		configurationData := cliClient.Configuration{Name: configurationName, Kind: configurationKind, Occurrences: len(validationResult.Errors())}

		if fileData, ok := failedRulesByFiles[fileName]; ok {
			if ruleData, ok := fileData[ruleSchema.RuleIdentifier]; ok {
				ruleData.Configurations = append(ruleData.Configurations, configurationData)
				failedRulesByFiles[fileName][ruleSchema.RuleIdentifier] = ruleData
			} else {
				failedRulesByFiles[fileName][ruleSchema.RuleIdentifier] = cliClient.FailedRule{Name: ruleSchema.RuleName, MessageOnFailure: ruleSchema.MessageOnFailure, Configurations: []cliClient.Configuration{configurationData}}
			}
		} else {
			failedRulesByFiles[fileName] = map[string]cliClient.FailedRule{ruleSchema.RuleIdentifier: cliClient.FailedRule{Name: ruleSchema.RuleName, MessageOnFailure: ruleSchema.MessageOnFailure, Configurations: []cliClient.Configuration{configurationData}}}
		}
	}

	return failedRulesByFiles
}
