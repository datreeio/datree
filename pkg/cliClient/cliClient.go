package cliClient

import (
	"encoding/json"
	"net/http"

	"github.com/datreeio/datree/pkg/httpClient"
	extractor "github.com/datreeio/datree/pkg/propertiesExtractor"
	"github.com/shirou/gopsutil/host"
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

type EvaluationRequest struct {
	CliId    string `json:"cliId"`
	Pattern  string `json:"pattern"`
	Metadata struct {
		CliVersion      string `json:"cliVersion"`
		Os              string `json:"os"`
		PlatformVersion string `json:"platformVersion"`
		KernelVersion   string `json:"kernelVersion"`
	} `json:"metadata"`
	Files []extractor.FileProperties `json:"files"`
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

func (c *CliClient) RequestEvaluation(pattern string, files []*extractor.FileProperties, cliId string) (EvaluationResponse, error) {
	evaluationRequest, err := c.createEvaluationRequest(pattern, files, cliId)
	if err != nil {
		return EvaluationResponse{}, err
	}
	res, err := c.httpClient.Request(http.MethodPost, "/cli/evaluate", evaluationRequest, nil)
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

func (c *CliClient) createEvaluationRequest(pattern string, files []*extractor.FileProperties, cliId string) (EvaluationRequest, error) {
	var filesProperties []extractor.FileProperties

	for _, file := range files {
		filesProperties = append(filesProperties, *file)
	}

	osInfo, err := host.Info()
	if err != nil {
		return EvaluationRequest{}, err
	}
	evaluationRequest := EvaluationRequest{
		CliId:   cliId,
		Pattern: pattern,
		Metadata: struct {
			CliVersion      string "json:\"cliVersion\""
			Os              string "json:\"os\""
			PlatformVersion string "json:\"platformVersion\""
			KernelVersion   string "json:\"kernelVersion\""
		}{
			CliVersion:      "0.0.1",
			Os:              osInfo.OS,
			PlatformVersion: osInfo.PlatformVersion,
			KernelVersion:   osInfo.KernelVersion,
		},
		Files: filesProperties,
	}
	return evaluationRequest, nil
}
