package evaluation

import (
	"github.com/datreeio/datree/bl/validation"
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
	// TODO: add actual tests
	tests := []*evaluateTestCase{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedCliClient := &mockCliClient{}

			mockedCliClient.On("CreateEvaluation", mock.Anything).Return(tt.mock.cliClient.createEvaluation.evaluationId, tt.mock.cliClient.createEvaluation.err)
			mockedCliClient.On("UpdateEvaluationValidation", mock.Anything).Return(nil)

			evaluator := &Evaluator{
				cliClient: mockedCliClient,
				osInfo:    tt.args.osInfo,
			}

			results, _, _, _, _ := evaluator.Evaluate(tt.args.validFilesChan, tt.args.invalidFilesChan, tt.args.evaluationId)

			if tt.expected.isRequestEvaluationCalled {
				mockedCliClient.AssertCalled(t, "RequestEvaluation", mock.Anything)
				assert.Equal(t, tt.expected.response.Summary, results.Summary)
				assert.Equal(t, tt.expected.response.FileNameRuleMapper, results.FileNameRuleMapper)
			}

			if tt.expected.isUpdateEvaluationValidationCalled {
				mockedCliClient.AssertCalled(t, "UpdateEvaluationValidation", mock.Anything)
			}

		})
	}
}

func TestCreateEvaluation(t *testing.T) {
	tests := []*evaluateTestCase{
		request_evaluation_all_invalid(),
		request_evaluation_all_valid(),
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

			// TODO: define and check the rest of the values
			results, _, _, _, _ := evaluator.Evaluate(tt.args.validFilesChan, tt.args.invalidFilesChan, tt.args.evaluationId)

			if tt.expected.isRequestEvaluationCalled {
				mockedCliClient.AssertCalled(t, "RequestEvaluation", mock.Anything)
				assert.Equal(t, tt.expected.response.Summary, results.Summary)
				assert.Equal(t, tt.expected.response.FileNameRuleMapper, results.FileNameRuleMapper)
			}

			if tt.expected.isUpdateEvaluationValidationCalled {
				mockedCliClient.AssertCalled(t, "UpdateEvaluationValidation", mock.Anything)
			}

		})
	}
}

type evaluateArgs struct {
	validFilesChan   chan string
	invalidFilesChan chan *validation.InvalidFile
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

func request_evaluation_all_valid() *evaluateTestCase {
	invalidPath := "path/path1/service.yaml"

	return &evaluateTestCase{
		name: "should request validation without invalid files",
		args: &evaluateArgs{
			validFilesChan:   newFilesChan(),
			invalidFilesChan: newInvalidFilesChan(invalidPath),
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
					TotalPassedCount int
				}{
					RulesCount:       0,
					TotalFailedRules: 0,
					FilesCount:       1,
					TotalPassedCount: 1,
				},
			},
			errors:                             []*Error{},
			err:                                nil,
			isRequestEvaluationCalled:          true,
			isUpdateEvaluationValidationCalled: false,
		},
	}
}

func request_evaluation_all_invalid() *evaluateTestCase {
	invalidPath := "path/path1/service.yaml"
	return &evaluateTestCase{
		name: "should request validation all files are invalid",
		args: &evaluateArgs{
			validFilesChan:   newFilesChan(),
			invalidFilesChan: newInvalidFilesChan(invalidPath),
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
					TotalPassedCount int
				}{
					RulesCount:       0,
					TotalFailedRules: 0,
					FilesCount:       1,
					TotalPassedCount: 1,
				},
			},
			errors:                             []*Error{},
			err:                                nil,
			isRequestEvaluationCalled:          false,
			isUpdateEvaluationValidationCalled: true,
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

func newInvalidFilesChan(invalidPath string) chan *validation.InvalidFile {
	invalidFilesChan := make(chan *validation.InvalidFile, 1)

	invalidFilesChan <- &validation.InvalidFile{
		Path:             invalidPath,
		ValidationErrors: []error{},
	}
	close(invalidFilesChan)
	return invalidFilesChan
}
