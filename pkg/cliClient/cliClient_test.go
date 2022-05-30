package cliClient

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/datreeio/datree/pkg/networkValidator"

	"github.com/datreeio/datree/pkg/ciContext"

	"github.com/datreeio/datree/pkg/fileReader"

	"github.com/datreeio/datree/bl/files"

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

type RequestEvaluationPrerunDataTestCase struct {
	name string
	args struct {
		token   string
		offline string
	}
	mock struct {
		response struct {
			status int
			body   *EvaluationPrerunDataResponse
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
		response    *EvaluationPrerunDataResponse
	}
}

type SendEvaluationResultTestCase struct {
	name string
	args struct {
		evaluationRequestData *EvaluationResultRequest
		offline               string
	}
	mock struct {
		response struct {
			status int
			body   *SendEvaluationResultsResponse
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
		response    *SendEvaluationResultsResponse
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

func TestRequestEvaluationPrerunDataSuccess(t *testing.T) {
	tests := []*RequestEvaluationPrerunDataTestCase{
		test_requestEvaluationPrerunData_success(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, tt.mock.response.error).Once()

			networkValidator := networkValidator.NewNetworkValidator()
			client := &CliClient{
				baseUrl:          "http://cli-service.test.io",
				httpClient:       &httpClientMock,
				networkValidator: networkValidator,
			}

			policyCheckData, _ := client.RequestEvaluationPrerunData(tt.args.token, false)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, policyCheckData)
		})
	}
}

func TestRequestEvaluationPrerunDataFail(t *testing.T) {
	tests := []*RequestEvaluationPrerunDataTestCase{
		test_requestEvaluationPrerunData_error(),
		test_requestEvaluationPrerunData_network_error("fail", errors.New("Failed since internet connection refused, you can use the following command to set your config to run offline:\ndatree config set offline local")),
		test_requestEvaluationPrerunData_network_error("local", nil),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, tt.mock.response.error).Once()

			validator := networkValidator.NewNetworkValidator()
			validator.SetOfflineMode(tt.args.offline)
			client := &CliClient{
				baseUrl:          "http://cli-service.test.io",
				httpClient:       &httpClientMock,
				networkValidator: validator,
			}

			preRunDataResp, err := client.RequestEvaluationPrerunData(tt.args.token, false)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, preRunDataResp)
			assert.Equal(t, tt.expected.responseErr, err)
		})
	}
}

func TestSendEvaluationResult(t *testing.T) {
	tests := []*SendEvaluationResultTestCase{
		test_sendEvaluationResult_success(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, nil).Once()

			networkValidator := networkValidator.NewNetworkValidator()
			client := &CliClient{
				baseUrl:          "http://cli-service.test.io",
				httpClient:       &httpClientMock,
				networkValidator: networkValidator,
			}

			sendEvaluationResultsResponse, _ := client.SendEvaluationResult(tt.args.evaluationRequestData)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, sendEvaluationResultsResponse)
		})
	}
}

func TestSendEvaluationResultFail(t *testing.T) {
	sendEvalResultResp := &SendEvaluationResultsResponse{}

	tests := []*SendEvaluationResultTestCase{
		test_sendEvaluationResult_network_error("fail", errors.New("Failed since internet connection refused, you can use the following command to set your config to run offline:\ndatree config set offline local"), &SendEvaluationResultsResponse{}),
		test_sendEvaluationResult_network_error("local", nil, sendEvalResultResp),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, tt.mock.response.error).Once()

			validator := networkValidator.NewNetworkValidator()
			validator.SetOfflineMode(tt.args.offline)
			client := &CliClient{
				baseUrl:          "http://cli-service.test.io",
				httpClient:       &httpClientMock,
				networkValidator: validator,
			}

			sendEvaluationResultsResponse, err := client.SendEvaluationResult(tt.args.evaluationRequestData)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, sendEvaluationResultsResponse)
			assert.Equal(t, tt.expected.responseErr, err)
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

func TestAddFlags(t *testing.T) {
	httpClientMock := mockHTTPClient{}
	flagsToAdd := map[string]interface{}{"flag-key": "flag-value", "flag-key-bool": true, "bad-flag-key": []int{1, 3}}
	cliClient := &CliClient{
		baseUrl:      "http://cli-service.test.io",
		httpClient:   &httpClientMock,
		flagsHeaders: make(map[string]string),
	}

	cliClient.AddFlags(flagsToAdd)
	assert.Equal(t, "flag-value", cliClient.flagsHeaders["x-cli-flags-flag-key"])
	assert.Equal(t, "true", cliClient.flagsHeaders["x-cli-flags-flag-key-bool"])
	assert.Equal(t, "", cliClient.flagsHeaders["x-cli-flags-bad-flag-key"])
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

func mockGetPreRunData() *EvaluationPrerunDataResponse {
	const policiesJsonPath = "../../internal/fixtures/policyAsCode/prerun.json"

	fileReader := fileReader.CreateFileReader(nil)
	policiesJsonStr, err := fileReader.ReadFileContent(policiesJsonPath)

	if err != nil {
		panic(err)
	}

	policiesJsonRawData := []byte(policiesJsonStr)

	var policiesJson *EvaluationPrerunDataResponse
	err = json.Unmarshal(policiesJsonRawData, &policiesJson)

	if err != nil {
		panic(err)
	}
	return policiesJson
}

func test_requestEvaluationPrerunData_success() *RequestEvaluationPrerunDataTestCase {
	preRunData := mockGetPreRunData()

	return &RequestEvaluationPrerunDataTestCase{
		name: "success - get prerun data for evaluation",
		args: struct {
			token   string
			offline string
		}{
			token: "internal_test_token",
		},
		mock: struct {
			response struct {
				status int
				body   *EvaluationPrerunDataResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *EvaluationPrerunDataResponse
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
			response    *EvaluationPrerunDataResponse
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodGet,
				uri:     "/cli/evaluation/tokens/internal_test_token/prerun?isCi=false",
				body:    nil,
				headers: nil,
			},
			response: preRunData,
		},
	}
}

func test_requestEvaluationPrerunData_error() *RequestEvaluationPrerunDataTestCase {
	return &RequestEvaluationPrerunDataTestCase{
		name: "fail - get prerun data for evaluation",
		args: struct {
			token   string
			offline string
		}{
			token: "internal_test_token",
		},
		mock: struct {
			response struct {
				status int
				body   *EvaluationPrerunDataResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *EvaluationPrerunDataResponse
				error  error
			}{
				status: http.StatusBadRequest,
				body:   nil,
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
			response    *EvaluationPrerunDataResponse
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodGet,
				uri:     "/cli/evaluation/tokens/internal_test_token/prerun?isCi=false",
				body:    nil,
				headers: nil,
			},
			responseErr: errors.New("error from cli-service"),
			response:    &EvaluationPrerunDataResponse{},
		},
	}
}

func test_requestEvaluationPrerunData_network_error(offlineMode string, expectedResponseErr error) *RequestEvaluationPrerunDataTestCase {
	emptyRes := &EvaluationPrerunDataResponse{}
	if offlineMode == "local" {
		emptyRes.IsPolicyAsCodeMode = true
	}

	return &RequestEvaluationPrerunDataTestCase{
		name: "fail - get prerun data for evaluation",
		args: struct {
			token   string
			offline string
		}{
			token:   "internal_test_token",
			offline: offlineMode,
		},
		mock: struct {
			response struct {
				status int
				body   *EvaluationPrerunDataResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *EvaluationPrerunDataResponse
				error  error
			}{
				status: http.StatusInternalServerError,
				body:   nil,
				error:  errors.New("network error"),
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
			response    *EvaluationPrerunDataResponse
		}{
			request: struct {
				method  string
				uri     string
				body    interface{}
				headers map[string]string
			}{
				method:  http.MethodGet,
				uri:     "/cli/evaluation/tokens/internal_test_token/prerun?isCi=false",
				body:    nil,
				headers: nil,
			},
			responseErr: expectedResponseErr,
			response:    emptyRes,
		},
	}
}

func test_sendEvaluationResult_success() *SendEvaluationResultTestCase {
	body := &EvaluationResultRequest{
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
	return &SendEvaluationResultTestCase{
		name: "success - send local evaluation result to server",
		args: struct {
			evaluationRequestData *EvaluationResultRequest
			offline               string
		}{
			evaluationRequestData: body,
		},
		mock: struct {
			response struct {
				status int
				body   *SendEvaluationResultsResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *SendEvaluationResultsResponse
				error  error
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
			responseErr error
			response    *SendEvaluationResultsResponse
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

func test_sendEvaluationResult_network_error(offlineMode string, expectedResponseErr error, sendEvalResultResp *SendEvaluationResultsResponse) *SendEvaluationResultTestCase {
	body := &EvaluationResultRequest{
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
	return &SendEvaluationResultTestCase{
		name: "fail - send local evaluation result to server",
		args: struct {
			evaluationRequestData *EvaluationResultRequest
			offline               string
		}{
			evaluationRequestData: body,
			offline:               offlineMode,
		},
		mock: struct {
			response struct {
				status int
				body   *SendEvaluationResultsResponse
				error  error
			}
		}{
			response: struct {
				status int
				body   *SendEvaluationResultsResponse
				error  error
			}{
				status: http.StatusInternalServerError,
				body:   nil,
				error:  errors.New("network error"),
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
			response    *SendEvaluationResultsResponse
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
			responseErr: expectedResponseErr,
			response:    sendEvalResultResp,
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
