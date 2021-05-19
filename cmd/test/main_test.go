package test

import (
	"testing"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/localConfig"
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

func (m *mockEvaluator) Evaluate(validFilesPathsChan <-chan string, invalidFilesPaths []*string, evaluationId int) (*evaluation.EvaluationResults, []*evaluation.Error, error) {
	args := m.Called(validFilesPathsChan, invalidFilesPaths, evaluationId)
	return args.Get(0).(*evaluation.EvaluationResults), args.Get(1).([]*evaluation.Error), args.Error(2)
}

func (m *mockEvaluator) CreateEvaluation(cliId string, cliVersion string) (int, error) {
	args := m.Called(cliId, cliVersion)
	return args.Get(0).(int), args.Error(1)
}

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string) {
	m.Called(messages, cliVersion)
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
}
func TestTestCommand(t *testing.T) {
	mockedEvaluator := &mockEvaluator{}

	mockedEvaluateResponse := &evaluation.EvaluationResults{
		FileNameRuleMapper: map[string]map[int]*evaluation.Rule{}, Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{RulesCount: 1, TotalFailedRules: 0, FilesCount: 0},
	}

	mockedEvaluator.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(mockedEvaluateResponse)
	mockedEvaluator.On("CreateEvaluation", mock.Anything, mock.Anything).Return(mockedEvaluateResponse)

	localConfigManager := &mockLocalConfigManager{}
	localConfigManager.On("GetConfiguration").Return(localConfig.LocalConfiguration{CliId: "134kh"}, nil)

	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything).Return(mockedMessagesChannel())

	ctx := &TestCommandContext{
		Evaluator:   mockedEvaluator,
		LocalConfig: &localConfig.LocalConfiguration{},
		Messager:    messager,
	}

	test_testCommand_no_flags(t, localConfigManager, mockedEvaluator, mockedEvaluateResponse, ctx)
	test_testCommand_json_output(t, localConfigManager, mockedEvaluator, mockedEvaluateResponse, ctx)
	test_testCommand_yaml_output(t, localConfigManager, mockedEvaluator, mockedEvaluateResponse, ctx)
}

func test_testCommand_no_flags(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *evaluation.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "")
}

func test_testCommand_json_output(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *evaluation.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "json"})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "json")
}

func test_testCommand_yaml_output(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *evaluation.EvaluationResults, ctx *TestCommandContext) {
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
