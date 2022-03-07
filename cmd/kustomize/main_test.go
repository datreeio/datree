package kustomize

import (
	"testing"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/cmd/test"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/executor"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/stretchr/testify/mock"
)

func TestKustomizeTestCommand(t *testing.T) {
	tests := []*somethingTestCase{
		test_kustomize_run_method_success(),
	}
	for _, tt := range tests {
		t.Skip(tt.name)
		// t.Run(tt.name, func(t *testing.T) {
		// 	cmd := New(tt.testCtx, tt.kustomizeCtx)
		// 	got := cmd.RunE(cmd, tt.args)
		// 	assert.Equal(t, tt.expected, got)
		// })
	}
}

type somethingTestCase struct {
	name         string
	args         []string
	testCtx      *test.TestCommandContext
	kustomizeCtx *KustomizeContext
	expected     error
}

func test_kustomize_run_method_success() *somethingTestCase {
	return &somethingTestCase{
		name: "should return nil when kustomize run method is successful",
		args: []string{"./kustomization.yaml"},
		testCtx: &test.TestCommandContext{
			K8sValidator: &k8sValidatorMock{},
			Evaluator:    &mockEvaluator{},
			LocalConfig:  &LocalConfigMock{},
			Messager:     &mockMessager{},
			Printer:      &PrinterMock{},
			Reader:       &ReaderMock{},
		},
		kustomizeCtx: &KustomizeContext{
			CommandRunner: &mockKustomizeExecuter{},
		},
		expected: nil,
	}
}

// --- Mocks ---------------------------------------------------------------
type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) Evaluate(filesConfigurationsChan []*extractor.FileConfigurations, evaluationResponse *cliClient.CreateEvaluationResponse, isInteractiveMode bool) (evaluation.ResultType, error) {
	args := m.Called(filesConfigurationsChan, evaluationResponse, isInteractiveMode)
	return args.Get(0).(evaluation.ResultType), args.Error(1)
}

func (m *mockEvaluator) CreateEvaluation(cliId string, cliVersion string, k8sVersion string, policyName string, ciContext *ciContext.CIContext) (*cliClient.CreateEvaluationResponse, error) {
	args := m.Called(cliId, cliVersion, k8sVersion, policyName)
	return args.Get(0).(*cliClient.CreateEvaluationResponse), args.Error(1)
}

func (m *mockEvaluator) UpdateFailedYamlValidation(invalidYamlFiles []*extractor.InvalidFile, evaluationId int, stopEvaluation bool) error {
	args := m.Called(invalidYamlFiles, evaluationId, stopEvaluation)
	return args.Error(0)
}

func (m *mockEvaluator) UpdateFailedK8sValidation(invalidK8sFiles []*extractor.InvalidFile, evaluationId int, stopEvaluation bool) error {
	args := m.Called(invalidK8sFiles, evaluationId, stopEvaluation)
	return args.Error(0)
}

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) LoadVersionMessages(cliVersion string) chan *messager.VersionMessage {
	messages := make(chan *messager.VersionMessage, 1)
	m.Called(cliVersion)
	return messages
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
}

type mockKustomizeExecuter struct {
	mock.Mock
}

func (m *mockKustomizeExecuter) ExecuteKustomizeBin(args []string) ([]byte, error) {
	_args := m.Called(args)
	return _args.Get(0).([]byte), _args.Error(1)
}

func (m *mockKustomizeExecuter) RunCommand(name string, args []string) (executor.CommandOutput, error) {
	_args := m.Called(args)
	return _args.Get(0).(executor.CommandOutput), _args.Error(1)
}

func (m *mockKustomizeExecuter) BuildCommandDescription(dir string, name string, args []string) string {
	_args := m.Called(dir, name)
	return _args.Get(0).(string)
}

func (m *mockKustomizeExecuter) CreateTempFile(tempFilePrefix string, content []byte) (string, error) {
	_args := m.Called(tempFilePrefix, content)
	return _args.Get(0).(string), _args.Error(1)
}

type k8sValidatorMock struct {
	mock.Mock
}

func (kv *k8sValidatorMock) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile) {
	args := kv.Called(filesConfigurationsChan, concurrency)
	return args.Get(0).(chan *extractor.FileConfigurations), args.Get(1).(chan *extractor.InvalidFile)
}

func (kv *k8sValidatorMock) GetK8sFiles(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.FileConfigurations) {
	args := kv.Called(filesConfigurationsChan, concurrency)
	return args.Get(0).(chan *extractor.FileConfigurations), args.Get(1).(chan *extractor.FileConfigurations)
}

func (kv *k8sValidatorMock) InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) {
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

func (p *PrinterMock) PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string) {
	p.Called(evaluationSummary)
}

func (p *PrinterMock) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

func (p *PrinterMock) PrintPromptMessage(promptMessage string) {
	p.Called(promptMessage)
}

func (p *PrinterMock) SetTheme(theme *printer.Theme) {
	p.Called(theme)
}

type ReaderMock struct {
	mock.Mock
}

func (rm *ReaderMock) FilterFiles(paths []string) ([]string, error) {
	args := rm.Called(paths)
	return args.Get(0).([]string), nil
}

type LocalConfigMock struct {
	mock.Mock
}

func (lc *LocalConfigMock) GetLocalConfiguration() (*localConfig.ConfigContent, error) {
	lc.Called()
	return &localConfig.ConfigContent{CliId: "134kh"}, nil
}
