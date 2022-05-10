package policy

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed validatePoliciesYamlFixtures/customRulesNull.yaml
var customRulesNull string

func Test_customRulesNull(t *testing.T) {
	err := validatePoliciesYaml([]byte(customRulesNull), "./validatePoliciesYamlFixtures/customRulesNull.yaml")
	assert.Nil(t, err)
}

//go:embed validatePoliciesYamlFixtures/missingApiVersion.yaml
var missingApiVersion string

func Test_missingApiVersion(t *testing.T) {
	err := validatePoliciesYaml([]byte(missingApiVersion), "./validatePoliciesYamlFixtures/missingApiVersion.yaml")
	assert.EqualError(t, err, "found errors in policies file ./validatePoliciesYamlFixtures/missingApiVersion.yaml:\nmissing properties: 'apiVersion'")
}
