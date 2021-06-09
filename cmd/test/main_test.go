package test

import (
	"testing"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/printer"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/mock"
)

type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) Evaluate(filesConfigurationsChan []*extractor.FileConfigurations, evaluationId int) (*evaluation.EvaluationResults, error) {
	args := m.Called(filesConfigurationsChan, evaluationId)
	return args.Get(0).(*evaluation.EvaluationResults), args.Error(1)
}

func (m *mockEvaluator) CreateEvaluation(cliId string, cliVersion string, k8sVersion string) (*cliClient.CreateEvaluationResponse, error) {
	args := m.Called(cliId, cliVersion, k8sVersion)
	return args.Get(0).(*cliClient.CreateEvaluationResponse), args.Error(1)
}

func (m *mockEvaluator) UpdateFailedYamlValidation(invalidYamlFiles []*validation.InvalidYamlFile, evaluationId int, stopEvaluation bool) error {
	args := m.Called(invalidYamlFiles, evaluationId, stopEvaluation)
	return args.Error(0)
}

func (m *mockEvaluator) UpdateFailedK8sValidation(invalidK8sFiles []*validation.InvalidK8sFile, evaluationId int, stopEvaluation bool) error {
	args := m.Called(invalidK8sFiles, evaluationId, stopEvaluation)
	return args.Error(0)
}

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string) {
	go func() {
		messages <- &messager.VersionMessage{
			CliVersion:   "1.2.3",
			MessageText:  "version message mock",
			MessageColor: "green"}
		close(messages)
	}()

	m.Called(messages, cliVersion)
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
}

type K8sValidatorMock struct {
	mock.Mock
}

func (kv *K8sValidatorMock) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *validation.InvalidK8sFile) {
	args := kv.Called(filesConfigurationsChan, concurrency)
	return args.Get(0).(chan *extractor.FileConfigurations), args.Get(1).(chan *validation.InvalidK8sFile)
}

func (kv *K8sValidatorMock) InitClient(k8sVersion string) {
}

type PrinterMock struct {
	mock.Mock
}

func (p *PrinterMock) PrintWarnings(warnings []printer.Warning) {
	p.Called(warnings)
}

func (p *PrinterMock) PrintSummaryTable(summary printer.Summary) {
	p.Called(summary)
}

func (p *PrinterMock) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

func (p *PrinterMock) PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string) {
	p.Called(evaluationSummary)
}

type ReaderMock struct {
	mock.Mock
}

func (rm *ReaderMock) FilterFiles(paths []string) ([]string, error) {
	args := rm.Called(paths)
	return args.Get(0).([]string), nil
}

func TestTestCommand(t *testing.T) {
	evaluationId := 444

	evaluationResults := &evaluation.EvaluationResults{
		FileNameRuleMapper: map[string]map[int]*evaluation.Rule{}, Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
			TotalPassedCount int
		}{RulesCount: 1, TotalFailedRules: 0, FilesCount: 0, TotalPassedCount: 1},
	}

	mockedEvaluator := &mockEvaluator{}
	mockedEvaluator.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(evaluationResults, nil)
	mockedEvaluator.On("CreateEvaluation", mock.Anything, mock.Anything, mock.Anything).Return(&cliClient.CreateEvaluationResponse{EvaluationId: evaluationId, K8sVersion: "1.18.0"}, nil)
	mockedEvaluator.On("UpdateFailedYamlValidation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockedEvaluator.On("UpdateFailedK8sValidation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything, mock.Anything)

	k8sValidatorMock := &K8sValidatorMock{}

	path := "valid/path"
	filesConfigurationsChan := newFilesConfigurationsChan(path)
	filesConfigurations := newFilesConfigurations(path)

	invelidK8sFilesChan := newInvalidK8sFilesChan()

	k8sValidatorMock.On("ValidateResources", mock.Anything, mock.Anything).Return(filesConfigurationsChan, invelidK8sFilesChan, newErrorsChan())
	k8sValidatorMock.On("InitClient", mock.Anything).Return()

	printerMock := &PrinterMock{}
	printerMock.On("PrintWarnings", mock.Anything)
	printerMock.On("PrintSummaryTable", mock.Anything)
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)
	printerMock.On("PrintEvaluationSummary", mock.Anything, mock.Anything)

	readerMock := &ReaderMock{}

	readerMock.On("FilterFiles", mock.Anything).Return([]string{"file/path"}, nil)

	ctx := &TestCommandContext{
		K8sValidator: k8sValidatorMock,
		Evaluator:    mockedEvaluator,
		LocalConfig:  &localConfig.LocalConfiguration{CliId: "134kh"},
		Messager:     messager,
		Printer:      printerMock,
		Reader:       readerMock,
	}

	test_testCommand_no_flags(t, mockedEvaluator, filesConfigurations, evaluationId, ctx)
	test_testCommand_json_output(t, mockedEvaluator, filesConfigurations, evaluationId, ctx)
	test_testCommand_yaml_output(t, mockedEvaluator, filesConfigurations, evaluationId, ctx)
}

func test_testCommand_no_flags(t *testing.T, evaluator *mockEvaluator, filesConfigurations []*extractor.FileConfigurations, evaluationId int, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{K8sVersion: "", Output: ""})

	evaluator.AssertCalled(t, "Evaluate", filesConfigurations, evaluationId)
}

func test_testCommand_json_output(t *testing.T, evaluator *mockEvaluator, filesConfigurations []*extractor.FileConfigurations, evaluationId int, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "json"})

	evaluator.AssertCalled(t, "Evaluate", filesConfigurations, evaluationId)
}

func test_testCommand_yaml_output(t *testing.T, evaluator *mockEvaluator, filesConfigurations []*extractor.FileConfigurations, evaluationId int, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "yaml"})

	evaluator.AssertCalled(t, "Evaluate", filesConfigurations, evaluationId)
}

func newFilesConfigurationsChan(path string) chan *extractor.FileConfigurations {
	filesConfigurationsChan := make(chan *extractor.FileConfigurations, 1)

	go func() {
		filesConfigurationsChan <- &extractor.FileConfigurations{
			FileName: path,
		}
		close(filesConfigurationsChan)
	}()

	return filesConfigurationsChan
}

func newFilesConfigurations(path string) []*extractor.FileConfigurations {
	var filesConfigurations []*extractor.FileConfigurations
	filesConfigurations = append(filesConfigurations, &extractor.FileConfigurations{
		FileName: path,
	})
	return filesConfigurations
}

func newInvalidK8sFilesChan() chan *validation.InvalidK8sFile {
	invalidFilesChan := make(chan *validation.InvalidK8sFile, 1)

	invalidFile := &validation.InvalidK8sFile{
		Path:             "invalid/path",
		ValidationErrors: []error{},
	}

	go func() {
		invalidFilesChan <- invalidFile
		close(invalidFilesChan)
	}()

	return invalidFilesChan
}

func newErrorsChan() chan error {
	invalidFilesChan := make(chan error, 1)

	close(invalidFilesChan)
	return invalidFilesChan
}
