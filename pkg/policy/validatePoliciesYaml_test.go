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

//go:embed validatePoliciesYamlFixtures/missingCustomRules.yaml
var missingCustomRules string

func Test_missingCustomRules(t *testing.T) {
	err := validatePoliciesYaml([]byte(missingCustomRules), "./validatePoliciesYamlFixtures/missingCustomRules.yaml")
	assert.Nil(t, err)
}

//go:embed validatePoliciesYamlFixtures/missingApiVersion.yaml
var missingApiVersion string

func Test_missingApiVersion(t *testing.T) {
	err := validatePoliciesYaml([]byte(missingApiVersion), "./validatePoliciesYamlFixtures/missingApiVersion.yaml")
	assert.EqualError(t, err, "found errors in policies file ./validatePoliciesYamlFixtures/missingApiVersion.yaml:\n(root): missing properties: 'apiVersion'")
}

//go:embed validatePoliciesYamlFixtures/missingPolicyName.yaml
var missingPolicyName string

func Test_missingPolicyName(t *testing.T) {
	err := validatePoliciesYaml([]byte(missingPolicyName), "./validatePoliciesYamlFixtures/missingPolicyName.yaml")
	assert.EqualError(t, err, "found errors in policies file ./validatePoliciesYamlFixtures/missingPolicyName.yaml:\n(root)/policies/0: missing properties: 'name'")
}
