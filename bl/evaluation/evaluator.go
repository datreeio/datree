package evaluation

import (
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
)

type CLIClient interface {
	RequestEvaluation(request *cliClient.EvaluationRequest) (*cliClient.EvaluationResponse, error)
	CreateEvaluation(request *cliClient.CreateEvaluationRequest) (int, error)
	UpdateEvaluationValidation(request *cliClient.UpdateEvaluationValidationRequest) error
}

type Evaluator struct {
	cliClient             CLIClient
	osInfo                *OSInfo
	extractConfigurations func(path string) (*extractor.FileConfiguration, *extractor.Error)
}

func New(c CLIClient) *Evaluator {
	return &Evaluator{
		cliClient: c,
		osInfo:    NewOSInfo(),
	}
}

type EvaluationResults struct {
	FileNameRuleMapper map[string]map[int]*Rule
	Summary            struct {
		RulesCount       int
		TotalFailedRules int
		FilesCount       int
	}
}

type Error struct {
	Message  string
	Filename string
}

func (e *Evaluator) CreateEvaluation(cliId string, cliVersion string) (int, error) {
	evaluationId, err := e.cliClient.CreateEvaluation(&cliClient.CreateEvaluationRequest{
		CliId: cliId,
		Metadata: &cliClient.Metadata{
			CliVersion:      cliVersion,
			Os:              e.osInfo.OS,
			PlatformVersion: e.osInfo.PlatformVersion,
			KernelVersion:   e.osInfo.KernelVersion,
		},
	})

	return evaluationId, err
}

func (e *Evaluator) Evaluate(validFilesPathsChan <-chan string, invalidFilesPaths []*string, evaluationId int) (*EvaluationResults, []*Error, error) {
	filesConfigurations, errors := e.extractFilesConfigurations(validFilesPathsChan)

	if len(invalidFilesPaths) > 0 {
		stopEvaluation := len(validFilesPathsChan) == 0 // NOTICE: validFilesPathsChan surely closed and empty
		err := e.cliClient.UpdateEvaluationValidation(&cliClient.UpdateEvaluationValidationRequest{
			EvaluationId:   evaluationId,
			InvalidFiles:   invalidFilesPaths,
			StopEvaluation: stopEvaluation,
		})

		return nil, errors, err
	}

	if len(filesConfigurations) > 0 {
		res, err := e.cliClient.RequestEvaluation(&cliClient.EvaluationRequest{
			EvaluationId: evaluationId,
			Files:        filesConfigurations,
		})
		if err != nil {
			return nil, errors, err
		}

		results := e.formatEvaluationResults(res.Results, len(filesConfigurations))
		return results, errors, nil
	}

	return nil, errors, nil
}

func (e *Evaluator) extractFilesConfigurations(filesPathsChan <-chan string) ([]*cliClient.FileConfiguration, []*Error) {
	var files []*cliClient.FileConfiguration
	var errors []*Error

	done := false
	for {
		path, ok := <-filesPathsChan
		if !ok {
			done = true
		} else {
			file, err := extractor.ExtractConfiguration(path)
			if file != nil {
				files = append(files, &cliClient.FileConfiguration{
					FileName:       file.FileName,
					Configurations: file.Configurations,
				})
			}

			if err != nil {
				errors = append(errors, &Error{
					Message:  err.Message,
					Filename: err.Filename,
				})
			}
		}

		if done {
			break
		}
	}

	return files, errors
}

func (e *Evaluator) formatEvaluationResults(evaluationResults []*cliClient.EvaluationResult, filesCount int) *EvaluationResults {
	mapper := make(map[string]map[int]*Rule)

	totalRulesCount := len(evaluationResults)
	totalFailedCount := 0

	for _, result := range evaluationResults {
		for _, match := range result.Results.Matches {
			// file not already exists in mapper
			if _, exists := mapper[match.FileName]; !exists {
				mapper[match.FileName] = make(map[int]*Rule)
			}

			// file and rule not already exists in mapper
			if _, exists := mapper[match.FileName][result.Rule.ID]; !exists {
				totalFailedCount++
				mapper[match.FileName][result.Rule.ID] = &Rule{ID: result.Rule.ID, Name: result.Rule.Name, FailSuggestion: result.Rule.FailSuggestion, Count: 0}
			}

			mapper[match.FileName][result.Rule.ID].IncrementCount()
		}
	}

	results := &EvaluationResults{
		FileNameRuleMapper: mapper,
		Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{
			RulesCount:       totalRulesCount,
			TotalFailedRules: totalFailedCount,
			FilesCount:       filesCount,
		},
	}

	return results
}
