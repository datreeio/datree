package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"

	internal_policy "github.com/datreeio/datree/pkg/policy"

	"github.com/datreeio/datree/pkg/cliClient"

	"github.com/stretchr/testify/assert"
)

const policiesJsonPath = "internal/fixtures/policyAsCode/policies.json"

func TestCreatePolicy(t *testing.T) {
	err := os.Chdir("../../")

	fileReader := fileReader.CreateFileReader(nil)
	policiesJsonStr, err := fileReader.ReadFileContent(policiesJsonPath)

	if err != nil {
		fmt.Errorf("can't read policies json")
	}

	policiesJsonRawData := []byte(policiesJsonStr)

	var policiesJson cliClient.PrerunDataForEvaluationResponse
	err = json.Unmarshal(policiesJsonRawData, &policiesJson)

	if err != nil {
		fmt.Errorf("can't marshel policies json")
	}

	t.Run("Test Create Policy With Default Policy", func(t *testing.T) {
		policy, err := CreatePolicy(policiesJson.PoliciesJson, "")
		var expectedRules []RuleSchema

		defaultRules, err := internal_policy.GetDefaultRules()

		if err != nil {
			fmt.Errorf("can't read default rules")
		}

		for _, defaultRule := range defaultRules.Rules {
			switch defaultRule.UniqueName {
			case "WORKLOAD_INCORRECT_NAMESPACE_VALUE_DEFAULT", "CONTAINERS_INCORRECT_PRIVILEGED_VALUE_TRUE":
				expectedRules = append(expectedRules, RuleSchema{RuleIdentifier: defaultRule.UniqueName, RuleName: defaultRule.Name, Schema: defaultRule.Schema, MessageOnFailure: defaultRule.MessageOnFailure})
				break
			}
		}

		customRuleJsonMap := make(map[string]interface{})
		customRuleSchemaStr := "{\"properties\":{\"metadata\":{\"properties\":{\"labels\":{\"additionalProperties\":false,\"patternProperties\":{\"^.*$\":{\"format\":\"hostname\"}}}}}}}"
		json.Unmarshal([]byte(customRuleSchemaStr), &customRuleJsonMap)

		expectedRules = append(expectedRules, RuleSchema{RuleIdentifier: "CUSTOM_WORKLOAD_INVALID_LABELS_VALUE", RuleName: "Ensure workload has valid label values [CUSTOM RULE]", Schema: customRuleJsonMap, MessageOnFailure: "All lables values must follow the RFC 1123 hostname standard (https://knowledge.broadcom.com/external/article/49542/restrictions-on-valid-host-names.html)"})

		expectedPolicy := Policy{Name: "labels_best_practices", Rules: expectedRules}

		assert.Equal(t, expectedPolicy, policy)
	})

	t.Run("Test Create Policy With Specific Policy", func(t *testing.T) {
		policy, err := CreatePolicy(policiesJson.PoliciesJson, "labels_best_practices2")
		var expectedRules []RuleSchema

		if err != nil {
			fmt.Errorf("can't read default rules")
		}

		customRuleJsonMap := make(map[string]interface{})
		customRuleSchemaStr := "{\"properties\":{\"metadata\":{\"properties\":{\"labels\":{\"additionalProperties\":false,\"patternProperties\":{\"^.*$\":{\"format\":\"hostname\"}}}}}}}"
		json.Unmarshal([]byte(customRuleSchemaStr), &customRuleJsonMap)

		expectedRules = append(expectedRules, RuleSchema{RuleIdentifier: "CUSTOM_WORKLOAD_INVALID_LABELS_VALUE", RuleName: "Ensure workload has valid label values [CUSTOM RULE]", Schema: customRuleJsonMap, MessageOnFailure: "All lables values must follow the RFC 1123 hostname standard (https://knowledge.broadcom.com/external/article/49542/restrictions-on-valid-host-names.html)"})

		expectedPolicy := Policy{Name: "labels_best_practices2", Rules: expectedRules}

		assert.Equal(t, expectedPolicy, policy)
	})
}
