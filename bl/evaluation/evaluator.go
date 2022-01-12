package evaluation

import (
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
)

type CLIClient interface {
	RequestEvaluation(request *cliClient.EvaluationRequest) (*cliClient.EvaluationResponse, error)
	CreateEvaluation(request *cliClient.CreateEvaluationRequest) (*cliClient.CreateEvaluationResponse, error)
	SendFailedYamlValidation(request *cliClient.UpdateEvaluationValidationRequest) error
	SendFailedK8sValidation(request *cliClient.UpdateEvaluationValidationRequest) error
}

type Evaluator struct {
	cliClient CLIClient
	osInfo    *OSInfo
}

func New(c CLIClient) *Evaluator {
	return &Evaluator{
		cliClient: c,
		osInfo:    NewOSInfo(),
	}
}

type FileNameRuleMapper map[string]map[int]*Rule

type ResultType struct {
	EvaluationResults               *EvaluationResults
	NonInteractiveEvaluationResults *NonInteractiveEvaluationResults
}
type EvaluationResults struct {
	FileNameRuleMapper FileNameRuleMapper
	Summary            struct {
		TotalFailedRules int
		FilesCount       int
		TotalPassedCount int
	}
}

func (e *Evaluator) CreateEvaluation(cliId string, cliVersion string, k8sVersion string, policyName string) (*cliClient.CreateEvaluationResponse, error) {
	createEvaluationResponse, err := e.cliClient.CreateEvaluation(&cliClient.CreateEvaluationRequest{
		K8sVersion: &k8sVersion,
		CliId:      cliId,
		PolicyName: policyName,
		Metadata: &cliClient.Metadata{
			CliVersion:      cliVersion,
			Os:              e.osInfo.OS,
			PlatformVersion: e.osInfo.PlatformVersion,
			KernelVersion:   e.osInfo.KernelVersion,
		},
	})

	return createEvaluationResponse, err
}

func (e *Evaluator) UpdateFailedYamlValidation(invalidFiles []*validation.InvalidYamlFile, evaluationId int, stopEvaluation bool) error {
	invalidFilesPaths := []*string{}
	for _, file := range invalidFiles {
		invalidFilesPaths = append(invalidFilesPaths, &file.Path)
	}
	err := e.cliClient.SendFailedYamlValidation(&cliClient.UpdateEvaluationValidationRequest{
		EvaluationId:   evaluationId,
		InvalidFiles:   invalidFilesPaths,
		StopEvaluation: stopEvaluation,
	})
	return err
}

func (e *Evaluator) UpdateFailedK8sValidation(invalidFiles []*validation.InvalidK8sFile, evaluationId int, stopEvaluation bool) error {
	invalidFilesPaths := []*string{}
	for _, file := range invalidFiles {
		invalidFilesPaths = append(invalidFilesPaths, &file.Path)
	}
	err := e.cliClient.SendFailedK8sValidation(&cliClient.UpdateEvaluationValidationRequest{
		EvaluationId:   evaluationId,
		InvalidFiles:   invalidFilesPaths,
		StopEvaluation: stopEvaluation,
	})
	return err
}

func (e *Evaluator) Evaluate(filesConfigurations []*extractor.FileConfigurations, evaluationResponse *cliClient.CreateEvaluationResponse, isInteractiveMode bool) (ResultType, error) {

	if len(filesConfigurations) == 0 {
		return ResultType{}, nil
	}

	res, err := e.cliClient.RequestEvaluation(&cliClient.EvaluationRequest{
		EvaluationId: evaluationResponse.EvaluationId,
		Files:        filesConfigurations,
	})
	if err != nil {
		return ResultType{}, err
	}
	resultType := ResultType{}
	resultType.EvaluationResults = e.formatEvaluationResults(res.Results, len(filesConfigurations))
	if !isInteractiveMode {
		resultType.NonInteractiveEvaluationResults = e.formatNonInteractiveEvaluationResults(resultType.EvaluationResults, res.Results, evaluationResponse.PolicyName, evaluationResponse.RulesCount)
	}
	return resultType, nil
}

// This method creates a NonInteractiveEvaluationResults structure
// from EvaluationResults.
func (e *Evaluator) formatNonInteractiveEvaluationResults(evaluationResults *EvaluationResults, listEvaluationResult []*cliClient.EvaluationResult, policyName string, totalRulesInPolicy int) *NonInteractiveEvaluationResults {
	fileNameRuleMapper := evaluationResults.FileNameRuleMapper
	ruleMapper := make(map[int]string)
	for _, result := range listEvaluationResult {
		ruleId := getRuleId(result)
		ruleMapper[ruleId] = result.Rule.Identifier
	}
	nonInteractiveEvaluationResults := NonInteractiveEvaluationResults{}

	for fileName, rules := range fileNameRuleMapper {
		formattedEvaluationResults := FormattedEvaluationResults{}
		formattedEvaluationResults.FileName = fileName

		for _, rule := range rules {
			ruleResult := RuleResult{Identifier: ruleMapper[rule.ID], Name: rule.Name, MessageOnFailure: rule.FailSuggestion, OccurrencesDetails: rule.OccurrencesDetails}
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
		TotalRulesFailed:   evaluationResults.Summary.TotalFailedRules,
		TotalPassedCount:   evaluationResults.Summary.TotalPassedCount,
	}

	return &nonInteractiveEvaluationResults
}

func (e *Evaluator) formatEvaluationResults(evaluationResults []*cliClient.EvaluationResult, filesCount int) *EvaluationResults {
	mapper := make(map[string]map[int]*Rule)

	totalFailedCount := 0
	totalPassedCount := filesCount

	for _, result := range evaluationResults {
		for _, match := range result.Results.Matches {
			// file not already exists in mapper
			if _, exists := mapper[match.FileName]; !exists {
				mapper[match.FileName] = make(map[int]*Rule)
				totalPassedCount = totalPassedCount - 1
			}

			ruleId := getRuleId(result)

			// file and rule not already exists in mapper
			if _, exists := mapper[match.FileName][ruleId]; !exists {
				totalFailedCount++
				mapper[match.FileName][ruleId] = &Rule{
					ID:                 ruleId,
					Name:               result.Rule.Name,
					FailSuggestion:     result.Rule.FailSuggestion,
					OccurrencesDetails: []OccurrenceDetails{},
				}
			}

			mapper[match.FileName][ruleId].OccurrencesDetails = append(
				mapper[match.FileName][ruleId].OccurrencesDetails,
				OccurrenceDetails{MetadataName: match.MetadataName, Kind: match.Kind},
			)
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

func getRuleId(evaluationResult *cliClient.EvaluationResult) int {
	var ruleId int
	if evaluationResult.Rule.Origin.Type == "default" {
		ruleId = *evaluationResult.Rule.Origin.DefaultRuleId
	} else {
		ruleId = *evaluationResult.Rule.Origin.CustomRuleId
	}

	return ruleId
}
