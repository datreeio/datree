package evaluation

import (
	"path/filepath"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCliClient struct {
	mock.Mock
}

func (m *mockCliClient) CreateEvaluation(createEvaluationRequest *cliClient.CreateEvaluationRequest) (int, error) {
	args := m.Called(createEvaluationRequest)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockCliClient) RequestEvaluation(evaluationRequest *cliClient.EvaluationRequest) (*cliClient.EvaluationResponse, error) {
	args := m.Called(evaluationRequest)
	return args.Get(0).(*cliClient.EvaluationResponse), args.Error(1)
}

func (m *mockCliClient) UpdateEvaluationValidation(request *cliClient.UpdateEvaluationValidationRequest) error {
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

func TestEvaluate(t *testing.T) {
	tests := []*evaluateTestCase{
		happy_flow_test(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedCliClient := &mockCliClient{}

			mockedCliClient.On("RequestEvaluation", mock.Anything).Return(tt.mock.cliClient.requestEvaluation.response, tt.mock.cliClient.requestEvaluation.err)
			mockedCliClient.On("UpdateEvaluationValidation", mock.Anything).Return(nil)

			evaluator := &Evaluator{
				cliClient: mockedCliClient,
				osInfo:    tt.args.osInfo,
			}

			actualResponse, _, _ := evaluator.Evaluate(tt.args.validFilesChan, tt.args.invalidFilesChan, tt.args.evaluationId)

			if tt.expected.isRequestEvaluationCalled {
				mockedCliClient.AssertCalled(t, "RequestEvaluation", mock.Anything)
			}

			if tt.expected.isUpdateEvaluationValidationCalled {
				mockedCliClient.AssertCalled(t, "UpdateEvaluationValidation")
			}

			assert.Equal(t, tt.expected.response.Summary, actualResponse.Summary)
			assert.Equal(t, tt.expected.response.FileNameRuleMapper, actualResponse.FileNameRuleMapper)
		})
	}
}

func TestCreateEvaluation(t *testing.T) {
	tests := []*evaluateTestCase{
		happy_flow_test(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedCliClient := &mockCliClient{}

			mockedCliClient.On("RequestEvaluation", mock.Anything).Return(tt.mock.cliClient.requestEvaluation.response, tt.mock.cliClient.requestEvaluation.err)
			mockedCliClient.On("UpdateEvaluationValidation", mock.Anything).Return(nil)

			evaluator := &Evaluator{
				cliClient: mockedCliClient,
				osInfo:    tt.args.osInfo,
			}

			actualResponse, _, _ := evaluator.Evaluate(tt.args.validFilesChan, tt.args.invalidFilesChan, tt.args.evaluationId)

			if tt.expected.isRequestEvaluationCalled {
				mockedCliClient.AssertCalled(t, "RequestEvaluation", mock.Anything)
			}

			if tt.expected.isUpdateEvaluationValidationCalled {
				mockedCliClient.AssertCalled(t, "UpdateEvaluationValidation")
			}

			assert.Equal(t, tt.expected.response.Summary, actualResponse.Summary)
			assert.Equal(t, tt.expected.response.FileNameRuleMapper, actualResponse.FileNameRuleMapper)
		})
	}
}

type evaluateArgs struct {
	validFilesChan   <-chan string
	invalidFilesChan <-chan string
	evaluationId     int
	osInfo           *OSInfo
}

type evaluateExpected struct {
	response                           *EvaluationResults
	errors                             []*Error
	err                                error
	isRequestEvaluationCalled          bool
	isCreateEvaluationCalled           bool
	isUpdateEvaluationValidationCalled bool
	isGetVersionMessageCalled          bool
}

type evaluateTestCase struct {
	name     string
	args     *evaluateArgs
	mock     *evaluatorMock
	expected *evaluateExpected
}

func happy_flow_test() *evaluateTestCase {
	return &evaluateTestCase{
		name: "should",
		args: &evaluateArgs{
			validFilesChan:   newFilesChan(),
			invalidFilesChan: newInvalidFilesChan(),
			evaluationId:     1,
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
					RulesCount       int
					TotalFailedRules int
					FilesCount       int
				}{
					RulesCount:       0,
					TotalFailedRules: 0,
					FilesCount:       1,
				},
			},
			errors:                             []*Error{},
			err:                                nil,
			isRequestEvaluationCalled:          true,
			isUpdateEvaluationValidationCalled: false,
		},
	}
}

func newFilesChan() chan string {
	files := make(chan string, 1)
	p, _ := filepath.Abs("../../internal/fixtures/kube/pass-all.yaml")
	files <- p
	close(files)
	return files
}

func newInvalidFilesChan() chan string {
	files := make(chan string, 1)
	files <- "path/path1/service.yaml"
	close(files)
	return files
}
