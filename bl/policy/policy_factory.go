package policy

import (
	"encoding/json"
	"fmt"
	"github.com/datreeio/datree/pkg/defaultPolicies"

	"github.com/datreeio/datree/pkg/defaultRules"
)

type Policy struct {
	Name  string
	Rules []RuleWithSchema
}

type RuleWithSchema struct {
	RuleIdentifier   string
	RuleName         string
	DocumentationUrl string
	Schema           interface{}
	MessageOnFailure string
}

func CreatePolicy(policies *defaultPolicies.EvaluationPrerunPolicies, policyName string, registrationURL string, defaultRules *defaultRules.DefaultRulesDefinitions, isAnonymous bool) (Policy, error) {
	if policies == nil {
		// policies should never be nil because of the fallback of defaultPolicies.yaml
		panic("policies is nil")
	}

	var err error

	var rules []RuleWithSchema

	var chosenPolicy *defaultPolicies.Policy

	if policyName == "" {
		for _, policy := range policies.Policies {
			if policy.IsDefault {
				chosenPolicy = policy
				policyName = chosenPolicy.Name
				break
			}
		}
	} else {
		for _, policy := range policies.Policies {
			if policy.Name == policyName {
				chosenPolicy = policy
				policyName = chosenPolicy.Name
				break
			}
		}
	}

	if chosenPolicy == nil {
		if isAnonymous {
			return Policy{}, fmt.Errorf("policy %s doesn't exist, sign in to the dashboard to customize your policies: %s", policyName, registrationURL)
		} else {
			return Policy{}, fmt.Errorf("policy %s doesn't exist", policyName)
		}
	}

	rules, err = populateRules(chosenPolicy.Rules, policies.CustomRules, defaultRules.Rules)

	if err != nil {
		return Policy{}, err
	}

	return Policy{policyName, rules}, nil
}

func populateRules(policyRules []defaultPolicies.Rule, customRules []*defaultPolicies.CustomRule, defaultRules []*defaultRules.DefaultRuleDefinition) ([]RuleWithSchema, error) {
	var rules = []RuleWithSchema{}

	if policyRules == nil {
		return rules, nil
	}

	for _, rule := range policyRules {
		customRule := getCustomRuleByIdentifier(customRules, rule.Identifier)

		if customRule != nil {
			if customRule.Schema == nil {
				schema := make(map[string]interface{})
				err := json.Unmarshal([]byte(customRule.JsonSchema), &schema)
				if err != nil {
					return nil, err
				}
				rules = append(rules, RuleWithSchema{rule.Identifier, customRule.Name, "", schema, rule.MessageOnFailure})
			} else {
				rules = append(rules, RuleWithSchema{rule.Identifier, customRule.Name, "", customRule.Schema, rule.MessageOnFailure})
			}
		} else {
			defaultRule := getDefaultRuleByIdentifier(defaultRules, rule.Identifier)

			if defaultRule != nil {
				rules = append(rules, RuleWithSchema{rule.Identifier, defaultRule.Name, defaultRule.DocumentationUrl, defaultRule.Schema, rule.MessageOnFailure})
			} else {
				rulesIsNotCustomNorDefaultErr := fmt.Errorf("rule %s is not custom nor default", rule.Identifier)
				return nil, rulesIsNotCustomNorDefaultErr
			}
		}
	}

	return rules, nil
}

func getDefaultRuleByIdentifier(defaultRules []*defaultRules.DefaultRuleDefinition, identifier string) *defaultRules.DefaultRuleDefinition {
	for _, defaultRule := range defaultRules {
		if identifier == defaultRule.UniqueName {
			return defaultRule
		}
	}

	return nil
}

func getCustomRuleByIdentifier(customRules []*defaultPolicies.CustomRule, identifier string) *defaultPolicies.CustomRule {
	for _, customRule := range customRules {
		if identifier == customRule.Identifier {
			return customRule
		}
	}

	return nil
}
