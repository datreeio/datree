package policy

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/stretchr/testify/assert"
)

//go:embed defaultRulesSchema.json
var defaultRulesSchemaContent string

func TestDefaultRulesFileExists(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"
	_, err := getFileFromPath(defaultRulesYamlPath)

	assert.Nil(t, err)
}

func TestDefaultRulesFileFitsSchema(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"

	err := validateYamlUsingJSONSchema(defaultRulesYamlPath, defaultRulesSchemaContent)

	assert.Nil(t, err)
}

func getFileFromPath(path string) (string, error) {
	fileReader := fileReader.CreateFileReader(nil)
	fileContent, err := fileReader.ReadFileContent(path)

	if err != nil {
		return "", err
	}

	return fileContent, nil
}

func validateYamlUsingJSONSchema(yamlFilePath string, schema string) error {
	fileContent, _ := getFileFromPath(yamlFilePath)
	jsonSchemaValidator := jsonSchemaValidator.New()
	result, err := jsonSchemaValidator.ValidateYamlSchema(schema, fileContent)

	if err != nil {
		fmt.Errorf("Failed to validate %s:\n", yamlFilePath)

		return err
	}

	if !result.Valid() {
		validationErrors := fmt.Errorf("Received validation errors for %s:\n", yamlFilePath)

		for _, validationError := range result.Errors() {
			validationErrors = fmt.Errorf("%s\n%s", validationErrors, validationError)
		}

		return validationErrors
	}

	return nil
}
