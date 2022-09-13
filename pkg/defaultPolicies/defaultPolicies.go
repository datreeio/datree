package defaultPolicies

import (
	_ "embed"

	yamlConvertor "github.com/ghodss/yaml"
)

type EvaluationPrerunPolicies struct {
	ApiVersion  string        `json:"apiVersion"`
	CustomRules []*CustomRule `json:"customRules"`
	Policies    []*Policy     `json:"policies"`
}

type CustomRule struct {
	Identifier              string      `json:"identifier"`
	Name                    string      `json:"name"`
	DefaultMessageOnFailure string      `json:"defaultMessageOnFailure"`
	Schema                  interface{} `json:"schema"`
	JsonSchema              string      `json:"jsonSchema"`
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

//go:embed defaultPolicies.yaml
var DefaultPoliciesFileContent string

func GetDefaultPoliciesStruct() *EvaluationPrerunPolicies {
	var defaultPolicies EvaluationPrerunPolicies
	err := yamlConvertor.Unmarshal([]byte(DefaultPoliciesFileContent), &defaultPolicies)
	if err != nil {
		panic(err)
	}
	return &defaultPolicies
}
