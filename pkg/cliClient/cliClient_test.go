package cliClient

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

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

func (c *mockHTTPClient) name()  {

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

			actualEvaluationId, _ := client.CreateEvaluation(tt.args.createEvaluationRequest)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response.EvaluationId, actualEvaluationId)

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
				baseUrl:    "http://cli-service.test.io",
				timeoutClient: &httpClientMock,
			}

			res, _ := client.GetVersionMessage(tt.args.cliVersion, 1000)
			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, tt.expected.response, res)

		})
	}
}

func readMock(path string) ([]extractor.K8sConfiguration, error) {
	var configurations []extractor.K8sConfiguration

	absPath, _ := filepath.Abs(path)
	content, err := ioutil.ReadFile(absPath)

	if err != nil {
		return []extractor.K8sConfiguration{}, err
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

func castPropertiesMock(fileName string, path string) []*extractor.FileConfiguration {
	configurations, _ := readMock(path)

	properties := []*extractor.FileConfiguration{
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
	return &CreateEvaluationTestCase{
		name: "success - create evaluation",
		args: struct {
			createEvaluationRequest *CreateEvaluationRequest
		}{
			createEvaluationRequest: &CreateEvaluationRequest{
				CliId: "cli_id",
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
					CliId: "cli_id",
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
			},
		},
	}
}
