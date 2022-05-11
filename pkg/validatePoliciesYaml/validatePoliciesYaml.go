package validatePoliciesYaml

import (
	_ "embed"
	"fmt"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
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

	return nil
}
