package defaultPolicies

import (
	_ "embed"
	"github.com/datreeio/datree/pkg/cliClient"
	yamlConvertor "github.com/ghodss/yaml"
)

//go:embed defaultPolicies.yaml
var defaultPoliciesFileContent string

func GetDefaultPoliciesYamlContent() string {
	return defaultPoliciesFileContent
}

func GetDefaultPoliciesStruct() *cliClient.EvaluationPrerunPolicies {
	var defaultPolicies cliClient.EvaluationPrerunPolicies
	err := yamlConvertor.Unmarshal([]byte(GetDefaultPoliciesYamlContent()), &defaultPolicies)
	if err != nil {
		panic(err)
	}
	return &defaultPolicies
}
