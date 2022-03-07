package evaluation

//
//import (
//	"encoding/json"
//
//	policy_factory "github.com/datreeio/datree/bl/policy"
//
//	"github.com/datreeio/datree/pkg/ciContext"
//	"github.com/datreeio/datree/pkg/cliClient"
//	"github.com/datreeio/datree/pkg/extractor"
//	"github.com/datreeio/datree/pkg/yamlSchemaValidator"
//	"github.com/xeipuuv/gojsonschema"
//)
//
//type LocalEvaluator struct {
//	cliClient CLIClient
//	osInfo    *OSInfo
//	ciContext *ciContext.CIContext
//}
//
//func NewLocalEvaluator(c CLIClient) *LocalEvaluator {
//	return &LocalEvaluator{
//		cliClient: c,
//		osInfo:    NewOSInfo(),
//		ciContext: ciContext.Extract(),
//	}
//}
//
//type Result = gojsonschema.Result
//
//type YamlSchemaValidator interface {
//	Validate(yamlSchema string, yaml string) (*Result, error)
//}
//
////type Configuration struct {
////	Name string
////	Kind string
////}
////
////type FailedRule struct {
////	Name             string
////	MessageOnFailure string
////	Occurrences      int
////	Configurations   []Configuration
////}
////
////type RuleSchema struct {
////	RuleName string
////	Schema   map[string]interface{}
////}
////
////type LocalFileNameRuleMapper map[string]map[string]*NewRule
////
////type LocalEvaluationResults struct {
////	FileNameRuleMapper LocalFileNameRuleMapper
////	Summary            struct {
////		TotalFailedRules int
////		FilesCount       int
////		TotalPassedCount int
////	}
////}
////
////type LocalResultType struct {
////	EvaluationResults               *LocalEvaluationResults
////	NonInteractiveEvaluationResults *NonInteractiveEvaluationResults
////}
//
//func (e *LocalEvaluator) LocalEvaluate(filesConfigurations []*extractor.FileConfigurations, isInteractiveMode bool, policyName string, policy policy_factory.Policy) (LocalResultType, int, error) {
//
//	if len(filesConfigurations) == 0 {
//		return LocalResultType{}, 0, nil
//	}
//
//	rulesCount := len(policy.Rules)
//
//	failedDict := make(map[string]map[string]FailedRule)
//
//	for _, filesConfiguration := range filesConfigurations {
//		for _, configuration := range filesConfiguration.Configurations {
//			for _, rulesSchema := range policy.Rules {
//
//				kind := configuration["kind"].(string)
//				metadata := configuration["metadata"]
//				name := metadata.(map[string]interface{})["name"].(string)
//
//				configurationJson, _ := json.Marshal(configuration)
//				yamlSchemaValidatorInst := yamlSchemaValidator.New()
//
//				ruleSchemaJson, _ := json.Marshal(rulesSchema.Schema)
//
//				result, err := yamlSchemaValidatorInst.Validate(string(ruleSchemaJson), string(configurationJson))
//
//				if err != nil {
//					return LocalResultType{}, 0, err
//				}
//
//				if len(result.Errors()) > 0 {
//					configurationData := Configuration{name, kind}
//
//					if fileData, ok := failedDict[filesConfiguration.FileName]; ok {
//						if ruleData, ok := fileData[rulesSchema.RuleIdentifier]; ok {
//							ruleData.Occurrences = ruleData.Occurrences + len(result.Errors())
//							ruleData.Configurations = append(ruleData.Configurations, configurationData)
//							failedDict[filesConfiguration.FileName][rulesSchema.RuleIdentifier] = ruleData
//						} else {
//							failedDict[filesConfiguration.FileName][rulesSchema.RuleIdentifier] = FailedRule{rulesSchema.RuleName, rulesSchema.MessageOnFailure, len(result.Errors()), []Configuration{configurationData}}
//						}
//					} else {
//						failedDict[filesConfiguration.FileName] = map[string]FailedRule{rulesSchema.RuleIdentifier: FailedRule{rulesSchema.RuleName, rulesSchema.MessageOnFailure, len(result.Errors()), []Configuration{configurationData}}}
//					}
//				}
//			}
//		}
//
//	}
//
//	resultType := LocalResultType{}
//	resultType.EvaluationResults = e.localFormatEvaluationResults(failedDict, len(filesConfigurations))
//	if !isInteractiveMode {
//		resultType.NonInteractiveEvaluationResults = e.localFormatNonInteractiveEvaluationResults(resultType.EvaluationResults, failedDict, policyName, rulesCount)
//	}
//	return resultType, rulesCount, nil
//}
//
//// This method creates a NonInteractiveEvaluationResults structure
//// from EvaluationResults.
//func (e *LocalEvaluator) localFormatNonInteractiveEvaluationResults(formattedEvaluationResults *LocalEvaluationResults, evaluationResults map[string]map[string]FailedRule, policyName string, totalRulesInPolicy int) *NonInteractiveEvaluationResults {
//	fileNameRuleMapper := formattedEvaluationResults.FileNameRuleMapper
//	ruleMapper := make(map[string]string)
//
//	for filePath := range evaluationResults {
//		for ruleUniqueName := range evaluationResults[filePath] {
//			ruleMapper[ruleUniqueName] = ruleUniqueName
//		}
//	}
//
//	nonInteractiveEvaluationResults := NonInteractiveEvaluationResults{}
//
//	for fileName, rules := range fileNameRuleMapper {
//		formattedEvaluationResults := FormattedEvaluationResults{}
//		formattedEvaluationResults.FileName = fileName
//
//		for _, rule := range rules {
//			ruleResult := RuleResult{Identifier: ruleMapper[rule.Identifier], Name: rule.Name, MessageOnFailure: rule.MessageOnFailure, OccurrencesDetails: rule.OccurrencesDetails}
//			formattedEvaluationResults.RuleResults = append(
//				formattedEvaluationResults.RuleResults,
//				&ruleResult,
//			)
//		}
//		nonInteractiveEvaluationResults.FormattedEvaluationResults = append(
//			nonInteractiveEvaluationResults.FormattedEvaluationResults,
//			&formattedEvaluationResults,
//		)
//	}
//	nonInteractiveEvaluationResults.PolicySummary = &PolicySummary{
//		PolicyName:         policyName,
//		TotalRulesInPolicy: totalRulesInPolicy,
//		TotalRulesFailed:   formattedEvaluationResults.Summary.TotalFailedRules,
//		TotalPassedCount:   formattedEvaluationResults.Summary.TotalPassedCount,
//	}
//
//	return &nonInteractiveEvaluationResults
//}
//
//func (e *LocalEvaluator) localFormatEvaluationResults(evaluationResults map[string]map[string]FailedRule, filesCount int) *LocalEvaluationResults {
//	mapper := make(map[string]map[string]*NewRule)
//
//	totalFailedCount := 0
//	totalPassedCount := filesCount
//
//	for filePath := range evaluationResults {
//		if _, exists := mapper[filePath]; !exists {
//			mapper[filePath] = make(map[string]*NewRule)
//			totalPassedCount = totalPassedCount - 1
//		}
//
//		for ruleUniqueName, failedRuleData := range evaluationResults[filePath] {
//			// file and rule not already exists in mapper
//			if _, exists := mapper[filePath][ruleUniqueName]; !exists {
//				totalFailedCount++
//				mapper[filePath][ruleUniqueName] = &NewRule{
//					Identifier:         ruleUniqueName,
//					Name:               failedRuleData.Name,
//					MessageOnFailure:   failedRuleData.MessageOnFailure,
//					OccurrencesDetails: []OccurrenceDetails{},
//				}
//			}
//
//			for _, configuration := range failedRuleData.Configurations {
//				mapper[filePath][ruleUniqueName].OccurrencesDetails = append(
//					mapper[filePath][ruleUniqueName].OccurrencesDetails,
//					OccurrenceDetails{MetadataName: configuration.Name, Kind: configuration.Kind},
//				)
//			}
//		}
//	}
//
//	results := &LocalEvaluationResults{
//		FileNameRuleMapper: mapper,
//		Summary: struct {
//			TotalFailedRules int
//			FilesCount       int
//			TotalPassedCount int
//		}{
//			TotalFailedRules: totalFailedCount,
//			FilesCount:       filesCount,
//			TotalPassedCount: totalPassedCount,
//		},
//	}
//
//	return results
//}
//
//func localGetRuleId(evaluationResult *cliClient.EvaluationResult) int {
//	var ruleId int
//	if evaluationResult.Rule.Origin.Type == "default" {
//		ruleId = *evaluationResult.Rule.Origin.DefaultRuleId
//	} else {
//		ruleId = *evaluationResult.Rule.Origin.CustomRuleId
//	}
//
//	return ruleId
//}
