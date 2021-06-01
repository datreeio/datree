package evaluation

import (
	"github.com/datreeio/datree/bl/files"
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

type EvaluationResults struct {
	FileNameRuleMapper map[string]map[int]*Rule
	Summary            struct {
		RulesCount       int
		TotalFailedRules int
		FilesCount       int
		TotalPassedCount int
	}
}

func (e *Evaluator) CreateEvaluation(cliId string, cliVersion string, k8sVersion string) (*cliClient.CreateEvaluationResponse, error) {
	createEvaluationResponse, err := e.cliClient.CreateEvaluation(&cliClient.CreateEvaluationRequest{
		K8sVersion: &k8sVersion,
		CliId:      cliId,
		Metadata: &cliClient.Metadata{
			CliVersion:      cliVersion,
			Os:              e.osInfo.OS,
			PlatformVersion: e.osInfo.PlatformVersion,
			KernelVersion:   e.osInfo.KernelVersion,
		},
	})

	return createEvaluationResponse, err
}

func (e *Evaluator) UpdateFailedYamlValidation(invalidFiles []*validation.InvalidFile, evaluationId int, stopEvaluation bool) error {
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

func (e *Evaluator) UpdateFailedK8sValidation(invalidFiles []*validation.InvalidFile, evaluationId int, stopEvaluation bool) error {
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

func (e *Evaluator) Evaluate(filesConfigurationsChan chan *extractor.FileConfigurations, evaluationId int) (*EvaluationResults, error) {

	var filesConfigurations []*extractor.FileConfigurations

	for fileConfigurations := range filesConfigurationsChan {
		filesConfigurations = append(filesConfigurations, fileConfigurations)
	}

	if len(filesConfigurations) == 0 {
		return &EvaluationResults{}, nil
	}

	res, err := e.cliClient.RequestEvaluation(&cliClient.EvaluationRequest{
		EvaluationId: evaluationId,
		Files:        filesConfigurations,
	})
	if err != nil {
		return nil, err
	}

	results := e.formatEvaluationResults(res.Results, len(filesConfigurations))
	return results, nil
}

func (e *Evaluator) ExtractFilesConfigurations(paths []string, concurrency int) (chan *extractor.FileConfigurations, chan *validation.InvalidFile) {

	filesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidFilesChan := make(chan *validation.InvalidFile, concurrency)

	go func() {
		defer func() {
			close(filesConfigurationsChan)
			close(invalidFilesChan)
		}()

		for _, path := range paths {

			absolutePath, err := files.ToAbsolutePath(path)
			if err != nil {
				invalidFilesChan <- &validation.InvalidFile{Path: path, ValidationErrors: []error{err}}
				continue
			}

			content, err := extractor.ReadFileContent(absolutePath)
			if err != nil {
				invalidFilesChan <- &validation.InvalidFile{Path: absolutePath, ValidationErrors: []error{err}}
				continue
			}

			configurations, err := extractor.ParseYaml(content)
			if err != nil {
				invalidFilesChan <- &validation.InvalidFile{Path: absolutePath, ValidationErrors: []error{err}}
				continue
			}

			filesConfigurationsChan <- &extractor.FileConfigurations{FileName: absolutePath, Configurations: *configurations}
		}
	}()

	return filesConfigurationsChan, invalidFilesChan
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
