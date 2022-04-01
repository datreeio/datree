package policy

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"

	internal_policy "github.com/datreeio/datree/pkg/policy"

	"github.com/datreeio/datree/pkg/cliClient"

	"github.com/stretchr/testify/assert"
)

const policiesJsonPath = "../../internal/fixtures/policyAsCode/prerun.json"

func TestCreatePolicy(t *testing.T) {
	policiesJson := mockGetPreRunData()

	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	t.Run("Test Create Policy With Default Policy", func(t *testing.T) {
		policy, _ := CreatePolicy(policiesJson.PoliciesJson, "")
		var expectedRules []RuleWithSchema

		defaultRules, err := internal_policy.GetDefaultRules()

		if err != nil {
			panic(err)
		}

		for _, defaultRule := range defaultRules.Rules {
			switch defaultRule.UniqueName {
			case "WORKLOAD_INCORRECT_NAMESPACE_VALUE_DEFAULT":
				expectedRules = append(expectedRules, RuleWithSchema{RuleIdentifier: defaultRule.UniqueName, RuleName: defaultRule.Name, DocumentationUrl: defaultRule.DocumentationUrl, Schema: defaultRule.Schema, MessageOnFailure: "Incorrect value for key `namespace` - use an explicit namespace instead of the default one (`default`)"})
			case "CONTAINERS_INCORRECT_PRIVILEGED_VALUE_TRUE":
				expectedRules = append(expectedRules, RuleWithSchema{RuleIdentifier: defaultRule.UniqueName, RuleName: defaultRule.Name, DocumentationUrl: defaultRule.DocumentationUrl, Schema: defaultRule.Schema, MessageOnFailure: "Incorrect value for key `privileged` - this mode will allow the container thenhjgjgj same access as processes running on the host"})
			}
		}

		customRuleJsonMap := make(map[string]interface{})
		customRuleSchemaStr := "{\"properties\":{\"metadata\":{\"properties\":{\"labels\":{\"additionalProperties\":false,\"patternProperties\":{\"^.*$\":{\"format\":\"hostname\"}}}}}}}"
		err = json.Unmarshal([]byte(customRuleSchemaStr), &customRuleJsonMap)
		if err != nil {
			panic(err)
		}

		expectedRules = append(expectedRules, RuleWithSchema{RuleIdentifier: "CUSTOM_WORKLOAD_INVALID_LABELS_VALUE", RuleName: "Ensure workload has valid label values [CUSTOM RULE]", Schema: customRuleJsonMap, MessageOnFailure: "All lables values must follow the RFC 1123 hostname standard (https://knowledge.broadcom.com/external/article/49542/restrictions-on-valid-host-names.html)"})

		expectedPolicy := Policy{Name: "labels_best_practices", Rules: expectedRules}

		assert.Equal(t, expectedPolicy, policy)
	})

	t.Run("Test Create Policy With Specific Policy", func(t *testing.T) {
		policy, err := CreatePolicy(policiesJson.PoliciesJson, "labels_best_practices2")
		var expectedRules []RuleWithSchema

		if err != nil {
			panic(err)
		}

		customRuleJsonMap := make(map[string]interface{})
		customRuleSchemaStr := "{\"properties\":{\"metadata\":{\"properties\":{\"labels\":{\"additionalProperties\":false,\"patternProperties\":{\"^.*$\":{\"format\":\"hostname\"}}}}}}}"
		err = json.Unmarshal([]byte(customRuleSchemaStr), &customRuleJsonMap)
		if err != nil {
			panic(err)
		}

		expectedRules = append(expectedRules, RuleWithSchema{RuleIdentifier: "CUSTOM_WORKLOAD_INVALID_LABELS_VALUE", RuleName: "Ensure workload has valid label values [CUSTOM RULE]", Schema: customRuleJsonMap, MessageOnFailure: "All lables values must follow the RFC 1123 hostname standard (https://knowledge.broadcom.com/external/article/49542/restrictions-on-valid-host-names.html)"})

		expectedPolicy := Policy{Name: "labels_best_practices2", Rules: expectedRules}

		assert.Equal(t, expectedPolicy, policy)
	})

	t.Run("Test Create Policy With Custom Rules", func(t *testing.T) {
		policy, err := CreatePolicy(policiesJson.PoliciesJson, "labels_best_practices3")
		var expectedRules []RuleWithSchema
		if err != nil {
			panic(err)
		}

		jsonSchemaStr := "{\"type\":\"object\",\"properties\":{\"apiVersion\":{\"type\":\"string\"}},\"required\":[\"apiVersion\"]}"
		customRuleJsonSchema := make(map[string]interface{})
		err = json.Unmarshal([]byte(jsonSchemaStr), &customRuleJsonSchema)
		if err != nil {
			panic(err)
		}
		expectedRules = append(expectedRules, RuleWithSchema{RuleIdentifier: "UNIQUE2", RuleName: "rule unique 2", Schema: customRuleJsonSchema, MessageOnFailure: "default message for rule fail number 2"})
		expectedRules = append(expectedRules, RuleWithSchema{RuleIdentifier: "UNIQUE3", RuleName: "rule unique 3", Schema: customRuleJsonSchema, MessageOnFailure: "default message for rule fail number 3"})

		assert.Equal(t, expectedRules, policy.Rules)
	})
}

func mockGetPreRunData() *cliClient.EvaluationPrerunDataResponse {
	fileReader := fileReader.CreateFileReader(nil)
	policiesJsonStr, err := fileReader.ReadFileContent(policiesJsonPath)

	if err != nil {
		panic(err)
	}

	policiesJsonRawData := []byte(policiesJsonStr)

	var policiesJson *cliClient.EvaluationPrerunDataResponse
	err = json.Unmarshal(policiesJsonRawData, &policiesJson)

	if err != nil {
		panic(err)
	}
	return policiesJson
}
