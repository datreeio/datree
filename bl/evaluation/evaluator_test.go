package evaluation

import (
	"path/filepath"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCliClient struct {
	mock.Mock
}

func (m *mockCliClient) CreateEvaluation(createEvaluationRequest *cliClient.CreateEvaluationRequest) (*cliClient.CreateEvaluationResponse, error) {
	args := m.Called(createEvaluationRequest)
	return args.Get(0).(*cliClient.CreateEvaluationResponse), args.Error(1)
}

func (m *mockCliClient) RequestEvaluation(evaluationRequest *cliClient.EvaluationRequest) (*cliClient.EvaluationResponse, error) {
	args := m.Called(evaluationRequest)
	return args.Get(0).(*cliClient.EvaluationResponse), args.Error(1)
}

func (m *mockCliClient) SendFailedYamlValidation(request *cliClient.UpdateEvaluationValidationRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *mockCliClient) SendFailedK8sValidation(request *cliClient.UpdateEvaluationValidationRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *mockCliClient) GetVersionMessage(cliVersion string, timeout int) (*cliClient.VersionMessage, error) {
	args := m.Called(cliVersion, timeout)
	return args.Get(0).(*cliClient.VersionMessage), args.Error(1)
}

type cliClientMockTestCase struct {
	createEvaluation struct {
		evaluationId int
		k8sVersion   string
		err          error
	}
	requestEvaluation struct {
		response *cliClient.EvaluationResponse
		err      error
	}
	updateEvaluationValidation struct {
		err error
	}
	getVersionMessage struct {
		response *cliClient.VersionMessage
		err      error
	}
}
type evaluatorMock struct {
	cliClient *cliClientMockTestCase
}

// TODO: add actual tests
func TestCreateEvaluation(t *testing.T) {
	t.Run("CreateEvaluation should succedd", func(t *testing.T) {
		mockedCliClient := &mockCliClient{}
		evaluator := &Evaluator{
			cliClient: mockedCliClient,
			osInfo: &OSInfo{
				OS:              "darwin",
				PlatformVersion: "1.2.3",
				KernelVersion:   "4.5.6",
			},
		}

		cliId := "test_token"
		cliVersion := "0.0.7"
		k8sVersion := "1.18.1"
		policyName := "Default"

		mockedCliClient.On("CreateEvaluation", mock.Anything).Return(&cliClient.CreateEvaluationResponse{EvaluationId: 1, K8sVersion: k8sVersion}, nil)

		expectedCreateEvaluationResponse := &cliClient.CreateEvaluationResponse{EvaluationId: 1, K8sVersion: k8sVersion}
		createEvaluationResponse, _ := evaluator.CreateEvaluation(cliId, cliVersion, k8sVersion, policyName)

		mockedCliClient.AssertCalled(t, "CreateEvaluation", mock.Anything)
		assert.Equal(t, expectedCreateEvaluationResponse, createEvaluationResponse)

	})
}

func TestEvaluate(t *testing.T) {
	tests := []*evaluateTestCase{
		request_evaluation_all_valid(),
		request_evaluation_all_invalid(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedCliClient := &mockCliClient{}

			mockedCliClient.On("RequestEvaluation", mock.Anything).Return(tt.mock.cliClient.requestEvaluation.response, tt.mock.cliClient.requestEvaluation.err)

			evaluator := &Evaluator{
				cliClient: mockedCliClient,
				osInfo:    tt.args.osInfo,
			}

			// TODO: define and check the rest of the values
			results, _ := evaluator.Evaluate(tt.args.validFilesConfigurations, tt.args.evaluationId)

			if tt.expected.isRequestEvaluationCalled {
				mockedCliClient.AssertCalled(t, "RequestEvaluation", mock.Anything)
				assert.Equal(t, tt.expected.response.Summary, results.Summary)
				assert.Equal(t, tt.expected.response.FileNameRuleMapper, results.FileNameRuleMapper)
			} else {
				mockedCliClient.AssertNotCalled(t, "RequestEvaluation")
			}
		})
	}
}

type evaluateArgs struct {
	validFilesConfigurations []*extractor.FileConfigurations
	evaluationId             int
	osInfo                   *OSInfo
	rulesCount               int
}

type evaluateExpected struct {
	response                  *EvaluationResults
	err                       error
	isRequestEvaluationCalled bool
	isCreateEvaluationCalled  bool
	isGetVersionMessageCalled bool
}

type evaluateTestCase struct {
	name     string
	args     *evaluateArgs
	mock     *evaluatorMock
	expected *evaluateExpected
}

func request_evaluation_all_valid() *evaluateTestCase {
	validFilePath := "../../internal/fixtures/kube/pass-all.yaml"

	return &evaluateTestCase{
		name: "should request validation without invalid files",
		args: &evaluateArgs{
			validFilesConfigurations: newFilesConfigurations(validFilePath),
			evaluationId:             1,
			osInfo: &OSInfo{
				OS:              "darwin",
				PlatformVersion: "1.2.3",
				KernelVersion:   "4.5.6",
			},
		},
		mock: &evaluatorMock{
			cliClient: &cliClientMockTestCase{
				createEvaluation: struct {
					evaluationId int
					k8sVersion   string
					err          error
				}{
					evaluationId: 1,
					err:          nil,
				},
				requestEvaluation: struct {
					response *cliClient.EvaluationResponse
					err      error
				}{
					response: &cliClient.EvaluationResponse{
						Results: []*cliClient.EvaluationResult{},
					},
					err: nil,
				},
				updateEvaluationValidation: struct{ err error }{
					err: nil,
				},
				getVersionMessage: struct {
					response *cliClient.VersionMessage
					err      error
				}{
					response: nil,
					err:      nil,
				},
			},
		},
		expected: &evaluateExpected{
			response: &EvaluationResults{
				FileNameRuleMapper: make(map[string]map[int]*Rule),
				Summary: struct {
					TotalFailedRules int
					FilesCount       int
					TotalPassedCount int
				}{
					TotalFailedRules: 0,
					FilesCount:       1,
					TotalPassedCount: 1,
				},
			},
			err:                       nil,
			isRequestEvaluationCalled: true,
		},
	}
}

func request_evaluation_all_invalid() *evaluateTestCase {
	return &evaluateTestCase{
		name: "should not request validation if there are no valid files",
		args: &evaluateArgs{
			validFilesConfigurations: []*extractor.FileConfigurations{},
			evaluationId:             1,
			osInfo: &OSInfo{
				OS:              "darwin",
				PlatformVersion: "1.2.3",
				KernelVersion:   "4.5.6",
			},
		},
		mock: &evaluatorMock{
			cliClient: &cliClientMockTestCase{
				createEvaluation: struct {
					evaluationId int
					k8sVersion   string
					err          error
				}{
					evaluationId: 1,
					err:          nil,
				},
				requestEvaluation: struct {
					response *cliClient.EvaluationResponse
					err      error
				}{
					response: &cliClient.EvaluationResponse{
						Results: []*cliClient.EvaluationResult{},
					},
					err: nil,
				},
				updateEvaluationValidation: struct{ err error }{
					err: nil,
				},
				getVersionMessage: struct {
					response *cliClient.VersionMessage
					err      error
				}{
					response: nil,
					err:      nil,
				},
			},
		},
		expected: &evaluateExpected{
			response: &EvaluationResults{
				FileNameRuleMapper: make(map[string]map[int]*Rule),
				Summary: struct {
					TotalFailedRules int
					FilesCount       int
					TotalPassedCount int
				}{
					TotalFailedRules: 0,
					FilesCount:       1,
					TotalPassedCount: 0,
				},
			},
			err:                       nil,
			isRequestEvaluationCalled: false,
		},
	}
}

func newFilesConfigurations(path string) []*extractor.FileConfigurations {
	var filesConfigurations []*extractor.FileConfigurations
	absolutePath, _ := filepath.Abs(path)
	filesConfigurations = append(filesConfigurations, &extractor.FileConfigurations{
		FileName: absolutePath,
	})
	return filesConfigurations
}
