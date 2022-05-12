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

	return validatePoliciesContent(jsonContent)
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

	err := checkIdentifierUniquenessInPolicy(policies)
	if err != nil {
		return err
	}

	err = checkIdentifierUniqueness(customRules)
	return err
}

func checkIdentifierUniquenessInPolicy(policies []*cliClient.Policy) error {
	for _, policy := range policies {
		propertyValuesExistenceMap := make(map[string]bool)
		rules := policy.Rules
		for _, rule := range rules {
			identifier := rule.Identifier
			if propertyValuesExistenceMap[identifier] {
				return fmt.Errorf("invalid policy file: identifier \"%s\" used atleast two times in policy \"%s\"", identifier, policy.Name)
			}
			propertyValuesExistenceMap[identifier] = true
		}
	}
	return nil
}

func checkIdentifierUniqueness(customRules []*cliClient.CustomRule) error {
	defaultRules, err := defaultRules.GetDefaultRules()
	if err != nil {
		return err
	}
	propertyValuesExistenceMap := make(map[string]bool)

	for _, item := range customRules {
		identifier := item.Identifier

		if propertyValuesExistenceMap[identifier] {
			return fmt.Errorf("invalid policy file: identifier \"%s\" is used in atleast two rules", identifier)
		}

		propertyValuesExistenceMap[identifier] = true
	}

	for _, item := range defaultRules.Rules {
		identifier := item.UniqueName

		if propertyValuesExistenceMap[identifier] {
			return fmt.Errorf("invalid policy file: a default rule with same identifier \"%s\" already exists", identifier)
		}
		propertyValuesExistenceMap[identifier] = true
	}

	return nil
}

func validateSingleDefaultPolicy(policies []*cliClient.Policy) error {
	sawDefault := false
	for _, policy := range policies {

		if policy.IsDefault {
			if sawDefault {
				return fmt.Errorf("invalid policy file: multiple policies are set to default")
			}
			sawDefault = true
		}
	}
	return nil
}

func validateSchemaField(customRules []*cliClient.CustomRule) error {
	for _, rule := range customRules {
		var err error
		var jsonContent string
		identifier := rule.Identifier
		if rule.Schema != nil {
			var content []byte
			schema := rule.Schema
			content, err = json.Marshal(schema)
			if err != nil {
				return err
			}
			jsonContent = string(content)
		} else {
			jsonContent = rule.JsonSchema
		}

		if jsonContent == "" {
			return fmt.Errorf("invalid policy file: rule identifier \"%s\" has no associated schema", identifier)
		}
		schemaLoader := gojsonschema.NewStringLoader(jsonContent)
		_, err = gojsonschema.NewSchemaLoader().Compile(schemaLoader)
		if err != nil {
			return fmt.Errorf("invalid policy file: rule \"%s\" has schema error: %s", rule.Identifier, err.Error())
		}
	}
	return nil
}
