package policy

import (
	"fmt"

	"github.com/pkg/errors"

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

func CreatePolicy(policies []*cliClient.Policy, customRules []*cliClient.CustomRule, policyName string) (Policy, error) {
	defaultRules, err := internal_policy.GetDefaultRules()
	if err != nil {
		return Policy{}, err
	}

	var rules []RuleSchema

	if policies == nil {
		err := errors.New("There are no policies to run on")
		return Policy{}, err
	}

	if policies != nil {
		var chosenPolicy *cliClient.Policy
		getDefaultPolicy := false

		if policyName == "" || policyName == "default" {
			getDefaultPolicy = true
		}

		for _, policy := range policies {
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

		for _, rule := range chosenPolicy.Rules {
			var isCustomRule bool

			for _, customRule := range customRules {
				if rule.Identifier == customRule.Identifier {
					rules = append(rules, RuleSchema{rule.Identifier, customRule.Name, customRule.Schema, customRule.DefaultMessageOnFailure})
					isCustomRule = true
				}
			}

			if !isCustomRule {
				for _, defaultRule := range defaultRules.Rules {
					if rule.Identifier == defaultRule.UniqueName {
						rules = append(rules, RuleSchema{rule.Identifier, defaultRule.Name, defaultRule.Schema, defaultRule.MessageOnFailure})
					}
				}
			}

		}
	}

	return Policy{policyName, rules}, nil
}
