package evaluation

import (
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"

	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
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

type FileNameRuleMapper map[string]map[string]*Rule
type FailedRulesByFiles map[string]map[string]cliClient.FailedRule

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

func (e *Evaluator) SendEvaluationResult(evaluationRequestData EvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error) {
	sendEvaluationResultsResponse, err := e.cliClient.SendEvaluationResult(&cliClient.EvaluationResultRequest{
		K8sVersion: evaluationRequestData.K8sVersion,
		ClientId:   evaluationRequestData.ClientId,
		Token:      evaluationRequestData.Token,
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
	RawResults       FailedRulesByFiles
	RulesCount       int
}

func (e *Evaluator) Evaluate(policyCheckData PolicyCheckData) (PolicyCheckResultData, error) {
	rulesCount := len(policyCheckData.Policy.Rules)

	if len(policyCheckData.FilesConfigurations) == 0 {
		return PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, FailedRulesByFiles{}, rulesCount}, nil
	}

	jsonSchemaValidator := jsonSchemaValidator.New()

	// map of files paths to map of rules to failed rule data
	failedRulesByFiles := make(FailedRulesByFiles)

	rulesData := []cliClient.RuleData{}
	var filesData []cliClient.FileData

	emptyPolicyCheckResult := PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}

	for _, filesConfiguration := range policyCheckData.FilesConfigurations {
		filesData = append(filesData, cliClient.FileData{FilePath: filesConfiguration.FileName, ConfigurationsCount: len(filesConfiguration.Configurations)})

		for _, configuration := range filesConfiguration.Configurations {
			for _, ruleWithSchema := range policyCheckData.Policy.Rules {
				rulesData = append(rulesData, cliClient.RuleData{Identifier: ruleWithSchema.RuleIdentifier, Name: ruleWithSchema.RuleName})

				configurationName, configurationKind := extractConfigurationInfo(configuration)

				configurationJson, err := json.Marshal(configuration)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				ruleSchemaJson, err := json.Marshal(ruleWithSchema.Schema)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				validationResult, err := jsonSchemaValidator.ValidateYamlSchema(string(ruleSchemaJson), string(configurationJson))

				if err != nil {
					return emptyPolicyCheckResult, err
				}

				failedRulesByFiles = calculateFailedRulesByFiles(failedRulesByFiles, validationResult, filesConfiguration.FileName, ruleWithSchema, configurationName, configurationKind)
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
	kind := configuration["kind"].(string)
	metadata := configuration["metadata"]

	nonStringName := metadata.(map[string]interface{})["name"]
	var name string
	if nonStringName != nil {
		name = nonStringName.(string)
	} else {
		name = ""
	}

	return name, kind
}

type Result = gojsonschema.Result

func calculateFailedRulesByFiles(currentFailedRulesByFiles FailedRulesByFiles, validationResult *Result, fileName string, rule policy_factory.RuleWithSchema, configurationName string, configurationKind string) map[string]map[string]cliClient.FailedRule {
	occurrences := countOccurrences(validationResult)
	if occurrences > 0 {
		configurationData := cliClient.Configuration{Name: configurationName, Kind: configurationKind, Occurrences: occurrences}

		if fileData, ok := currentFailedRulesByFiles[fileName]; ok {
			if ruleData, ok := fileData[rule.RuleIdentifier]; ok {
				ruleData.Configurations = append(ruleData.Configurations, configurationData)
				currentFailedRulesByFiles[fileName][rule.RuleIdentifier] = ruleData
			} else {
				currentFailedRulesByFiles[fileName][rule.RuleIdentifier] = cliClient.FailedRule{Name: rule.RuleName, MessageOnFailure: rule.MessageOnFailure, Configurations: []cliClient.Configuration{configurationData}}
			}
		} else {
			currentFailedRulesByFiles[fileName] = map[string]cliClient.FailedRule{rule.RuleIdentifier: {Name: rule.RuleName, MessageOnFailure: rule.MessageOnFailure, Configurations: []cliClient.Configuration{configurationData}}}
		}
	}

	return currentFailedRulesByFiles
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
