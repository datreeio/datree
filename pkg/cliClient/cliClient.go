package cliClient

import (
	"encoding/json"
	"net/http"

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

type Metadata struct {
	CliVersion      string `json:"cliVersion"`
	Os              string `json:"os"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
}

type EvaluationRequest struct {
	CliId    string                     `json:"cliId"`
	Metadata Metadata                   `json:"metadata"`
	Files    []extractor.FileProperties `json:"files"`
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
	EvaluationId int                `json:"evaluationId"`
	Results      []EvaluationResult `json:"results"`
	Status       string             `json:"status"`
}

type VersionMessage struct {
	CliVersion   string `json:"cliVersion"`
	MessageText  string `json:"messageText"`
	MessageColor string `json:"messageColor"`
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

func (c *CliClient) GetVersionMessage(cliVersion string) (*VersionMessage, error) {
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
