package defaultPolicies_test

import (
	_ "embed"
	"github.com/datreeio/datree/pkg/defaultPolicies"
	"testing"

	"github.com/datreeio/datree/pkg/validatePoliciesYaml"
)

func TestMain(m *testing.M) {
	err := validatePoliciesYaml.ValidatePoliciesYaml([]byte(defaultPolicies.DefaultPoliciesFileContent), "./defaultPolicies.yaml")
	if err != nil {
		panic(err)
	}
}
