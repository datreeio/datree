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

//go:embed test_fixtures/noDefaultPolicy.yaml
var noDefaultPolicy string

//go:embed test_fixtures/wrongApiVersion.yaml
var wrongApiVersion string

//go:embed test_fixtures/bothSchemaAndJsonSchemaDefined.yaml
var bothSchemaAndJsonSchemaDefined string

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

//go:embed test_fixtures/customRuleJsonSchemaInvalidJson.yaml
var customRuleJsonSchemaInvalidJson string

//go:embed test_fixtures/identifierNotDefined.yaml
var identifierNotDefined string

//go:embed test_fixtures/customRuleIdentifierNotUnique.yaml
var customRuleIdentifierNotUnique string

//go:embed test_fixtures/customRuleIdentifierMatchDefaultRule.yaml
var customRuleIdentifierMatchDefaultRule string

//go:embed test_fixtures/duplicateRuleIdentifier.yaml
var duplicateRuleIdentifier string

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
	assertValidationResult(t, missingSchema, "./test_fixtures/missingSchema.yaml", errors.New("found errors in policies file ./test_fixtures/missingSchema.yaml:\n(root)/customRules/1: Exactly one of [schema,jsonSchema] should be defined per custom rule"))
	assertValidationResult(t, multipleDefaultPolicy, "./test_fixtures/multipleDefaultPolicy.yaml", errors.New("found errors in policies file ./test_fixtures/multipleDefaultPolicy.yaml:\n(root)/policies: Should have exactly one policy set as default"))
	assertValidationResult(t, noDefaultPolicy, "./test_fixtures/noDefaultPolicy.yaml", errors.New("found errors in policies file ./test_fixtures/noDefaultPolicy.yaml:\n(root)/policies: Should have exactly one policy set as default"))
	assertValidationResult(t, wrongApiVersion, "./test_fixtures/wrongApiVersion.yaml", errors.New("found errors in policies file ./test_fixtures/wrongApiVersion.yaml:\n(root)/apiVersion: value must be \"v1\""))

	// customRule
	assertValidationResult(t, customRuleMissingIdentifier, "./test_fixtures/customRuleMissingIdentifier.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleMissingIdentifier.yaml:\n(root)/customRules/0: missing properties: 'identifier'"))
	assertValidationResult(t, customRuleMissingName, "./test_fixtures/customRuleMissingName.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleMissingName.yaml:\n(root)/customRules/0: missing properties: 'name'"))
	assertValidationResult(t, customRuleMissingDefaultMessageOnFailure, "./test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleMissingDefaultMessageOnFailure.yaml:\n(root)/customRules/0: missing properties: 'defaultMessageOnFailure'"))
	assertValidationResult(t, customRuleValidSchema, "./test_fixtures/customRuleValidSchema.yaml", nil)
	assertValidationResult(t, customRuleInvalidSchema, "./test_fixtures/customRuleInvalidSchema.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleInvalidSchema.yaml:\n(root)/customRules/1/schema: has a primitive type that is NOT VALID -- given: /arrayi/ Expected valid values are:[array boolean integer number null object string]"))
	assertValidationResult(t, customRuleInvalidJsonSchema, "./test_fixtures/customRuleInvalidJsonSchema.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleInvalidJsonSchema.yaml:\n(root)/customRules/1/jsonSchema: Invalid type. Expected: array of schemas, given: definitions"))
	assertValidationResult(t, bothSchemaAndJsonSchemaDefined, "./test_fixtures/bothSchemaAndJsonSchemaDefined.yaml", errors.New("found errors in policies file ./test_fixtures/bothSchemaAndJsonSchemaDefined.yaml:\n(root)/customRules/0: Exactly one of [schema,jsonSchema] should be defined per custom rule"))
	assertValidationResult(t, identifierNotDefined, "./test_fixtures/identifierNotDefined.yaml", errors.New("found errors in policies file ./test_fixtures/identifierNotDefined.yaml:\n(root)/policies/0/rules: identifier \"SOME_IDENTIFIER_NAME\" is neither custom nor default"))
	assertValidationResult(t, customRuleIdentifierNotUnique, "./test_fixtures/customRuleIdentifierNotUnique.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleIdentifierNotUnique.yaml:\n(root)/customRules: identifier \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" is used in more than one custom rule"))
	assertValidationResult(t, customRuleIdentifierMatchDefaultRule, "./test_fixtures/customRuleIdentifierMatchDefaultRule.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleIdentifierMatchDefaultRule.yaml:\n(root)/customRules/0: a default rule with same identifier \"RESOURCE_MISSING_NAME\" already exists"))
	assertValidationResult(t, duplicateRuleIdentifier, "./test_fixtures/duplicateRuleIdentifier.yaml", errors.New("found errors in policies file ./test_fixtures/duplicateRuleIdentifier.yaml:\n(root)/policies/0/rules: identifier \"PODDISRUPTIONBUDGET_DENY_ZERO_VOLUNTARY_DISRUPTION\" is used more than once in policy"))
	assertValidationResult(t, customRuleJsonSchemaInvalidJson, "./test_fixtures/customRuleJsonSchemaInvalidJson.yaml", errors.New("found errors in policies file ./test_fixtures/customRuleJsonSchemaInvalidJson.yaml:\n(root)/customRules/1/jsonSchema: invalid character '2' looking for beginning of object key string"))
}
