package policy

import (
	_ "embed"
	"fmt"

	"github.com/ghodss/yaml"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
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

//go:embed policiesSchema.json
var policiesSchemaContent string

func CreatePolicy(policies *cliClient.EvaluationPrerunPolicies, policyName string) (Policy, error) {
	defaultRules, err := internal_policy.GetDefaultRules()

	if err != nil {
		return Policy{}, err
	}

	var rules []RuleWithSchema

	if policies != nil {
		var chosenPolicy *cliClient.Policy

		for _, policy := range policies.Policies {
			if policyName == "" && policy.IsDefault {
				chosenPolicy = policy
				policyName = chosenPolicy.Name
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

		rules, err = populateRules(chosenPolicy.Rules, policies.CustomRules, defaultRules.Rules)

		if err != nil {
			return Policy{}, err
		}
	} else {
		policy := createDefaultPolicy(defaultRules)
		return policy, nil
	}

	return Policy{policyName, rules}, nil
}

func GetPoliciesFileFromPath(path string) (*cliClient.EvaluationPrerunPolicies, error) {
	fileReader := fileReader.CreateFileReader(nil)
	policiesStr, err := fileReader.ReadFileContent(path)
	if err != nil {
		return nil, err
	}

	err = validatePoliciesYaml(policiesStr, path)
	if err != nil {
		return nil, err
	}

	var policies *cliClient.EvaluationPrerunPolicies
	policiesBytes, err := yaml.YAMLToJSON([]byte(policiesStr))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(policiesBytes, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func populateRules(policyRules []cliClient.Rule, customRules []*cliClient.CustomRule, defaultRules []*internal_policy.DefaultRuleDefinition) ([]RuleWithSchema, error) {
	var rules = []RuleWithSchema{}

	if policyRules == nil {
		return rules, nil
	}

	for _, rule := range policyRules {
		customRule := getCustomRuleByIdentifier(customRules, rule.Identifier)

		if customRule != nil {
			rules = append(rules, RuleWithSchema{rule.Identifier, customRule.Name, customRule.Schema, rule.MessageOnFailure})
		} else {
			defaultRule := getDefaultRuleByIdentifier(defaultRules, rule.Identifier)

			if defaultRule != nil {
				rules = append(rules, RuleWithSchema{rule.Identifier, defaultRule.Name, defaultRule.Schema, rule.MessageOnFailure})
			} else {
				rulesIsNotCustomNorDefaultErr := fmt.Errorf("rule %s is not custom nor default", rule.Identifier)
				return nil, rulesIsNotCustomNorDefaultErr
			}
		}
	}

	return rules, nil
}

func getDefaultRuleByIdentifier(defaultRules []*internal_policy.DefaultRuleDefinition, identifier string) *internal_policy.DefaultRuleDefinition {
	for _, defaultRule := range defaultRules {
		if identifier == defaultRule.UniqueName {
			return defaultRule
		}
	}

	return nil
}

func getCustomRuleByIdentifier(customRules []*cliClient.CustomRule, identifier string) *cliClient.CustomRule {
	for _, customRule := range customRules {
		if identifier == customRule.Identifier {
			return customRule
		}
	}

	return nil
}

func createDefaultPolicy(defaultRules *internal_policy.DefaultRulesDefinitions) Policy {
	var rules []RuleWithSchema

	for _, defaultRule := range defaultRules.Rules {
		if defaultRule.EnabledByDefault {
			rules = append(rules, RuleWithSchema{defaultRule.UniqueName, defaultRule.Name, defaultRule.Schema, defaultRule.MessageOnFailure})
		}
	}

	return Policy{"Default", rules}
}

func validatePoliciesYaml(content string, policyYamlPath string) error {
	jsonSchemaValidator := jsonSchemaValidator.New()
	result, err := jsonSchemaValidator.Validate(policiesSchemaContent, content)

	if err != nil {
		return err
	}

	if !result.Valid() {
		validationErrors := fmt.Errorf("Found errors in policies file %s:\n", policyYamlPath)

		for _, validationError := range result.Errors() {
			validationErrors = fmt.Errorf("%s\n%s", validationErrors, validationError)
		}

		return validationErrors
	}

	return nil
}
