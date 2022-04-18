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
	_, fileReadError := getFileFromPath(defaultRulesYamlPath)

	assert.Nil(t, fileReadError)
}

func TestDefaultRulesFileFitsJSONSchema(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"

	validationError := validateYamlUsingJSONSchema(defaultRulesYamlPath, defaultRulesSchemaContent)

	assert.Nil(t, validationError)
}

func TestDefaultRulesHasUniqueNamesInRules(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)

	defaultRulesFileContentJSON, jsonParseError := convertYamlFileContentToJSON(defaultRulesFileContent)
	assert.Nil(t, jsonParseError)

	validationError := validateUniqueStringPropertyValuesInArray("uniqueName", defaultRulesFileContentJSON["rules"])
	assert.Nil(t, validationError)
}

func TestDefaultRulesHasUniqueIDsInRules(t *testing.T) {
	defaultRulesYamlPath := "./defaultRules.yaml"
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)

	defaultRulesFileContentJSON, jsonParseError := convertYamlFileContentToJSON(defaultRulesFileContent)
	assert.Nil(t, jsonParseError)

	validationError := validateUniqueFloat64PropertyValuesInArray("id", defaultRulesFileContentJSON["rules"])

	assert.Nil(t, validationError)
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

func convertYamlFileContentToJSON(yamlFileContent string) (map[string][]interface{}, error) {
	yamlFileContentRawJSON, err := yaml.YAMLToJSON([]byte(yamlFileContent))

	if err != nil {
		return map[string][]interface{}{}, err
	}

	var yamlFileContentJSON map[string][]interface{}
	json.Unmarshal(yamlFileContentRawJSON, &yamlFileContentJSON)

	return yamlFileContentJSON, nil
}

func validateUniqueStringPropertyValuesInArray(propertyName string, array []interface{}) error {
	propertyValuesExistenceMap := make(map[string]bool)

	for _, item := range array {
		itemObject := item.(map[string]interface{})
		propertyValue := itemObject[propertyName].(string)

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("Property %s has duplicate value %s", propertyName, propertyValue)
		}

		propertyValuesExistenceMap[propertyValue] = true
	}

	return nil
}

func validateUniqueFloat64PropertyValuesInArray(propertyName string, array []interface{}) error {
	propertyValuesExistenceMap := make(map[float64]bool)

	for _, item := range array {
		itemObject := item.(map[string]interface{})
		propertyValue := itemObject[propertyName].(float64)

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("Property %s has duplicate value %f", propertyName, propertyValue)
		}

		propertyValuesExistenceMap[propertyValue] = true
	}

	return nil
}
