package policy

import (
	"fmt"

	"github.com/datreeio/datree/pkg/cliClient"
	internal_policy "github.com/datreeio/datree/pkg/policy"
)

type Policy struct {
	Name  string
	Rules []RuleSchema
}

type RuleSchema struct {
	RuleIdentifier   string
	RuleName         string
	Schema           map[string]interface{}
	MessageOnFailure string
}

func CreatePolicy(policies *cliClient.PrerunPoliciesForEvaluation, policyName string) (Policy, error) {
	defaultRules, err := internal_policy.GetDefaultRules()

	if err != nil {
		return Policy{}, err
	}

	var rules []RuleSchema

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
			return Policy{policyName, []RuleSchema{}}, nil
		}

		rules, err = populateRules(chosenPolicy.Rules, policies.CustomRules, defaultRules.Rules)

		if err != nil {
			return Policy{}, err
		}
	} else {
		for _, defaultRule := range defaultRules.Rules {
			rules = append(rules, RuleSchema{defaultRule.UniqueName, defaultRule.Name, defaultRule.Schema, defaultRule.MessageOnFailure})
		}
		policyName = "Default"
	}

	return Policy{policyName, rules}, nil
}

func populateRules(policyRules []cliClient.Rule, customRules []*cliClient.CustomRule, defaultRules []internal_policy.DefaultRuleDefinition) ([]RuleSchema, error) {
	var rules []RuleSchema

	for _, rule := range policyRules {
		var isCustomRule bool

		for _, customRule := range customRules {
			if rule.Identifier == customRule.Identifier {
				rules = append(rules, RuleSchema{rule.Identifier, customRule.Name, customRule.Schema, customRule.DefaultMessageOnFailure})
				isCustomRule = true
			}
		}

		if !isCustomRule {
			for _, defaultRule := range defaultRules {
				if rule.Identifier == defaultRule.UniqueName {
					rules = append(rules, RuleSchema{rule.Identifier, defaultRule.Name, defaultRule.Schema, defaultRule.MessageOnFailure})
				} else {
					rulesIsNotCustomNotDefaultErr := fmt.Errorf("rule %s is not custom nor default", rule.Identifier)
					return nil, rulesIsNotCustomNotDefaultErr
				}
			}
		}
	}
	return rules, nil
}
