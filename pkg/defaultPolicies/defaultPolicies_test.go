package defaultPolicies_test

import (
	"github.com/datreeio/datree/pkg/defaultPolicies"
	"github.com/datreeio/datree/pkg/validatePoliciesYaml"
	"testing"
)

func TestMain(m *testing.M) {
	defaultPoliciesFileContent := defaultPolicies.GetDefaultPoliciesYamlContent()
	err := validatePoliciesYaml.ValidatePoliciesYaml([]byte(defaultPoliciesFileContent), "./defaultPolicies.yaml")
	if err != nil {
		panic(err)
	}
}
