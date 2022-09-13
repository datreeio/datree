package defaultPolicies

import (
	_ "embed"
	"github.com/datreeio/datree/pkg/validatePoliciesYaml"
	"testing"
)

//go:embed defaultPolicies.yaml
var defaultPoliciesFileContent string

var defaultRulesSchemaFileContent string

func TestMain(m *testing.M) {
	err := validatePoliciesYaml.ValidatePoliciesYaml([]byte(defaultPoliciesFileContent), "./defaultPolicies.yaml")
	if err != nil {
		panic(err)
	}
}
