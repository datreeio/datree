package test

import (
	"testing"

	"github.com/datreeio/datree/bl/evaluator"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/propertiesExtractor"
	"github.com/stretchr/testify/mock"
)

type mockLocalConfigManager struct {
	mock.Mock
}

func (m *mockLocalConfigManager) GetConfiguration() (localConfig.LocalConfiguration, error) {
	args := m.Called()
	return args.Get(0).(localConfig.LocalConfiguration), args.Error(1)
}

type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) PrintResults(results *evaluator.EvaluationResults, cliId string, output string) error {
	m.Called(results, cliId, output)
	return nil
}

func (m *mockEvaluator) PrintFileParsingErrors(errors []propertiesExtractor.FileError) {
	m.Called(errors)
}

func (m *mockEvaluator) Evaluate(paths []string, cliId string, evaluationConc int, cliVersion string) (*evaluator.EvaluationResults, []propertiesExtractor.FileError, error) {
	args := m.Called(paths, cliId, evaluationConc)
	return args.Get(0).(*evaluator.EvaluationResults), args.Get(1).([]propertiesExtractor.FileError), args.Error(2)
}

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) LoadVersionMessages(cliVersion string) <-chan *messager.VersionMessage {
	args := m.Called(cliVersion)
	return args.Get(0).(<-chan *messager.VersionMessage)
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
}
func TestTestCommand(t *testing.T) {
	mockedEvaluator := &mockEvaluator{}

	mockedEvaluateResponse := &evaluator.EvaluationResults{
		FileNameRuleMapper: map[string]map[int]*evaluator.Rule{}, Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{RulesCount: 1, TotalFailedRules: 0, FilesCount: 0},
	}

	mockedEvaluator.On("Evaluate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedEvaluateResponse, []propertiesExtractor.FileError{}, nil)
	mockedEvaluator.On("PrintFileParsingErrors", mock.Anything).Return()
	mockedEvaluator.On("PrintResults", mock.Anything, mock.Anything, mock.Anything).Return()

	localConfigManager := &mockLocalConfigManager{}
	localConfigManager.On("GetConfiguration").Return(localConfig.LocalConfiguration{CliId: "134kh"}, nil)

	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything).Return(mockedMessagesChannel())

	ctx := &TestCommandContext{
		Evaluator:   mockedEvaluator,
		LocalConfig: localConfigManager,
		Messager:    messager,
	}

	test_testCommand_no_flags(t, localConfigManager, mockedEvaluator, mockedEvaluateResponse, ctx)
	test_testCommand_json_output(t, localConfigManager, mockedEvaluator, mockedEvaluateResponse, ctx)
	test_testCommand_yaml_output(t, localConfigManager, mockedEvaluator, mockedEvaluateResponse, ctx)
}

func test_testCommand_no_flags(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *evaluator.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "")
}

func test_testCommand_json_output(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *evaluator.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "json"})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "json")
}

func test_testCommand_yaml_output(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *evaluator.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "yaml"})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "yaml")
}

func mockedMessagesChannel() <-chan *cliClient.VersionMessage {
	mock := make(chan *cliClient.VersionMessage)
	mock <- &cliClient.VersionMessage{
		CliVersion:   "1.2.3",
		MessageText:  "version message mock",
		MessageColor: "green"}
	close(mock)

	return mock
}
