package cliClient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/datreeio/datree/pkg/httpClient"
	extractor "github.com/datreeio/datree/pkg/propertiesExtractor"
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

type RequestEvaluationTestCase struct {
	name string
	args struct {
		pattern    string
		cliId      string
		properties []*extractor.FileProperties
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

			res, _ := client.RequestEvaluation(tt.args.pattern, tt.args.properties, tt.args.cliId)

			httpClientMock.AssertCalled(t, "Request", tt.expected.request.method, tt.expected.request.uri, *tt.expected.request.body, tt.expected.request.headers)
			assert.Equal(t, *tt.expected.response, res)

		})
	}
}

func readMock(path string, data interface{}) error {
	absPath, _ := filepath.Abs(path)
	f, err := ioutil.ReadFile(absPath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(f, data)
	if err != nil {
		return err
	}
	return nil
}

func castPropertiesMock(fileName string, path string) []extractor.FileProperties {
	var fileProperties map[string]interface{}
	_ = readMock(path, &fileProperties)

	properties := []extractor.FileProperties{
		{
			FileName:       fileName,
			Configurations: []extractor.K8sConfiguration{fileProperties},
		}}

	return properties
}

func castPropertiesPointersMock(fileName string, path string) []*extractor.FileProperties {
	var filesProperties []*extractor.FileProperties
	props := castPropertiesMock("service_mock", "mocks/service_mock.yaml")
	for _, p := range props {
		filesProperties = append(filesProperties, &p)
	}

	return filesProperties

}

func test_requestEvaluation_success() *RequestEvaluationTestCase {
	return &RequestEvaluationTestCase{
		name: "success - request evaluation",
		args: struct {
			pattern    string
			cliId      string
			properties []*extractor.FileProperties
		}{
			pattern:    "pattern",
			properties: castPropertiesPointersMock("service_mock", "mocks/service_mock.yaml"),
			cliId:      "cli-id-test",
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
					CliId:   "cli-id-test",
					Pattern: "pattern",
					Files:   castPropertiesMock("service_mock", "mocks/service_mock.yaml"),
					Metadata: struct {
						CliVersion string "json:\"cliVersion\""
						Os         string "json:\"os\""
					}{
						CliVersion: "0.0.1",
						Os:         runtime.GOOS,
					},
				},
				headers: nil,
			},
			response: &EvaluationResponse{},
		},
	}
}
