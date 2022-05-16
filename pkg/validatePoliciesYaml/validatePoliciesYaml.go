package validatePoliciesYaml

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/defaultRules"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed policiesSchema.json
var policiesSchemaContent string

func ValidatePoliciesYaml(content []byte, policyYamlPath string) error {
	jsonSchemaValidator := jsonSchemaValidator.New()
	jsonContent, _ := yaml.YAMLToJSON(content)
	errorsResult, err := jsonSchemaValidator.Validate(policiesSchemaContent, jsonContent)

	if err != nil {
		return err
	}

	errorPrefix := fmt.Errorf("found errors in policies file %s:", policyYamlPath)
	if errorsResult != nil {
		validationErrors := errorPrefix

		for _, validationError := range errorsResult {
			validationErrors = fmt.Errorf("%s\n(root)%s: %s", validationErrors, validationError.InstanceLocation, validationError.Error)
		}

		return validationErrors
	}

	err = validatePoliciesContent(jsonContent)
	if err != nil {
		return fmt.Errorf("%s\n%s", errorPrefix, err)
	}
	return nil
}

func validatePoliciesContent(content []byte) error {
	// unmarshal the content
	var schema *cliClient.EvaluationPrerunPolicies
	err := json.Unmarshal(content, &schema)
	if err != nil {
		return err
	}

	// validate that exactly one policy is set to default
	err = validateSingleDefaultPolicy(schema.Policies)
	if err != nil {
		return err
	}

	// validate if the policy file has any identifier related issues
	err = validateIdentifier(schema.Policies, schema.CustomRules)
	if err != nil {
		return err
	}

	// validate the schema of each rule
	err = validateSchemaField(schema.CustomRules)
	return err
}

func validateIdentifier(policies []*cliClient.Policy, customRules []*cliClient.CustomRule) error {

	err := checkIdentifierInPolicy(policies, customRules)
	if err != nil {
		return err
	}

	err = checkCustomRulesIdentifiersUniqueness(customRules)
	return err
}

func checkIdentifierInPolicy(policies []*cliClient.Policy, customRules []*cliClient.CustomRule) error {
	err := checkIdentifierUniquenessInPolicy(policies)
	if err != nil {
		return err
	}

	err = checkIdentifierExistence(policies, customRules)
	if err != nil {
		return err
	}
	return nil
}

func checkIdentifierExistence(policies []*cliClient.Policy, customRules []*cliClient.CustomRule) error {
	defaultRules, err := defaultRules.GetDefaultRules()
	if err != nil {
		return err
	}

	for index, policy := range policies {
		rules := policy.Rules

		for _, rule := range rules {
			found := false
			identifier := rule.Identifier
			for _, customRule := range customRules {
				if identifier == customRule.Identifier {
					found = true
					break
				}
			}
			if found {
				continue
			}
			for _, defaultRule := range defaultRules.Rules {
				if identifier == defaultRule.UniqueName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("(root)/policies/%d/rules: identifier \"%s\" is neither custom nor default", index, identifier)
			}
		}
	}

	return nil
}

func checkIdentifierUniquenessInPolicy(policies []*cliClient.Policy) error {
	for index, policy := range policies {
		propertyValuesExistenceMap := make(map[string]bool)
		rules := policy.Rules
		for _, rule := range rules {
			identifier := rule.Identifier
			if propertyValuesExistenceMap[identifier] {
				return fmt.Errorf("(root)/policies/%d/rules: identifier \"%s\" is used more than once in policy", index, identifier)
			}
			propertyValuesExistenceMap[identifier] = true
		}
	}
	return nil
}

func checkCustomRulesIdentifiersUniqueness(customRules []*cliClient.CustomRule) error {
	defaultRules, err := defaultRules.GetDefaultRules()
	if err != nil {
		return err
	}

	defaultRulesIdentifierToExistenceMap := make(map[string]bool)
	for _, defaultRule := range defaultRules.Rules {
		defaultRulesIdentifierToExistenceMap[defaultRule.UniqueName] = true
	}

	customRulesIdentifierToExistenceMap := make(map[string]bool)
	for index, customRule := range customRules {
		identifier := customRule.Identifier
		if customRulesIdentifierToExistenceMap[identifier] {
			return fmt.Errorf("(root)/customRules: identifier \"%s\" is used in more than one custom rule", identifier)
		}
		customRulesIdentifierToExistenceMap[customRule.Identifier] = true

		if defaultRulesIdentifierToExistenceMap[identifier] {
			return fmt.Errorf("(root)/customRules/%d: a default rule with same identifier \"%s\" already exists", index, identifier)
		}
	}

	return nil
}

func validateSingleDefaultPolicy(policies []*cliClient.Policy) error {
	sawDefault := false
	for _, policy := range policies {

		if policy.IsDefault {
			if sawDefault {
				return fmt.Errorf("(root)/policies: Should have exactly one policy set as default")
			}
			sawDefault = true
		}
	}

	if !sawDefault {
		return fmt.Errorf("(root)/policies: Should have exactly one policy set as default")
	}
	return nil
}

func validateSchemaField(customRules []*cliClient.CustomRule) error {
	for index, rule := range customRules {
		var err error
		var jsonContent string
		if rule.Schema != nil && rule.JsonSchema != "" {
			return fmt.Errorf("(root)/customRules/%d: Exactly one of [schema,jsonSchema] should be defined per custom rule", index)
		}
		var schemaKeyUsed string
		if rule.Schema != nil {
			var content []byte
			schema := rule.Schema
			content, err = json.Marshal(schema)
			if err != nil {
				return fmt.Errorf("(root)/customRules/%d: %s", index, err.Error())
			}
			jsonContent = string(content)
			schemaKeyUsed = "schema"
		} else {
			jsonContent = rule.JsonSchema
			schemaKeyUsed = "jsonSchema"
		}

		if jsonContent == "" {
			return fmt.Errorf("(root)/customRules/%d: Exactly one of [schema,jsonSchema] should be defined per custom rule", index)
		}
		schemaLoader := gojsonschema.NewStringLoader(jsonContent)
		_, err = gojsonschema.NewSchemaLoader().Compile(schemaLoader)
		if err != nil {
			return fmt.Errorf("(root)/customRules/%v/%s: %s", index, schemaKeyUsed, err.Error())
		}
	}
	return nil
}
