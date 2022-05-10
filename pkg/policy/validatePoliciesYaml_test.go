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

//go:embed validatePoliciesYamlFixtures/policyRulesNull.yaml
var policyRulesNull string

func Test_policyRulesNull(t *testing.T) {
	err := validatePoliciesYaml([]byte(policyRulesNull), "./validatePoliciesYamlFixtures/policyRulesNull.yaml")
	assert.Nil(t, err)
}

//go:embed validatePoliciesYamlFixtures/missingCustomRules.yaml
var missingCustomRules string

func Test_missingCustomRules(t *testing.T) {
	err := validatePoliciesYaml([]byte(missingCustomRules), "./validatePoliciesYamlFixtures/missingCustomRules.yaml")
	assert.Nil(t, err)
}

//go:embed validatePoliciesYamlFixtures/missingPolicyRules.yaml
var missingPolicyRules string

func Test_missingPolicyRules(t *testing.T) {
	err := validatePoliciesYaml([]byte(missingPolicyRules), "./validatePoliciesYamlFixtures/missingPolicyRules.yaml")
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

//go:embed validatePoliciesYamlFixtures/wrongApiVersion.yaml
var wrongApiVersion string

func Test_wrongApiVersion(t *testing.T) {
	err := validatePoliciesYaml([]byte(wrongApiVersion), "./validatePoliciesYamlFixtures/wrongApiVersion.yaml")
	assert.EqualError(t, err, "found errors in policies file ./validatePoliciesYamlFixtures/wrongApiVersion.yaml:\n(root)/apiVersion: value must be \"v1\"")
}

// customRule

//go:embed validatePoliciesYamlFixtures/customRuleMissingIdentifier.yaml
var customRuleMissingIdentifier string

func Test_customRuleMissingIdentifier(t *testing.T) {
	err := validatePoliciesYaml([]byte(customRuleMissingIdentifier), "./validatePoliciesYamlFixtures/customRuleMissingIdentifier.yaml")
	assert.EqualError(t, err, "found errors in policies file ./validatePoliciesYamlFixtures/customRuleMissingIdentifier.yaml:\n(root)/customRules/0: missing properties: 'identifier'")
}

//go:embed validatePoliciesYamlFixtures/customRuleMissingName.yaml
var customRuleMissingName string

func Test_customRuleMissingName(t *testing.T) {
	err := validatePoliciesYaml([]byte(customRuleMissingName), "./validatePoliciesYamlFixtures/customRuleMissingName.yaml")
	assert.EqualError(t, err, "found errors in policies file ./validatePoliciesYamlFixtures/customRuleMissingName.yaml:\n(root)/customRules/0: missing properties: 'name'")
}
