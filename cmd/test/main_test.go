package test

import (
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/printer"
	"testing"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/mock"
)

type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) Evaluate(validFilesPathsChan chan string, invalidFilesPaths chan *validation.InvalidFile, evaluationId int) (*evaluation.EvaluationResults, []*validation.InvalidFile, []*extractor.FileConfiguration, []*evaluation.Error, error) {
	args := m.Called(validFilesPathsChan, invalidFilesPaths, evaluationId)
	return args.Get(0).(*evaluation.EvaluationResults), args.Get(1).([]*validation.InvalidFile), args.Get(2).([]*extractor.FileConfiguration), args.Get(3).([]*evaluation.Error), args.Error(4)
}

func (m *mockEvaluator) CreateEvaluation(cliId string, cliVersion string, k8sVersion string) (int, error) {
	args := m.Called(cliId, cliVersion, k8sVersion)
	return args.Get(0).(int), args.Error(1)
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

func (kv *K8sValidatorMock) ValidateResources(paths []string) (chan string, chan *validation.InvalidFile, chan error) {
	args := kv.Called(paths)
	return args.Get(0).(chan string), args.Get(1).(chan *validation.InvalidFile), args.Get(2).(chan error)
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

func (p *PrinterMock) PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary) {
	p.Called(evaluationSummary)
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

	invalidFiles := []*validation.InvalidFile{}
	var fileConfigurations []*extractor.FileConfiguration
	var evaluationErrors []*evaluation.Error

	mockedEvaluator := &mockEvaluator{}
	mockedEvaluator.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(evaluationResults, invalidFiles, fileConfigurations, evaluationErrors, nil)
	mockedEvaluator.On("CreateEvaluation", mock.Anything, mock.Anything, mock.Anything).Return(evaluationId, nil)

	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything, mock.Anything)

	k8sValidatorMock := &K8sValidatorMock{}

	validFilesPathChan := newValidFilesPathsChan()
	invalidFilesChan := newInvalidFilesChan()

	k8sValidatorMock.On("ValidateResources", mock.Anything).Return(validFilesPathChan, invalidFilesChan, newErrorsChan())

	printerMock := &PrinterMock{}
	printerMock.On("PrintWarnings", mock.Anything)
	printerMock.On("PrintSummaryTable", mock.Anything)
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)
	printerMock.On("PrintEvaluationSummary", mock.Anything)

	ctx := &TestCommandContext{
		K8sValidator: k8sValidatorMock,
		Evaluator:    mockedEvaluator,
		LocalConfig:  &localConfig.LocalConfiguration{CliId: "134kh"},
		Messager:     messager,
		Printer:      printerMock,
	}

	test_testCommand_no_flags(t, mockedEvaluator, validFilesPathChan, invalidFilesChan, evaluationId, ctx)
	test_testCommand_json_output(t, mockedEvaluator, validFilesPathChan, invalidFilesChan, evaluationId, ctx)
	test_testCommand_yaml_output(t, mockedEvaluator, validFilesPathChan, invalidFilesChan, evaluationId, ctx)
}

func test_testCommand_no_flags(t *testing.T, evaluator *mockEvaluator, validFilesPathChan chan string, invalidFilesChan chan *validation.InvalidFile, evaluationId int, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{})

	evaluator.AssertCalled(t, "Evaluate", validFilesPathChan, invalidFilesChan, evaluationId)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
}

func test_testCommand_json_output(t *testing.T, evaluator *mockEvaluator, validFilesPathChan chan string, invalidFilesChan chan *validation.InvalidFile, evaluationId int, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "json"})

	evaluator.AssertCalled(t, "Evaluate", validFilesPathChan, invalidFilesChan, evaluationId)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
}

func test_testCommand_yaml_output(t *testing.T, evaluator *mockEvaluator, validFilesPathChan chan string, invalidFilesChan chan *validation.InvalidFile, evaluationId int, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "yaml"})

	evaluator.AssertCalled(t, "Evaluate", validFilesPathChan, invalidFilesChan, evaluationId)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
}

func newValidFilesPathsChan() chan string {
	validFilesChan := make(chan string, 1)

	go func() {
		validFilesChan <- "valid/path"
		close(validFilesChan)
	}()

	return validFilesChan
}

func newInvalidFilesChan() chan *validation.InvalidFile {
	invalidFilesChan := make(chan *validation.InvalidFile, 1)

	invalidFile := &validation.InvalidFile{
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
