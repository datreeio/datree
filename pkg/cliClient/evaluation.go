package cliClient

import (
	"encoding/json"
	"net/http"

	"github.com/datreeio/datree/pkg/extractor"
)

type Metadata struct {
	CliVersion      string `json:"cliVersion"`
	Os              string `json:"os"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
}

type CreateEvaluationRequest struct {
	CliId      string    `json:"cliId"`
	Metadata   *Metadata `json:"metadata"`
	K8sVersion *string   `json:"k8sVersion"`
	PolicyName string    `json:"policyName"`
}

type CreateEvaluationResponse struct {
	EvaluationId  int    `json:"evaluationId"`
	K8sVersion    string `json:"k8sVersion"`
	RulesCount    int    `json:"rulesCount"`
	PolicyName    string `json:"policyName"`
	PromptMessage string `json:"promptMessage"`
}

func (c *CliClient) CreateEvaluation(request *CreateEvaluationRequest) (*CreateEvaluationResponse, error) {
	httpRes, err := c.httpClient.Request(http.MethodPost, "/cli/evaluation/create", request, nil)
	if err != nil {
		return nil, err
	}

	var res = &CreateEvaluationResponse{}
	err = json.Unmarshal(httpRes.Body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type Match struct {
	FileName     string `json:"fileName"`
	Path         string `json:"path"`
	Value        string `json:"value"`
	MetadataName string `json:"metadataName"`
	Kind         string `json:"kind"`
}

type EvaluationResult struct {
	Passed  bool `json:"passed"`
	Results struct {
		Matches    []*Match `json:"matches"`
		Mismatches []*Match `json:"mismatches"`
	} `json:"results"`
	Rule struct {
		Name           string     `json:"name"`
		FailSuggestion string     `json:"failSuggestion"`
		Origin         RuleOrigin `json:"origin"`
	} `json:"rule"`
}

type RuleOrigin struct {
	DefaultRuleId *int     `json:"defaultRuleId,omitempty"`
	CustomRuleId  *int     `json:"customRuleId,omitempty"`
	Type          RuleType `json:"type"`
}

type RuleType string

const (
	Default RuleType = "default"
	Custom           = "custom"
)

type EvaluationResponse struct {
	Results []*EvaluationResult `json:"results"`
	Status  string              `json:"status"`
}

type EvaluationRequest struct {
	EvaluationId int                             `json:"evaluationId"`
	Files        []*extractor.FileConfigurations `json:"files"`
}

func (c *CliClient) RequestEvaluation(request *EvaluationRequest) (*EvaluationResponse, error) {
	res, err := c.httpClient.Request(http.MethodPost, "/cli/evaluate", request, nil)
	if err != nil {
		return &EvaluationResponse{}, err
	}

	var evaluationResponse = &EvaluationResponse{}
	err = json.Unmarshal(res.Body, &evaluationResponse)
	if err != nil {
		return &EvaluationResponse{}, err
	}

	return evaluationResponse, nil
}

type UpdateEvaluationValidationRequest struct {
	EvaluationId   int       `json:"evaluationId"`
	InvalidFiles   []*string `json:"failedFiles"`
	StopEvaluation bool      `json:"stopEvaluation"`
}

func (c *CliClient) SendFailedYamlValidation(request *UpdateEvaluationValidationRequest) error {
	_, err := c.httpClient.Request(http.MethodPost, "/cli/evaluation/validation/yaml", request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *CliClient) SendFailedK8sValidation(request *UpdateEvaluationValidationRequest) error {
	_, err := c.httpClient.Request(http.MethodPost, "/cli/evaluation/validation/k8s", request, nil)
	if err != nil {
		return err
	}

	return nil
}
