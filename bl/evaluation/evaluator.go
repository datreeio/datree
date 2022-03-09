package evaluation

import (
	"encoding/json"

	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/yamlSchemaValidator"
)

type CLIClient interface {
	RequestEvaluation(request *cliClient.EvaluationRequest) (*cliClient.EvaluationResponse, error)
	SendLocalEvaluationResult(request *cliClient.LocalEvaluationResultRequest) (*cliClient.SendEvaluationResultsResponse, error)
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

type LocalEvaluationRequestData struct {
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

func (e *Evaluator) SendLocalEvaluationResult(localEvaluationRequestData LocalEvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error) {
	sendLocalEvaluationResultsResponse, err := e.cliClient.SendLocalEvaluationResult(&cliClient.LocalEvaluationResultRequest{
		K8sVersion: localEvaluationRequestData.K8sVersion,
		ClientId:   localEvaluationRequestData.CliId,
		Token:      localEvaluationRequestData.CliId,
		PolicyName: localEvaluationRequestData.PolicyName,
		Metadata: &cliClient.Metadata{
			CliVersion:      localEvaluationRequestData.CliVersion,
			Os:              e.osInfo.OS,
			PlatformVersion: e.osInfo.PlatformVersion,
			KernelVersion:   e.osInfo.KernelVersion,
			CIContext:       localEvaluationRequestData.CiContext,
		},
		FailedYamlFiles:    localEvaluationRequestData.FailedYamlFiles,
		FailedK8sFiles:     localEvaluationRequestData.FailedK8sFiles,
		AllExecutedRules:   localEvaluationRequestData.RulesData,
		AllEvaluatedFiles:  localEvaluationRequestData.FilesData,
		PolicyCheckResults: localEvaluationRequestData.PolicyCheckResults,
	})

	return sendLocalEvaluationResultsResponse, err
}

type DataForPolicyCheck struct {
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

func (e *Evaluator) Evaluate(dataForEvaluation DataForPolicyCheck) (PolicyCheckResultData, error) {

	if len(dataForEvaluation.FilesConfigurations) == 0 {
		return PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}, nil
	}

	rulesCount := len(dataForEvaluation.Policy.Rules)

	// map of files paths to map of rules to failed rule data
	failedDict := make(map[string]map[string]cliClient.FailedRule)
	rulesData := []cliClient.RuleData{}
	var filesData []cliClient.FileData

	emptyPolicyCheckResult := PolicyCheckResultData{FormattedResults{}, []cliClient.RuleData{}, []cliClient.FileData{}, nil, 0}

	for _, filesConfiguration := range dataForEvaluation.FilesConfigurations {
		filesData = append(filesData, cliClient.FileData{FilePath: filesConfiguration.FileName, ConfigurationsCount: len(filesConfiguration.Configurations)})

		for _, configuration := range filesConfiguration.Configurations {
			for _, rulesSchema := range dataForEvaluation.Policy.Rules {
				rulesData = append(rulesData, cliClient.RuleData{Identifier: rulesSchema.RuleIdentifier, Name: rulesSchema.RuleName})

				kind := configuration["kind"].(string)
				metadata := configuration["metadata"]
				name := metadata.(map[string]interface{})["name"].(string)

				configurationJson, err := json.Marshal(configuration)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				yamlSchemaValidatorInst := yamlSchemaValidator.New()

				ruleSchemaJson, err := json.Marshal(rulesSchema.Schema)
				if err != nil {
					return emptyPolicyCheckResult, err
				}

				result, err := yamlSchemaValidatorInst.Validate(string(ruleSchemaJson), string(configurationJson))

				if err != nil {
					return emptyPolicyCheckResult, err
				}

				if len(result.Errors()) > 0 {
					configurationData := cliClient.Configuration{Name: name, Kind: kind, Occurrences: len(result.Errors())}

					if fileData, ok := failedDict[filesConfiguration.FileName]; ok {
						if ruleData, ok := fileData[rulesSchema.RuleIdentifier]; ok {
							ruleData.Configurations = append(ruleData.Configurations, configurationData)
							failedDict[filesConfiguration.FileName][rulesSchema.RuleIdentifier] = ruleData
						} else {
							failedDict[filesConfiguration.FileName][rulesSchema.RuleIdentifier] = cliClient.FailedRule{Name: rulesSchema.RuleName, MessageOnFailure: rulesSchema.MessageOnFailure, Configurations: []cliClient.Configuration{configurationData}}
						}
					} else {
						failedDict[filesConfiguration.FileName] = map[string]cliClient.FailedRule{rulesSchema.RuleIdentifier: cliClient.FailedRule{Name: rulesSchema.RuleName, MessageOnFailure: rulesSchema.MessageOnFailure, Configurations: []cliClient.Configuration{configurationData}}}
					}
				}
			}
		}

	}

	formattedResults := FormattedResults{}
	formattedResults.EvaluationResults = e.formatEvaluationResults(failedDict, len(dataForEvaluation.FilesConfigurations))

	if !dataForEvaluation.IsInteractiveMode {
		formattedResults.NonInteractiveEvaluationResults = e.formatNonInteractiveEvaluationResults(formattedResults.EvaluationResults, failedDict, dataForEvaluation.PolicyName, rulesCount)
	}

	return PolicyCheckResultData{formattedResults, rulesData, filesData, failedDict, rulesCount}, nil
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
