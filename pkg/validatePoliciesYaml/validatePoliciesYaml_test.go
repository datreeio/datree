package validatePoliciesYaml

import (
	_ "embed"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_fixtures/customRulesNull.yaml
var customRulesNull string

//go:embed test_fixtures/policyRulesNull.yaml
var policyRulesNull string

//go:embed test_fixtures/missingCustomRules.yaml
var missingCustomRules string

//go:embed test_fixtures/missingPolicyRules.yaml
var missingPolicyRules string

//go:embed test_fixtures/missingApiVersion.yaml
var missingApiVersion string

//go:embed test_fixtures/missingPolicyName.yaml
var missingPolicyName string

//go:embed test_fixtures/missingSchema.yaml
var missingSchema string

//go:embed test_fixtures/multipleDefaultPolicy.yaml
var multipleDefaultPolicy string

//go:embed test_fixtures/wrongApiVersion.yaml
var wrongApiVersion string

//go:embed test_fixtures/customRuleMissingIdentifier.yaml
var customRuleMissingIdentifier string

//go:embed test_fixtures/customRuleMissingName.yaml
var customRuleMissingName string

//go:embed test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml
var customRuleMissingDefaultMessageOnFailure string

//go:embed test_fixtures/customRuleValidSchema.yaml
var customRuleValidSchema string

//go:embed test_fixtures/customRuleInvalidSchema.yaml
var customRuleInvalidSchema string

//go:embed test_fixtures/customRuleInvalidJsonSchema.yaml
var customRuleInvalidJsonSchema string

//go:embed test_fixtures/customRuleIdentifierNotUnique1.yaml
var customRuleIdentifierNotUnique1 string

//go:embed test_fixtures/customRuleIdentifierNotUnique2.yaml
var customRuleIdentifierNotUnique2 string

func assertValidationResult(t *testing.T, policiesFile string, policiesFilePath string, expectedError error) {
	err := ValidatePoliciesYaml([]byte(policiesFile), policiesFilePath)
	assert.Equal(t, err, expectedError)
}

func TestValidatePoliciesYaml(t *testing.T) {
	assertValidationResult(t, customRulesNull, "./test_fixtures/customRulesNull.yaml", nil)
	assertValidationResult(t, policyRulesNull, "./test_fixtures/policyRulesNull.yaml", nil)
	assertValidationResult(t, missingCustomRules, "./test_fixtures/missingCustomRules.yaml", nil)
	assertValidationResult(t, missingPolicyRules, "./test_fixtures/missingPolicyRules.yaml", nil)
	assertValidationResult(t, missingApiVersion, "./test_fixtures/missingApiVersion.yaml", errors.New("found errors in policies file ./test_fixtures/missingApiVersion.yaml:\n(root): missing properties: 'apiVersion'"))
	assertValidationResult(t, missingPolicyName, "./test_fixtures/missingPolicyName.yaml", errors.New("found errors in policies file ./test_fixtures/missingPolicyName.yaml:\n(root)/policies/0: missing properties: 'name'"))
	assertValidationResult(t, missingSchema, "./test_fixtures/missingSchema.yaml", errors.New("invalid policy file: rule identifier \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" has no associated schema"))
	assertValidationResult(t, multipleDefaultPolicy, "./test_fixtures/multipleDefaultPolicy.yaml", errors.New("invalid policy file: multiple policies are set to default"))
	assertValidationResult(t, wrongApiVersion, "./test_fixtures/wrongApiVersion.yaml", errors.New("found errors in policies file ./test_fixtures/wrongApiVersion.yaml:\n(root)/apiVersion: value must be \"v1\""))

	// customRule
	assertValidationResult(t, customRuleMissingIdentifier, "./test_fixtures/customRuleMissingIdentifier.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleMissingIdentifier.yaml:\n(root)/customRules/0: missing properties: 'identifier'"))
	assertValidationResult(t, customRuleMissingName, "./test_fixtures/customRuleMissingName.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleMissingName.yaml:\n(root)/customRules/0: missing properties: 'name'"))
	assertValidationResult(t, customRuleMissingDefaultMessageOnFailure, "./test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml:\n(root)/customRules/0: missing properties: 'defaultMessageOnFailure'"))
	assertValidationResult(t, customRuleValidSchema, "./test_fixtures/customRuleValidSchema.yaml", nil)
	assertValidationResult(t, customRuleInvalidSchema, "./test_fixtures/customRuleInvalidSchema.yaml", errors.New("invalid policy file: rule \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" has schema error: has a primitive type that is NOT VALID -- given: /arrayi/ Expected valid values are:[array boolean integer number null object string]"))
	assertValidationResult(t, customRuleInvalidJsonSchema, "./test_fixtures/customRuleInvalidJsonSchema.yaml", errors.New("invalid policy file: rule \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" has schema error: Invalid type. Expected: array of schemas, given: definitions"))
	assertValidationResult(t, customRuleIdentifierNotUnique1, "./test_fixtures/customRuleIdentifierNotUnique1.yaml", errors.New("invalid policy file: identifier \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" is used in atleast two rules"))
	assertValidationResult(t, customRuleIdentifierNotUnique2, "./test_fixtures/customRuleIdentifierNotUnique2.yaml", errors.New("invalid policy file: identifier \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" used atleast two times in policy \"Default\""))
}
