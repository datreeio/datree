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
	cliClient CLIClient
	osInfo    *OSInfo
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

func (e *Evaluator) Evaluate(validFilesChan <-chan string, invalidFilesChan <-chan string, evaluationId int) (*EvaluationResults, []*Error, error) {
	validFiles, invalidFiles, errors := e.aggregateFiles(validFilesChan, invalidFilesChan)
	stopEvaluation := len(validFiles) == 0

	if len(invalidFiles) > 0 {
		e.cliClient.UpdateEvaluationValidation(&cliClient.UpdateEvaluationValidationRequest{
			EvaluationId:   evaluationId,
			InvalidFiles:   invalidFiles,
			StopEvaluation: stopEvaluation,
		})
	}

	if !stopEvaluation {
		res, err := e.cliClient.RequestEvaluation(&cliClient.EvaluationRequest{
			EvaluationId: evaluationId,
			Files:        validFiles,
		})
		if err != nil {
			return nil, errors, err
		}

		results := e.aggregateEvaluationResults(res.Results, len(validFiles))
		return results, errors, nil
	}

	return nil, errors, nil
}

func (e *Evaluator) aggregateFiles(validFilesChan <-chan string, invalidFilesChan <-chan string) ([]*cliClient.FileConfiguration, []*string, []*Error) {
	var validFiles []*cliClient.FileConfiguration
	var invalidFiles []*string
	var errors []*Error

	for {
		validDone, invalidDone := false, false
		select {
		case validFile, ok := <-validFilesChan:
			if !ok {
				validDone = true
			} else {
				go e.appendValidFileConfiguration(validFile, validFiles, errors)
			}
		case invalidFile, ok := <-invalidFilesChan:
			if !ok {
				invalidDone = true
			} else {
				go e.appendInvalidFile(invalidFile, invalidFiles)
			}
		}
		if invalidDone && validDone {
			break
		}
	}

	return validFiles, invalidFiles, errors
}

func (e *Evaluator) appendValidFileConfiguration(path string, files []*cliClient.FileConfiguration, errors []*Error) {
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

func (e *Evaluator) appendInvalidFile(path string, files []*string) {
	files = append(files, &path)
}

func (e *Evaluator) aggregateEvaluationResults(evaluationResults []*cliClient.EvaluationResult, filesCount int) *EvaluationResults {
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
