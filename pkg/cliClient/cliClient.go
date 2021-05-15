package cliClient

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/datreeio/datree/pkg/httpClient"
	extractor "github.com/datreeio/datree/pkg/propertiesExtractor"
)

type HTTPClient interface {
	Request(method string, resourceURI string, body interface{}, headers map[string]string) (httpClient.Response, error)
}
type CliClient struct {
	baseUrl    string
	httpClient HTTPClient
}

func NewCliClient(url string) *CliClient {
	httpClient := httpClient.NewClient(url, nil)
	return &CliClient{
		baseUrl:    url,
		httpClient: httpClient,
	}
}

type CreateEvaluationRequest struct {
	CliId    string   `json:"cliId"`
	Metadata Metadata `json:"metadata"`
}

type CreateEvaluationResponse struct {
	EvaluationId int `json:"evaluationId"`
}

func (c *CliClient) CreateEvaluation(request CreateEvaluationRequest) (int, error) {
	httpRes, err := c.httpClient.Request(http.MethodPost, "/cli/evaluate/create", request, nil)
	if err != nil {
		return 0, err
	}

	var res = &CreateEvaluationResponse{}
	err = json.Unmarshal(httpRes.Body, &res)
	if err != nil {
		return 0, err
	}

	return res.EvaluationId, nil
}

type Metadata struct {
	CliVersion      string `json:"cliVersion"`
	Os              string `json:"os"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
}

type EvaluationRequest struct {
	EvaluationId int                         `json:"evaluationId"`
	Files        []*extractor.FileProperties `json:"files"`
}

type Match struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
	Value    string `json:"value"`
}

type EvaluationResult struct {
	Passed  bool `json:"passed"`
	Results struct {
		Matches    []Match `json:"matches"`
		Mismatches []Match `json:"mismatches"`
	} `json:"results"`
	Rule struct {
		ID             int    `json:"defaultRuleId"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		FailSuggestion string `json:"failSuggestion"`
	} `json:"rule"`
}

type EvaluationResponse struct {
	Results []EvaluationResult `json:"results"`
	Status  string             `json:"status"`
}

func (c *CliClient) RequestEvaluation(request EvaluationRequest) (EvaluationResponse, error) {
	res, err := c.httpClient.Request(http.MethodPost, "/cli/evaluate", request, nil)
	if err != nil {
		return EvaluationResponse{}, err
	}

	var evaluationResponse = &EvaluationResponse{}
	err = json.Unmarshal(res.Body, &evaluationResponse)
	if err != nil {
		return EvaluationResponse{}, err
	}

	return *evaluationResponse, nil
}

type VersionMessage struct {
	CliVersion   string `json:"cliVersion"`
	MessageText  string `json:"messageText"`
	MessageColor string `json:"messageColor"`
}

type VersionMessageClient struct {
	baseUrl    string
	httpClient HTTPClient
}

func NewVersionMessageClient(url string) VersionMessageClient {
	httpClient := httpClient.NewClientTimeout(url, nil, 900*time.Millisecond)
	return VersionMessageClient{
		baseUrl:    url,
		httpClient: httpClient,
	}
}

func (c VersionMessageClient) GetVersionMessage(cliVersion string) (*VersionMessage, error) {
	res, err := c.httpClient.Request(http.MethodGet, "/cli/messages/versions/"+cliVersion, nil, nil)
	if err != nil {
		return nil, err
	}

	var response = &VersionMessage{}
	err = json.Unmarshal(res.Body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
