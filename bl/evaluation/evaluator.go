package evaluation

import (
	"github.com/datreeio/datree/bl/validation"
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
		TotalPassedCount int
	}
}

type Error struct {
	Message  string
	Filename string
}

func (e *Evaluator) CreateEvaluation(cliId string, cliVersion string, k8sVersion string) (int, error) {
	evaluationId, err := e.cliClient.CreateEvaluation(&cliClient.CreateEvaluationRequest{
		K8sVersion: k8sVersion,
		CliId:      cliId,
		Metadata: &cliClient.Metadata{
			CliVersion:      cliVersion,
			Os:              e.osInfo.OS,
			PlatformVersion: e.osInfo.PlatformVersion,
			KernelVersion:   e.osInfo.KernelVersion,
		},
	})

	return evaluationId, err
}

func (e *Evaluator) Evaluate(validFilesPathsChan chan string, invalidFilesPathsChan chan *validation.InvalidFile, evaluationId int) (*EvaluationResults, []*validation.InvalidFile, []*extractor.FileConfiguration, []*Error, error) {
	filesConfigurations, invalidFiles, extractionErrors := e.extractFilesConfigurations(validFilesPathsChan, invalidFilesPathsChan)

	invalidFilesPaths := []*string{}
	for _, file := range invalidFiles {
		invalidFilesPaths = append(invalidFilesPaths, &file.Path)
	}

	if len(invalidFiles) > 0 {
		stopEvaluation := len(filesConfigurations) == 0 // NOTICE: validFilesPathsChan surely closed and empty
		err := e.cliClient.UpdateEvaluationValidation(&cliClient.UpdateEvaluationValidationRequest{
			EvaluationId:   evaluationId,
			InvalidFiles:   invalidFilesPaths,
			StopEvaluation: stopEvaluation,
		})

		if stopEvaluation {
			return nil, invalidFiles, filesConfigurations, extractionErrors, err
		}
	}

	if len(filesConfigurations) > 0 {
		res, err := e.cliClient.RequestEvaluation(&cliClient.EvaluationRequest{
			EvaluationId: evaluationId,
			Files:        filesConfigurations,
		})
		if err != nil {
			return nil, invalidFiles, filesConfigurations, extractionErrors, err
		}

		results := e.formatEvaluationResults(res.Results, len(filesConfigurations))
		return results, invalidFiles, filesConfigurations, extractionErrors, nil
	}

	return nil, invalidFiles, filesConfigurations, extractionErrors, nil
}

func (e *Evaluator) extractFilesConfigurations(validFilesPathsChan chan string, invalidFilesPathsChan chan *validation.InvalidFile) ([]*extractor.FileConfiguration, []*validation.InvalidFile, []*Error) {
	invalidFiles := []*validation.InvalidFile{}
	var files []*extractor.FileConfiguration
	var errors []*Error

	readFromValidDone := false
	readFromInvalidDone := false
	for {
		select {
		case path, ok := <-validFilesPathsChan:
			if !ok {
				readFromValidDone = true
			} else {
				file, err := extractor.ExtractConfiguration(path)
				if file != nil {
					files = append(files, &extractor.FileConfiguration{
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
		case file, ok := <-invalidFilesPathsChan:
			if !ok {
				readFromInvalidDone = true
			} else {
				invalidFiles = append(invalidFiles, file)
			}
		}
		if readFromValidDone && readFromInvalidDone {
			break
		}
	}

	return files, invalidFiles, errors
}

func (e *Evaluator) formatEvaluationResults(evaluationResults []*cliClient.EvaluationResult, filesCount int) *EvaluationResults {
	mapper := make(map[string]map[int]*Rule)

	totalRulesCount := len(evaluationResults)
	totalFailedCount := 0
	totalPassedCount := filesCount

	for _, result := range evaluationResults {
		for _, match := range result.Results.Matches {
			// file not already exists in mapper
			if _, exists := mapper[match.FileName]; !exists {
				mapper[match.FileName] = make(map[int]*Rule)
				totalPassedCount = totalPassedCount - 1
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
			TotalPassedCount int
		}{
			RulesCount:       totalRulesCount,
			TotalFailedRules: totalFailedCount,
			FilesCount:       filesCount,
			TotalPassedCount: totalPassedCount,
		},
	}

	return results
}
