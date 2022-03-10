package policy

import (
	"fmt"

	"github.com/datreeio/datree/pkg/cliClient"
	internal_policy "github.com/datreeio/datree/pkg/policy"
)

type Policy struct {
	Name  string
	Rules []RuleWithSchema
}

type RuleWithSchema struct {
	RuleIdentifier   string
	RuleName         string
	Schema           map[string]interface{}
	MessageOnFailure string
}

func CreatePolicy(policies *cliClient.EvaluationPrerunPolicies, policyName string) (Policy, error) {
	defaultRules, err := internal_policy.GetDefaultRules()

	if err != nil {
		return Policy{}, err
	}

	var rules []RuleWithSchema

	if policies != nil {
		var chosenPolicy *cliClient.Policy
		getDefaultPolicy := false

		if policyName == "" {
			getDefaultPolicy = true
		}

		for _, policy := range policies.Policies {
			if getDefaultPolicy && policy.IsDefault {
				chosenPolicy = policy
				break
			} else if policy.Name == policyName {
				chosenPolicy = policy
				break
			}
		}

		if chosenPolicy == nil {
			err := fmt.Errorf("policy %s doesn't exist", policyName)
			return Policy{}, err
		}

		policyName = chosenPolicy.Name

		if chosenPolicy.Rules == nil {
			return Policy{policyName, []RuleWithSchema{}}, nil
		}

		rules, err = populateRules(chosenPolicy.Rules, policies.CustomRules, defaultRules.Rules)

		if err != nil {
			return Policy{}, err
		}
	} else {
		for _, defaultRule := range defaultRules.Rules {
			rules = append(rules, RuleWithSchema{defaultRule.UniqueName, defaultRule.Name, defaultRule.Schema, defaultRule.MessageOnFailure})
		}
		policyName = "Default"
	}

	return Policy{policyName, rules}, nil
}

func populateRules(policyRules []cliClient.Rule, customRules []*cliClient.CustomRule, defaultRules []internal_policy.DefaultRuleDefinition) ([]RuleWithSchema, error) {
	var rules []RuleWithSchema

	for _, rule := range policyRules {
		var isCustomRule bool
		var isDefaultRule bool

		for _, customRule := range customRules {
			if rule.Identifier == customRule.Identifier {
				isCustomRule = true
				rules = append(rules, RuleWithSchema{rule.Identifier, customRule.Name, customRule.Schema, rule.MessageOnFailure})
				break
			}
		}

		if !isCustomRule {
			for _, defaultRule := range defaultRules {
				if rule.Identifier == defaultRule.UniqueName {
					isDefaultRule = true
					rules = append(rules, RuleWithSchema{rule.Identifier, defaultRule.Name, defaultRule.Schema, rule.MessageOnFailure})
					break
				}
			}
			if !isDefaultRule {
				rulesIsNotCustomNorDefaultErr := fmt.Errorf("rule %s is not custom nor default", rule.Identifier)
				return nil, rulesIsNotCustomNorDefaultErr
			}
		}
	}

	return rules, nil
}
