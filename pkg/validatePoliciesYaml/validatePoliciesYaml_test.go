package validatePoliciesYaml

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_fixtures/customRulesNull.yaml
var customRulesNull string

func Test_customRulesNull(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(customRulesNull), "./test_fixtures/customRulesNull.yaml")
	assert.Nil(t, err)
}

//go:embed test_fixtures/policyRulesNull.yaml
var policyRulesNull string

func Test_policyRulesNull(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(policyRulesNull), "./test_fixtures/policyRulesNull.yaml")
	assert.Nil(t, err)
}

//go:embed test_fixtures/missingCustomRules.yaml
var missingCustomRules string

func Test_missingCustomRules(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(missingCustomRules), "./test_fixtures/missingCustomRules.yaml")
	assert.Nil(t, err)
}

//go:embed test_fixtures/missingPolicyRules.yaml
var missingPolicyRules string

func Test_missingPolicyRules(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(missingPolicyRules), "./test_fixtures/missingPolicyRules.yaml")
	assert.Nil(t, err)
}

//go:embed test_fixtures/missingApiVersion.yaml
var missingApiVersion string

func Test_missingApiVersion(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(missingApiVersion), "./test_fixtures/missingApiVersion.yaml")
	assert.EqualError(t, err, "found errors in policies file ./test_fixtures/missingApiVersion.yaml:\n(root): missing properties: 'apiVersion'")
}

//go:embed test_fixtures/missingPolicyName.yaml
var missingPolicyName string

func Test_missingPolicyName(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(missingPolicyName), "./test_fixtures/missingPolicyName.yaml")
	assert.EqualError(t, err, "found errors in policies file ./test_fixtures/missingPolicyName.yaml:\n(root)/policies/0: missing properties: 'name'")
}

//go:embed test_fixtures/wrongApiVersion.yaml
var wrongApiVersion string

func Test_wrongApiVersion(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(wrongApiVersion), "./test_fixtures/wrongApiVersion.yaml")
	assert.EqualError(t, err, "found errors in policies file ./test_fixtures/wrongApiVersion.yaml:\n(root)/apiVersion: value must be \"v1\"")
}

// customRule

//go:embed test_fixtures/customRuleMissingIdentifier.yaml
var customRuleMissingIdentifier string

func Test_customRuleMissingIdentifier(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(customRuleMissingIdentifier), "./test_fixtures/customRuleMissingIdentifier.yaml")
	assert.EqualError(t, err, "found errors in policies file ./test_fixtures/customRuleMissingIdentifier.yaml:\n(root)/customRules/0: missing properties: 'identifier'")
}

//go:embed test_fixtures/customRuleMissingName.yaml
var customRuleMissingName string

func Test_customRuleMissingName(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(customRuleMissingName), "./test_fixtures/customRuleMissingName.yaml")
	assert.EqualError(t, err, "found errors in policies file ./test_fixtures/customRuleMissingName.yaml:\n(root)/customRules/0: missing properties: 'name'")
}

//go:embed test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml
var customRuleMissingDefaultMessageOnFailure string

func Test_customRuleMissingDefaultMessageOnFailure(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(customRuleMissingDefaultMessageOnFailure), "./test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml")
	assert.EqualError(t, err, "found errors in policies file ./test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml:\n(root)/customRules/0: missing properties: 'defaultMessageOnFailure'")
}

//go:embed test_fixtures/customRuleValidSchema.yaml
var customRuleValidSchema string

func Test_customRuleValidSchema(t *testing.T) {
	err := ValidatePoliciesYaml([]byte(customRuleValidSchema), "./test_fixtures/customRuleValidSchema.yaml")
	assert.Nil(t, err)
}
