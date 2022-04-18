package policy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

//go:embed defaultRulesSchema.json
var defaultRulesSchemaContent string

func TestDefaultRulesFileExists(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"
	_, err := getFileFromPath(defaultRulesYamlPath)

	assert.Nil(t, err)
}

func TestDefaultRulesFileFitsJSONSchema(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"

	err := validateYamlUsingJSONSchema(defaultRulesYamlPath, defaultRulesSchemaContent)

	assert.Nil(t, err)
}

func TestDefaultRulesHasUniqueNamesInRules(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)
	defaultRulesFileContentRawJSON, _ := yaml.YAMLToJSON([]byte(defaultRulesFileContent))

	var defaultRulesFileContentJSON map[string][]interface{}
	json.Unmarshal(defaultRulesFileContentRawJSON, &defaultRulesFileContentJSON)

	err := validateUniqueStringPropertyValuesInArray("uniqueName", defaultRulesFileContentJSON["rules"])

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

func validateUniqueStringPropertyValuesInArray(propertyName string, array []interface{}) error {
	propertyValues := make(map[string]bool)

	for _, object := range array {
		objectMap := object.(map[string]interface{})

		propertyValue := objectMap[propertyName].(string)

		if propertyValues[propertyValue] {
			return fmt.Errorf("Property %s has duplicate value %s", propertyName, propertyValue)
		}

		propertyValues[propertyValue] = true
	}

	return nil
}
