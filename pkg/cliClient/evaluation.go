package cliClient

import (
	"encoding/json"
	"net/http"

	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/extractor"
)

type Metadata struct {
	CliVersion      string               `json:"cliVersion"`
	Os              string               `json:"os"`
	PlatformVersion string               `json:"platformVersion"`
	KernelVersion   string               `json:"kernelVersion"`
	CIContext       *ciContext.CIContext `json:"ciContext"`
}

type SendEvaluationResultsResponse struct {
	EvaluationId  int    `json:"evaluationId"`
	PromptMessage string `json:"promptMessage,omitempty"`
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
		Identifier     string     `json:"identifier"`
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

type CustomRule struct {
	Identifier              string                 `json:"identifier"`
	Name                    string                 `json:"name"`
	DefaultMessageOnFailure string                 `json:"defaultMessageOnFailure"`
	Schema                  map[string]interface{} `json:"schema"`
}

type Rule struct {
	Identifier       string `json:"identifier"`
	MessageOnFailure string `json:"messageOnFailure"`
}

type Policy struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault,omitempty"`
	Rules     []Rule `json:"rules"`
}

type EvaluationPrerunPolicies struct {
	ApiVersion  string        `json:"apiVersion"`
	CustomRules []*CustomRule `json:"customRules"`
	Policies    []*Policy     `json:"policies"`
}

type EvaluationPrerunDataResponse struct {
	PoliciesJson          *EvaluationPrerunPolicies `json:"policiesJson"`
	DefaultK8sVersion     string                    `json:"defaultK8sVersion"`
	AccountExists         bool                      `json:"accountExists"`
	RegistrationURL       string                    `json:"registrationURL"`
	PromptRegistrationURL string                    `json:"promptRegistrationURL"`
}

const badRequestStatusCode = 400

func (c *CliClient) RequestEvaluationPrerunData(tokenId string) (*EvaluationPrerunDataResponse, error) {
	res, err := c.httpClient.Request(http.MethodGet, "/cli/evaluation/tokens/"+tokenId+"/prerun", nil, nil)

	if err != nil && (res.StatusCode >= badRequestStatusCode || res.Body == nil) {
		c.networkValidator.SetIsBackendAvailable(err.Error())
		return &EvaluationPrerunDataResponse{}, err
	}

	var evaluationPrerunDataResponse = &EvaluationPrerunDataResponse{}
	err = json.Unmarshal(res.Body, &evaluationPrerunDataResponse)
	if err != nil {
		return &EvaluationPrerunDataResponse{}, err
	}

	return evaluationPrerunDataResponse, nil
}

type RuleData struct {
	Identifier string `json:"ruleIdentifier"`
	Name       string `json:"ruleName"`
}

type FileData struct {
	FilePath            string `json:"filepath"`
	ConfigurationsCount int    `json:"configurationsCount"`
}

type Configuration struct {
	Name        string `json:"metadataName"`
	Kind        string `json:"kind"`
	Occurrences int    `json:"occurrences"`
}

type FailedRule struct {
	Name             string          `json:"ruleName"`
	MessageOnFailure string          `json:"messageOnFailure"`
	Configurations   []Configuration `json:"configurations"`
}

type EvaluationResultRequest struct {
	ClientId           string                           `json:"clientId"`
	Token              string                           `json:"token"`
	Metadata           *Metadata                        `json:"metadata"`
	K8sVersion         string                           `json:"k8sVersion"`
	PolicyName         string                           `json:"policyName"`
	FailedYamlFiles    []string                         `json:"failedYamlFiles"`
	FailedK8sFiles     []string                         `json:"failedK8sFiles"`
	AllExecutedRules   []RuleData                       `json:"allExecutedRules"`
	AllEvaluatedFiles  []FileData                       `json:"allEvaluatedFiles"`
	PolicyCheckResults map[string]map[string]FailedRule `json:"policyCheckResults"`
}

func (c *CliClient) SendEvaluationResult(request *EvaluationResultRequest) (*SendEvaluationResultsResponse, error) {
	httpRes, err := c.httpClient.Request(http.MethodPost, "/cli/evaluation/result", request, nil)
	if err != nil {
		return nil, err
	}

	var res = &SendEvaluationResultsResponse{}
	err = json.Unmarshal(httpRes.Body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
