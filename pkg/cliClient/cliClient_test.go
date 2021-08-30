package cliClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/datreeio/datree/bl/files"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

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

type RequestEvaluationTestCase struct {
	name string
	args struct {
		evaluationRequest *EvaluationRequest
	}
	mock struct {
		response struct {
			status int
			body   *EvaluationResponse
		}
	}
	expected struct {
		request struct {
			method  string
			uri     string
			body    *EvaluationRequest
			headers map[string]string
		}
		response *EvaluationResponse
	}
}

type CreateEvaluationTestCase struct {
	name string
	args struct {
		createEvaluationRequest *CreateEvaluationRequest
	}
	mock struct {
		response struct {
			status int
			body   *CreateEvaluationResponse
		}
	}
	expected struct {
		request struct {
			method  string
			uri     string
			body    *CreateEvaluationRequest
			headers map[string]string
		}
		response *CreateEvaluationResponse
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
		cliId                 string
	}
	mockResponse struct {
		status int
		body   struct {
			message string
		}
		error error
	}

	expected struct {
		request struct {
			method  string
			uri     string
			body    files.UnknownStruct
			headers map[string]string
		}
		response error
	}
}

func TestRequestEvaluation(t *testing.T) {
	tests := []*RequestEvaluationTestCase{
		test_requestEvaluation_success(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, nil)

			client := &CliClient{
				baseUrl:    "http://cli-service.test.io",
				httpClient: &httpClientMock,
			}

			res, _ := client.RequestEvaluation(tt.args.evaluationRequest)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, res)

		})
	}
}

func TestCreateRequestEvaluation(t *testing.T) {
	tests := []*CreateEvaluationTestCase{
		test_createEvaluation_success(),
	}

	httpClientMock := mockHTTPClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.mock.response.body)
			mockedHTTPResponse := httpClient.Response{StatusCode: tt.mock.response.status, Body: body}
			httpClientMock.On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockedHTTPResponse, nil)

			client := &CliClient{
				baseUrl:    "http://cli-service.test.io",
				httpClient: &httpClientMock,
			}

			res, _ := client.CreateEvaluation(tt.args.createEvaluationRequest)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, res)
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

			actualResponse := client.PublishPolicies(tt.args.policiesConfiguration, tt.args.cliId)
			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, actualResponse)
		})
	}
}

func readMock(path string) ([]extractor.Configuration, error) {
	var configurations []extractor.Configuration

	absPath, _ := filepath.Abs(path)
	content, err := ioutil.ReadFile(absPath)

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

func test_requestEvaluation_success() *RequestEvaluationTestCase {
	return &RequestEvaluationTestCase{
		name: "success - request evaluation",
		args: struct {
			evaluationRequest *EvaluationRequest
		}{
			evaluationRequest: &EvaluationRequest{
				EvaluationId: 321,
				Files:        castPropertiesMock("service_mock", "mocks/service_mock.yaml"),
			},
		},
		mock: struct {
			response struct {
				status int
				body   *EvaluationResponse
			}
		}{
			response: struct {
				status int
				body   *EvaluationResponse
			}{
				status: http.StatusOK,
				body:   &EvaluationResponse{},
			},
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    *EvaluationRequest
				headers map[string]string
			}
			response *EvaluationResponse
		}{
			request: struct {
				method  string
				uri     string
				body    *EvaluationRequest
				headers map[string]string
			}{
				method: http.MethodPost,
				uri:    "/cli/evaluate",
				body: &EvaluationRequest{
					EvaluationId: 321,
					Files:        castPropertiesMock("service_mock", "mocks/service_mock.yaml"),
				},
				headers: nil,
			},
			response: &EvaluationResponse{},
		},
	}
}

func test_createEvaluation_success() *CreateEvaluationTestCase {
	k8sVersion := "1.18.0"

	return &CreateEvaluationTestCase{
		name: "success - create evaluation",
		args: struct {
			createEvaluationRequest *CreateEvaluationRequest
		}{
			createEvaluationRequest: &CreateEvaluationRequest{
				K8sVersion: &k8sVersion,
				CliId:      "cli_id",
				PolicyName: "Default",
				Metadata: &Metadata{
					CliVersion:      "0.0.1",
					Os:              "darwin",
					PlatformVersion: "1.2.3",
					KernelVersion:   "4.5.6",
				},
			},
		},
		mock: struct {
			response struct {
				status int
				body   *CreateEvaluationResponse
			}
		}{
			response: struct {
				status int
				body   *CreateEvaluationResponse
			}{
				status: http.StatusOK,
				body: &CreateEvaluationResponse{
					EvaluationId: 123,
					K8sVersion:   k8sVersion,
				},
			},
		},
		expected: struct {
			request struct {
				method  string
				uri     string
				body    *CreateEvaluationRequest
				headers map[string]string
			}
			response *CreateEvaluationResponse
		}{
			request: struct {
				method  string
				uri     string
				body    *CreateEvaluationRequest
				headers map[string]string
			}{
				method: http.MethodPost,
				uri:    "/cli/evaluation/create",
				body: &CreateEvaluationRequest{
					K8sVersion: &k8sVersion,
					CliId:      "cli_id",
					PolicyName: "Default",
					Metadata: &Metadata{
						CliVersion:      "0.0.1",
						Os:              "darwin",
						PlatformVersion: "1.2.3",
						KernelVersion:   "4.5.6",
					},
				},
				headers: nil,
			},
			response: &CreateEvaluationResponse{
				EvaluationId: 123,
				K8sVersion:   k8sVersion,
			},
		},
	}
}

func test_publishPolicies_success() *PublishPoliciesTestCase {
	expectedPublishHeaders := map[string]string{"x-cli-id": "cli_id"}

	requestPoliciesConfigurationArg := files.UnknownStruct{}
	return &PublishPoliciesTestCase{
		name: "success - publish policies",
		args: struct {
			policiesConfiguration files.UnknownStruct
			cliId                 string
		}{
			policiesConfiguration: requestPoliciesConfigurationArg,
			cliId:                 "cli_id",
		},
		mockResponse: struct {
			status int
			body   struct {
				message string
			}
			error error
		}{
			status: http.StatusCreated,
			body: struct {
				message string
			}{
				message: "",
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
			response error
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
			response: nil,
		},
	}
}

func test_publishPolicies_schemaError() *PublishPoliciesTestCase {
	expectedPublishHeaders := map[string]string{"x-cli-id": "cli_id"}

	requestPoliciesConfigurationArg := files.UnknownStruct{}
	return &PublishPoliciesTestCase{
		name: "schema error - publish policies",
		args: struct {
			policiesConfiguration files.UnknownStruct
			cliId                 string
		}{
			policiesConfiguration: requestPoliciesConfigurationArg,
			cliId:                 "cli_id",
		},
		mockResponse: struct {
			status int
			body   struct {
				message string
			}
			error error
		}{
			status: http.StatusBadRequest,
			body: struct {
				message string
			}{
				message: "",
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
			response error
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
			response: errors.New("error from cli-service"),
		},
	}
}
