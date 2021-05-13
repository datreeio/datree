package test

import (
	"testing"

	"github.com/datreeio/datree/bl"
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

func (m *mockEvaluator) PrintResults(results *bl.EvaluationResults, cliId string, output string) error {
	m.Called(results, cliId, output)
	return nil
}

func (m *mockEvaluator) PrintFileParsingErrors(errors []propertiesExtractor.FileError) {
	m.Called(errors)
}

func (m *mockEvaluator) Evaluate(paths []string, cliId string, evaluationConc int, cliVersion string) (*bl.EvaluationResults, []propertiesExtractor.FileError, error) {
	args := m.Called(paths, cliId, evaluationConc)
	return args.Get(0).(*bl.EvaluationResults), args.Get(1).([]propertiesExtractor.FileError), args.Error(2)
}

type mockVersionMessageClient struct {
	mock.Mock
}

func (m *mockVersionMessageClient) GetVersionMessage(cliVersion string) (*cliClient.VersionMessage, error) {
	args := m.Called(cliVersion)
	return args.Get(0).(*cliClient.VersionMessage), nil
}
func TestTestCommand(t *testing.T) {
	evaluator := &mockEvaluator{}
	mockedEvaluateResponse := &bl.EvaluationResults{
		FileNameRuleMapper: map[string]map[int]*bl.Rule{}, Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{RulesCount: 1, TotalFailedRules: 0, FilesCount: 0},
	}
	evaluator.On("Evaluate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedEvaluateResponse, []propertiesExtractor.FileError{}, nil)
	evaluator.On("PrintFileParsingErrors", mock.Anything).Return()
	evaluator.On("PrintResults", mock.Anything, mock.Anything, mock.Anything).Return()

	localConfigManager := &mockLocalConfigManager{}
	localConfigManager.On("GetConfiguration").Return(localConfig.LocalConfiguration{CliId: "134kh"}, nil)

	versionMessageClient := &mockVersionMessageClient{}
	versionMessageClient.On("GetVersionMessage", mock.Anything).Return(
		&cliClient.VersionMessage{
			CliVersion:   "1.2.3",
			MessageText:  "version message mock",
			MessageColor: "green"},
	)

	ctx := &TestCommandContext{
		Evaluator:            evaluator,
		LocalConfig:          localConfigManager,
		VersionMessageClient: versionMessageClient,
	}

	test_testCommand_no_flags(t, localConfigManager, evaluator, mockedEvaluateResponse, ctx)
	test_testCommand_json_output(t, localConfigManager, evaluator, mockedEvaluateResponse, ctx)
	test_testCommand_yaml_output(t, localConfigManager, evaluator, mockedEvaluateResponse, ctx)
}

func test_testCommand_no_flags(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *bl.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "")
}

func test_testCommand_json_output(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *bl.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "json"})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "json")
}

func test_testCommand_yaml_output(t *testing.T, localConfigManager *mockLocalConfigManager, evaluator *mockEvaluator, mockedEvaluateResponse *bl.EvaluationResults, ctx *TestCommandContext) {
	test(ctx, []string{"8/*"}, TestCommandFlags{Output: "yaml"})
	localConfigManager.AssertCalled(t, "GetConfiguration")

	evaluator.AssertCalled(t, "Evaluate", []string{"8/*"}, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh", "yaml")
}
