package policy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

//go:embed defaultRulesSchema.json
var defaultRulesSchemaContent string

const defaultRulesYamlPath = "./defaultRules.yaml"

func TestDefaultRulesFileExists(t *testing.T) {
	_, fileReadError := getFileFromPath(defaultRulesYamlPath)

	assert.Nil(t, fileReadError)
}

func TestDefaultRulesFileFitsJSONSchema(t *testing.T) {
	validationError := validateYamlUsingJSONSchema(defaultRulesYamlPath, defaultRulesSchemaContent)

	assert.Nil(t, validationError)
}

func TestDefaultRulesHasUniqueNamesInRules(t *testing.T) {
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)

	defaultRulesFileContentJSON, conversionToJSONError := convertYamlFileContentToJSON(defaultRulesFileContent)
	assert.Nil(t, conversionToJSONError)

	uniquenessValidationError := validateUniqueStringPropertyValuesInArray("uniqueName", defaultRulesFileContentJSON["rules"])
	assert.Nil(t, uniquenessValidationError)
}

func TestDefaultRulesHasUniqueIDsInRules(t *testing.T) {
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)

	defaultRulesFileContentJSON, conversionToJSONError := convertYamlFileContentToJSON(defaultRulesFileContent)
	assert.Nil(t, conversionToJSONError)

	uniquenessValidationError := validateUniqueFloat64PropertyValuesInArray("id", defaultRulesFileContentJSON["rules"])

	assert.Nil(t, uniquenessValidationError)
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
	yamlFileContentRawJSON, yamlParseError := yaml.YAMLToJSON([]byte(yamlFileContent))

	if yamlParseError != nil {
		return map[string][]interface{}{}, yamlParseError
	}

	var yamlFileContentJSON map[string][]interface{}
	jsonMarshallingError := json.Unmarshal(yamlFileContentRawJSON, &yamlFileContentJSON)

	var jsonMarshallingNotFailed = jsonMarshallingError != nil && reflect.TypeOf(yamlFileContentJSON) != reflect.TypeOf(map[string][]interface{}{})
	if jsonMarshallingNotFailed {
		return map[string][]interface{}{}, jsonMarshallingError
	}

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
