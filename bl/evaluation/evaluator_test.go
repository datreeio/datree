package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/fileReader"

	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCliClient struct {
	mock.Mock
}

func (m *mockCliClient) RequestPrerunDataForEvaluation(token string) (*cliClient.PrerunDataForEvaluationResponse, error) {
	args := m.Called(token)
	return args.Get(0).(*cliClient.PrerunDataForEvaluationResponse), args.Error(1)
}

func (m *mockCliClient) SendLocalEvaluationResult(localEvaluationResultRequest *cliClient.LocalEvaluationResultRequest) (*cliClient.SendEvaluationResultsResponse, error) {
	args := m.Called(localEvaluationResultRequest)
	return args.Get(0).(*cliClient.SendEvaluationResultsResponse), args.Error(1)
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
func TestSendLocalEvaluationResult(t *testing.T) {
	t.Run("SendLocalEvaluationResult should succeed", func(t *testing.T) {
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
		promptMessage := ""
		policyName := "Default"
		ciContext := &ciContext.CIContext{
			IsCI: true,
			CIMetadata: &ciContext.CIMetadata{
				CIEnvValue:       "travis",
				ShouldHideEmojis: false,
			},
		}

		mockedCliClient.On("SendLocalEvaluationResult", mock.Anything).Return(&cliClient.SendEvaluationResultsResponse{EvaluationId: 1, PromptMessage: promptMessage}, nil)

		expectedSendEvaluationResultsResponse := &cliClient.SendEvaluationResultsResponse{EvaluationId: 1, PromptMessage: promptMessage}

		localEvaluationRequestData := LocalEvaluationRequestData{
			CliId:              cliId,
			CliVersion:         cliVersion,
			K8sVersion:         k8sVersion,
			PolicyName:         policyName,
			CiContext:          ciContext,
			RulesData:          []cliClient.RuleData{},
			FilesData:          []cliClient.FileData{},
			FailedYamlFiles:    []string{},
			FailedK8sFiles:     []string{},
			PolicyCheckResults: nil,
		}

		sendEvaluationResultsResponse, _ := evaluator.SendLocalEvaluationResult(localEvaluationRequestData)

		sendLocalEvaluationResultRequestData := &cliClient.LocalEvaluationResultRequest{
			K8sVersion: localEvaluationRequestData.K8sVersion,
			ClientId:   localEvaluationRequestData.CliId,
			Token:      localEvaluationRequestData.CliId,
			PolicyName: localEvaluationRequestData.PolicyName,
			Metadata: &cliClient.Metadata{
				CliVersion:      localEvaluationRequestData.CliVersion,
				Os:              evaluator.osInfo.OS,
				PlatformVersion: evaluator.osInfo.PlatformVersion,
				KernelVersion:   evaluator.osInfo.KernelVersion,
				CIContext:       localEvaluationRequestData.CiContext,
			},
			FailedYamlFiles:    localEvaluationRequestData.FailedYamlFiles,
			FailedK8sFiles:     localEvaluationRequestData.FailedK8sFiles,
			AllExecutedRules:   localEvaluationRequestData.RulesData,
			AllEvaluatedFiles:  localEvaluationRequestData.FilesData,
			PolicyCheckResults: localEvaluationRequestData.PolicyCheckResults,
		}
		mockedCliClient.AssertCalled(t, "SendLocalEvaluationResult", sendLocalEvaluationResultRequestData)
		assert.Equal(t, expectedSendEvaluationResultsResponse, sendEvaluationResultsResponse)

	})
}

func TestEvaluate(t *testing.T) {
	os.Chdir("../../")

	tests := []*evaluateTestCase{
		//request_evaluation_all_valid(),
		request_evaluation_all_invalid(),
	}

	prerunData := mockGetPreRunData()
	policy, _ := policy_factory.CreatePolicy(prerunData.PoliciesJson, "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedCliClient := &mockCliClient{}

			evaluator := &Evaluator{
				cliClient: mockedCliClient,
				osInfo:    tt.args.osInfo,
			}

			dataForEvaluation := DataForPolicyCheck{
				FilesConfigurations: tt.args.dataForEvaluation.FilesConfigurations,
				IsInteractiveMode:   tt.args.dataForEvaluation.IsInteractiveMode,
				PolicyName:          policy.Name,
				Policy:              policy,
			}

			policyCheckResultData, err := evaluator.Evaluate(dataForEvaluation)
			if err != nil {
				fmt.Errorf(err.Error())
			}

			if len(dataForEvaluation.FilesConfigurations) > 0 {
				assert.Equal(t, tt.expected.policyCheckResultData.FormattedResults.EvaluationResults.Summary, policyCheckResultData.FormattedResults.EvaluationResults.Summary)
				assert.Equal(t, tt.expected.policyCheckResultData.FormattedResults.EvaluationResults.FileNameRuleMapper, policyCheckResultData.FormattedResults.EvaluationResults.FileNameRuleMapper)
			} else {
				assert.Equal(t, tt.expected.policyCheckResultData.FormattedResults, policyCheckResultData.FormattedResults)
			}
		})
	}
}

type evaluateArgs struct {
	dataForEvaluation DataForPolicyCheck
	osInfo            *OSInfo
}

type evaluateExpected struct {
	policyCheckResultData PolicyCheckResultData
	err                   error
}

type evaluateTestCase struct {
	name     string
	args     *evaluateArgs
	mock     *evaluatorMock
	expected *evaluateExpected
}

func request_evaluation_all_valid() *evaluateTestCase {
	validFilePath := "internal/fixtures/kube/pass-all.yaml"

	prerunData := mockGetPreRunData()
	policy, _ := policy_factory.CreatePolicy(prerunData.PoliciesJson, "")

	return &evaluateTestCase{
		name: "should request validation without invalid files",
		args: &evaluateArgs{
			dataForEvaluation: DataForPolicyCheck{
				FilesConfigurations: newFilesConfigurations(validFilePath),
				IsInteractiveMode:   true,
				PolicyName:          "Default",
				Policy:              policy,
			},
			osInfo: &OSInfo{
				OS:              "darwin",
				PlatformVersion: "1.2.3",
				KernelVersion:   "4.5.6",
			},
		},
		mock: &evaluatorMock{
			cliClient: &cliClientMockTestCase{
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
			policyCheckResultData: PolicyCheckResultData{
				FormattedResults: FormattedResults{
					EvaluationResults: &EvaluationResults{
						FileNameRuleMapper: make(map[string]map[string]*Rule),
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
					NonInteractiveEvaluationResults: nil,
				},
			},
			err: nil,
		},
	}
}

func request_evaluation_all_invalid() *evaluateTestCase {
	prerunData := mockGetPreRunData()
	policy, _ := policy_factory.CreatePolicy(prerunData.PoliciesJson, "")

	return &evaluateTestCase{
		name: "should not request validation if there are no valid files",
		args: &evaluateArgs{
			dataForEvaluation: DataForPolicyCheck{
				FilesConfigurations: []*extractor.FileConfigurations{},
				IsInteractiveMode:   true,
				PolicyName:          "Default",
				Policy:              policy,
			},
			osInfo: &OSInfo{
				OS:              "darwin",
				PlatformVersion: "1.2.3",
				KernelVersion:   "4.5.6",
			},
		},
		mock: &evaluatorMock{
			cliClient: &cliClientMockTestCase{
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
			policyCheckResultData: PolicyCheckResultData{
				FormattedResults: FormattedResults{},
			},
			err: nil,
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

func mockGetPreRunData() *cliClient.PrerunDataForEvaluationResponse {
	const policiesJsonPath = "internal/fixtures/policyAsCode/policies.json"

	fileReader := fileReader.CreateFileReader(nil)
	policiesJsonStr, err := fileReader.ReadFileContent(policiesJsonPath)

	if err != nil {
		fmt.Errorf("can't read policies json")
	}

	policiesJsonRawData := []byte(policiesJsonStr)

	var policiesJson *cliClient.PrerunDataForEvaluationResponse
	err = json.Unmarshal(policiesJsonRawData, &policiesJson)

	if err != nil {
		fmt.Errorf("can't marshel policies json")
	}

	return policiesJson
}
