package cliClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/datreeio/datree/pkg/ciContext"

	"github.com/datreeio/datree/pkg/fileReader"

	"github.com/datreeio/datree/bl/files"

	"gopkg.in/yaml.v3"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/httpClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHTTPClient struct {
	mock.Mock
}

func (c *mockHTTPClient) Request(method string, resourceURI string, body interface{}, headers map[string]string) (httpClient.Response, error) {
	args := c.Called(method, resourceURI, body, headers)

	return args.Get(0).(httpClient.Response), args.Error(1)
}

func (c *mockHTTPClient) name() {

}

type RequestPrerunDataForEvaluationTestCase struct {
	name string
	args struct {
		token string
	}
	mock struct {
		response struct {
			status int
			body   *PrerunDataForEvaluationResponse
			error  error
		}
	}
	expected struct {
		request struct {
			method  string
			uri     string
			body    interface{}
			headers map[string]string
		}
		responseErr error
		response    *PrerunDataForEvaluationResponse
	}
}

type SendLocalEvaluationResultTestCase struct {
	name string
	args struct {
		localEvaluationRequestData *LocalEvaluationResultRequest
	}
	mock struct {
		response struct {
			status int
			body   *SendEvaluationResultsResponse
		}
	}
	expected struct {
		request struct {
			method  string
			uri     string
			body    interface{}
			headers map[string]string
		}
		response *SendEvaluationResultsResponse
	}
}

type GetVersionMessageTestCase struct {
	name string
	args struct {
		cliVersion string
	}
	mock struct {
		response struct {
			status int
			body   *VersionMessage
		}
	}
	expected struct {
		request struct {
			method  string
			uri     string
			body    interface{}
			headers map[string]string
		}
		response *VersionMessage
	}
}

type PublishPoliciesTestCase struct {
	name string
	args struct {
		policiesConfiguration files.UnknownStruct
		token                 string
	}
	mockResponse struct {
		status int
		body   PublishFailedResponse
		error  error
	}

	expected struct {
		request struct {
			method  string
			uri     string
			body    files.UnknownStruct
			headers map[string]string
		}
		responseErr           error
		publishFailedResponse *PublishFailedResponse
	}
}

func TestRequestPrerunDataForEvaluationSuccess(t *testing.T) {
	tests := []*RequestPrerunDataForEvaluationTestCase{
		test_requestPrerunDataForEvaluation_success(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, tt.mock.response.error).Once()

			client := &CliClient{
				baseUrl:    "http://cli-service.test.io",
				httpClient: &httpClientMock,
			}

			prerunDataForEvaluation, _, _ := client.RequestPrerunDataForEvaluation(tt.args.token)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, prerunDataForEvaluation)
		})
	}
}

func TestRequestPrerunDataForEvaluationFail(t *testing.T) {
	tests := []*RequestPrerunDataForEvaluationTestCase{
		test_requestPrerunDataForEvaluation_error(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, tt.mock.response.error).Once()

			client := &CliClient{
				baseUrl:    "http://cli-service.test.io",
				httpClient: &httpClientMock,
			}

			_, _, err := client.RequestPrerunDataForEvaluation(tt.args.token)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.responseErr, err)
		})
	}
}

func TestSendLocalEvaluationResult(t *testing.T) {
	tests := []*SendLocalEvaluationResultTestCase{
		test_sendLocalEvaluationResult_success(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, nil).Once()

			client := &CliClient{
				baseUrl:    "http://cli-service.test.io",
				httpClient: &httpClientMock,
			}

			sendEvaluationResultsResponse, _ := client.SendLocalEvaluationResult(tt.args.localEvaluationRequestData)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, sendEvaluationResultsResponse)
		})
	}
}

func TestGetVersionMessage(t *testing.T) {
	tests := []*GetVersionMessageTestCase{
		test_getVersionMessage_success(),
	}
	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, nil)

			client := &CliClient{
				baseUrl:       "http://cli-service.test.io",
				timeoutClient: &httpClientMock,
			}

			res, _ := client.GetVersionMessage(tt.args.cliVersion, 1000)
			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, res)

		})
	}
}

func TestPublishPolicies(t *testing.T) {
	tests := []*PublishPoliciesTestCase{
		test_publishPolicies_success(),
		test_publishPolicies_schemaError(),
	}
	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mockResponse.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mockResponse.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, tt.mockResponse.error).Once()

			client := &CliClient{
				baseUrl:    "http://cli-service.test.io",
				httpClient: &httpClientMock,
			}

			publishFailedResponse, err := client.PublishPolicies(tt.args.policiesConfiguration, tt.args.token)
			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.publishFailedResponse, publishFailedResponse)
			assert.Equal(t, tt.expected.responseErr, err)
		})
	}
}

func readMock(path string) ([]extractor.Configuration, error) {
	var configurations []extractor.Configuration

	absPath, _ := filepath.Abs(path)
	content, err := os.ReadFile(absPath)

	if err != nil {
		return []extractor.Configuration{}, err
	}

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(content))

	for {
		var doc = map[string]interface{}{}
		err = yamlDecoder.Decode(&doc)
		if err != nil {
			break
		}
		configurations = append(configurations, doc)
	}

	return configurations, nil
}

func castPropertiesMock(fileName string, path string) []*extractor.FileConfigurations {
	configurations, _ := readMock(path)

	properties := []*extractor.FileConfigurations{
		{
			FileName:       fileName,
			Configurations: configurations,
		}}

	return properties
}

func test_getVersionMessage_success() *GetVersionMessageTestCase {
	return &GetVersionMessageTestCase{
		name: "success - get version message",
		args: struct {
			cliVersion string
		}{
			cliVersion: "0.0.1",
		},
		mock: struct {
			response struct {
				status int
				body   *VersionMessage
			}
		}{
			response: struct {
				status int
				body   *VersionMessage
			}{
				status: http.StatusOK,
				body:   &VersionMessage{},
			},
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}
			response *VersionMessage
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodGet,
				uri:     "/cli/messages/versions/0.0.1",
				body:    nil,
				headers: nil,
			},
			response: &VersionMessage{},
		},
	}
}

func mockGetPreRunData() *PrerunDataForEvaluationResponse {
	err := os.Chdir("../../")

	const policiesJsonPath = "internal/fixtures/policyAsCode/policies.json"

	fileReader := fileReader.CreateFileReader(nil)
	policiesJsonStr, err := fileReader.ReadFileContent(policiesJsonPath)

	if err != nil {
		fmt.Errorf("can't read policies json")
	}

	policiesJsonRawData := []byte(policiesJsonStr)

	var policiesJson *PrerunDataForEvaluationResponse
	err = json.Unmarshal(policiesJsonRawData, &policiesJson)

	if err != nil {
		fmt.Errorf("can't marshel policies json")
	}
	return policiesJson
}

func test_requestPrerunDataForEvaluation_success() *RequestPrerunDataForEvaluationTestCase {
	preRunData := mockGetPreRunData()

	return &RequestPrerunDataForEvaluationTestCase{
		name: "success - get prerun data for evaluation",
		args: struct {
			token string
		}{
			token: "internal_test_token",
		},
		mock: struct {
			response struct {
				status int
				body   *PrerunDataForEvaluationResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *PrerunDataForEvaluationResponse
				error  error
			}{
				status: http.StatusOK,
				body:   preRunData,
			},
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}
			responseErr error
			response    *PrerunDataForEvaluationResponse
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodGet,
				uri:     "/cli/evaluation/tokens/internal_test_token/prerun",
				body:    nil,
				headers: nil,
			},
			response: preRunData,
		},
	}
}

func test_requestPrerunDataForEvaluation_error() *RequestPrerunDataForEvaluationTestCase {
	preRunData := mockGetPreRunData()

	return &RequestPrerunDataForEvaluationTestCase{
		name: "success - get prerun data for evaluation",
		args: struct {
			token string
		}{
			token: "internal_test_token",
		},
		mock: struct {
			response struct {
				status int
				body   *PrerunDataForEvaluationResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *PrerunDataForEvaluationResponse
				error  error
			}{
				status: http.StatusOK,
				body:   preRunData,
				error:  errors.New("error from cli-service"),
			},
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}
			responseErr error
			response    *PrerunDataForEvaluationResponse
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodGet,
				uri:     "/cli/evaluation/tokens/internal_test_token/prerun",
				body:    nil,
				headers: nil,
			},
			responseErr: errors.New("error from cli-service"),
			response:    preRunData,
		},
	}
}

func test_sendLocalEvaluationResult_success() *SendLocalEvaluationResultTestCase {
	body := &LocalEvaluationResultRequest{
		ClientId: "internal_cliId_test",
		Token:    "internal_cliId_test",
		Metadata: &Metadata{
			CliVersion:      "0.0.01",
			Os:              "darwin",
			PlatformVersion: "1.2.3",
			KernelVersion:   "4.5.6",
			CIContext: &ciContext.CIContext{
				IsCI: true,
				CIMetadata: &ciContext.CIMetadata{
					CIEnvValue:       "travis",
					ShouldHideEmojis: false,
				},
			},
		},
		K8sVersion:         "1.18.0",
		PolicyName:         "Default",
		FailedYamlFiles:    []string{},
		FailedK8sFiles:     []string{},
		AllExecutedRules:   []RuleData{},
		AllEvaluatedFiles:  []FileData{},
		PolicyCheckResults: nil,
	}
	return &SendLocalEvaluationResultTestCase{
		name: "success - send local evaluation result to server",
		args: struct {
			localEvaluationRequestData *LocalEvaluationResultRequest
		}{
			localEvaluationRequestData: body,
		},
		mock: struct {
			response struct {
				status int
				body   *SendEvaluationResultsResponse
			}
		}{
			response: struct {
				status int
				body   *SendEvaluationResultsResponse
			}{
				status: http.StatusOK,
				body: &SendEvaluationResultsResponse{
					EvaluationId:  1234,
					PromptMessage: "",
				},
			},
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}
			response *SendEvaluationResultsResponse
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodPost,
				uri:     "/cli/evaluation/result",
				body:    body,
				headers: nil,
			},
			response: &SendEvaluationResultsResponse{
				EvaluationId:  1234,
				PromptMessage: "",
			},
		},
	}
}

func test_publishPolicies_success() *PublishPoliciesTestCase {
	expectedPublishHeaders := map[string]string{"x-datree-token": "token"}

	requestPoliciesConfigurationArg := files.UnknownStruct{}
	return &PublishPoliciesTestCase{
		name: "success - publish policies",
		args: struct {
			policiesConfiguration files.UnknownStruct
			token                 string
		}{
			policiesConfiguration: requestPoliciesConfigurationArg,
			token:                 "token",
		},
		mockResponse: struct {
			status int
			body   PublishFailedResponse
			error  error
		}{
			status: http.StatusCreated,
			body: PublishFailedResponse{
				Code:    "mocked code",
				Message: "error from cli-service",
				Payload: []string{"error from cli-service"},
			},
			error: nil,
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    files.UnknownStruct
				headers map[string]string
			}
			responseErr           error
			publishFailedResponse *PublishFailedResponse
		}{
			request: struct {
				method  string
				uri     string
				body    files.UnknownStruct
				headers map[string]string
			}{
				method:  http.MethodPut,
				uri:     "/cli/policy/publish",
				body:    requestPoliciesConfigurationArg,
				headers: expectedPublishHeaders,
			},
			responseErr:           nil,
			publishFailedResponse: nil,
		},
	}
}

func test_publishPolicies_schemaError() *PublishPoliciesTestCase {
	expectedPublishHeaders := map[string]string{"x-datree-token": "token"}

	requestPoliciesConfigurationArg := files.UnknownStruct{}
	return &PublishPoliciesTestCase{
		name: "schema error - publish policies",
		args: struct {
			policiesConfiguration files.UnknownStruct
			token                 string
		}{
			policiesConfiguration: requestPoliciesConfigurationArg,
			token:                 "token",
		},
		mockResponse: struct {
			status int
			body   PublishFailedResponse
			error  error
		}{
			status: http.StatusBadRequest,
			body: PublishFailedResponse{
				Code:    "mocked code",
				Message: "error from cli-service",
				Payload: []string{"error from cli-service"},
			},
			error: errors.New("error from cli-service"),
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    files.UnknownStruct
				headers map[string]string
			}
			responseErr           error
			publishFailedResponse *PublishFailedResponse
		}{
			request: struct {
				method  string
				uri     string
				body    files.UnknownStruct
				headers map[string]string
			}{
				method:  http.MethodPut,
				uri:     "/cli/policy/publish",
				body:    requestPoliciesConfigurationArg,
				headers: expectedPublishHeaders,
			},
			responseErr: errors.New("error from cli-service"),
			publishFailedResponse: &PublishFailedResponse{
				Code:    "mocked code",
				Message: "error from cli-service",
				Payload: []string{"error from cli-service"},
			},
		},
	}
}
