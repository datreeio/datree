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

	if errorsResult != nil {
		validationErrors := fmt.Errorf("found errors in policies file %s:", policyYamlPath)

		for _, validationError := range errorsResult {
			validationErrors = fmt.Errorf("%s\n(root)%s: %s", validationErrors, validationError.InstanceLocation, validationError.Error)
		}

		return validationErrors
	}

	return validatePoliciesContent(jsonContent, policyYamlPath)
}

func validatePoliciesContent(content []byte, policyYamlPath string) error {
	// unmarshal the content
	var schema *cliClient.EvaluationPrerunPolicies
	err := json.Unmarshal(content, &schema)
	if err != nil {
		return err
	}

	// validate that exactly one policy is set to default
	err = validateSingleDefaultPolicy(schema.Policies, policyYamlPath)
	if err != nil {
		return err
	}

	// validate if the policy file has any identifier related issues
	err = validateIdentifier(schema.Policies, schema.CustomRules, policyYamlPath)
	if err != nil {
		return err
	}

	// validate the schema of each rule
	err = validateSchemaField(schema.CustomRules, policyYamlPath)
	return err
}

func validateIdentifier(policies []*cliClient.Policy, customRules []*cliClient.CustomRule, policyYamlPath string) error {

	err := checkIdentifierInPolicy(policies, customRules, policyYamlPath)
	if err != nil {
		return err
	}

	err = checkIdentifierUniqueness(customRules, policyYamlPath)
	return err
}

func checkIdentifierInPolicy(policies []*cliClient.Policy, customRules []*cliClient.CustomRule, policyYamlPath string) error {
	err := checkIdentifierUniquenessInPolicy(policies, policyYamlPath)
	if err != nil {
		return err
	}

	err = checkIdentifierExistence(policies, customRules, policyYamlPath)
	if err != nil {
		return err
	}
	return nil
}

func checkIdentifierExistence(policies []*cliClient.Policy, customRules []*cliClient.CustomRule, policyYamlPath string) error {
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
				return fmt.Errorf("found errors in policies file %s:\n(root)/policies/%d/rules: identifier \"%s\" is neither custom nor default", policyYamlPath, index, identifier)
			}
		}
	}

	return nil
}

func checkIdentifierUniquenessInPolicy(policies []*cliClient.Policy, policyYamlPath string) error {
	for index, policy := range policies {
		propertyValuesExistenceMap := make(map[string]bool)
		rules := policy.Rules
		for _, rule := range rules {
			identifier := rule.Identifier
			if propertyValuesExistenceMap[identifier] {
				return fmt.Errorf("found errors in policies file %s:\n(root)/policies/%d/rules: identifier \"%s\" is used more than once in policy", policyYamlPath, index, identifier)
			}
			propertyValuesExistenceMap[identifier] = true
		}
	}
	return nil
}

func checkIdentifierUniqueness(customRules []*cliClient.CustomRule, policyYamlPath string) error {
	defaultRules, err := defaultRules.GetDefaultRules()
	if err != nil {
		return err
	}
	propertyValuesExistenceMap := make(map[string]bool)

	for _, item := range customRules {
		identifier := item.Identifier

		if propertyValuesExistenceMap[identifier] {
			return fmt.Errorf("found errors in policies file %s:\n(root)/customRules: identifier \"%s\" is used in more than one custom rule", policyYamlPath, identifier)
		}

		propertyValuesExistenceMap[identifier] = true
	}

	for _, item := range defaultRules.Rules {
		identifier := item.UniqueName

		if propertyValuesExistenceMap[identifier] {
			return fmt.Errorf("found errors in policies file %s:\n(root)/customRules: a default rule with same identifier \"%s\" already exists", policyYamlPath, identifier)
		}
		propertyValuesExistenceMap[identifier] = true
	}

	return nil
}

func validateSingleDefaultPolicy(policies []*cliClient.Policy, policyYamlPath string) error {
	sawDefault := false
	for _, policy := range policies {

		if policy.IsDefault {
			if sawDefault {
				return fmt.Errorf("found errors in policies file %s:\n(root)/policies: Should have exactly one policy set as default", policyYamlPath)
			}
			sawDefault = true
		}
	}

	if !sawDefault {
		return fmt.Errorf("found errors in policies file %s:\n(root)/policies: Should have exactly one policy set as default", policyYamlPath)
	}
	return nil
}

func validateSchemaField(customRules []*cliClient.CustomRule, policyYamlPath string) error {
	for index, rule := range customRules {
		var err error
		var jsonContent string
		if rule.Schema != nil && rule.JsonSchema != "" {
			return fmt.Errorf("found errors in policies file %s:\n(root)/customRules/%d: Exactly one of [schema,jsonSchema] should be defined per custom rule", policyYamlPath, index)
		}
		if rule.Schema != nil {
			var content []byte
			schema := rule.Schema
			content, err = json.Marshal(schema)
			if err != nil {
				return fmt.Errorf("found errors in policies file %s:\n(root)/customRules/%d: %s", policyYamlPath, index, err.Error())
			}
			jsonContent = string(content)
		} else {
			jsonContent = rule.JsonSchema
		}

		if jsonContent == "" {
			return fmt.Errorf("found errors in policies file %s:\n(root)/customRules/%d: Exactly one of [schema,jsonSchema] should be defined per custom rule", policyYamlPath, index)
		}
		schemaLoader := gojsonschema.NewStringLoader(jsonContent)
		_, err = gojsonschema.NewSchemaLoader().Compile(schemaLoader)
		if err != nil {
			return fmt.Errorf("found errors in policies file %s:\n(root)/customRules/%v/schema: %s", policyYamlPath, index, err.Error())
		}
	}
	return nil
}
