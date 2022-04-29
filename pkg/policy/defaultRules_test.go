package policy

import (
	_ "embed"
	"encoding/json"
	"errors"
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

	defaultRulesFileContentJSON, conversionToJSONError := convertYamlFileToMap(defaultRulesFileContent)
	assert.Nil(t, conversionToJSONError)

	uniquenessValidationError := validateUniqueStringPropertyValuesInArray("uniqueName", defaultRulesFileContentJSON["rules"])
	assert.Nil(t, uniquenessValidationError)
}

func TestDefaultRulesHasUniqueIDsInRules(t *testing.T) {
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)

	defaultRulesFileContentJSON, conversionToJSONError := convertYamlFileToMap(defaultRulesFileContent)
	assert.Nil(t, conversionToJSONError)

	uniquenessValidationError := validateUniqueFloat64PropertyValuesInArray("id", defaultRulesFileContentJSON["rules"])

	assert.Nil(t, uniquenessValidationError)
}

func getFileFromPath(path string) (string, error) {
	fileReader := fileReader.CreateFileReader(nil)
	fileContent, readFileError := fileReader.ReadFileContent(path)

	if readFileError != nil {
		return "", readFileError
	}

	return fileContent, nil
}

func validateYamlUsingJSONSchema(yamlFilePath string, schema string) error {
	fileContent, _ := getFileFromPath(yamlFilePath)
	jsonSchemaValidator := jsonSchemaValidator.New()
	schemaValidationResult, schemaValidationError := jsonSchemaValidator.ValidateYamlSchema(schema, fileContent)

	if schemaValidationError != nil {
		panic(errors.New("can't validate yaml file using json schema"))
	}

	if schemaValidationResult != nil {
		validationErrors := fmt.Errorf("Received validation errors for %s:\n", yamlFilePath)

		for _, validationError := range schemaValidationResult {
			validationErrors = fmt.Errorf("%s\n%s", validationErrors, validationError.Error)
		}

		return validationErrors
	}

	return nil
}

func convertYamlFileToMap(yamlFileContent string) (map[string][]interface{}, error) {
	yamlFileContentRawJSON, yamlParseError := yaml.YAMLToJSON([]byte(yamlFileContent))

	if yamlParseError != nil {
		return nil, yamlParseError
	}

	var yamlFileContentJSON map[string][]interface{}
	jsonUnmarshallingError := json.Unmarshal(yamlFileContentRawJSON, &yamlFileContentJSON)

	var jsonUnmarshallingFailed = jsonUnmarshallingError != nil && reflect.TypeOf(yamlFileContentJSON) != reflect.TypeOf(map[string][]interface{}{})
	if jsonUnmarshallingFailed {
		return nil, jsonUnmarshallingError
	}

	return yamlFileContentJSON, nil
}

func validateUniqueStringPropertyValuesInArray(propertyName string, array []interface{}) error {
	propertyValuesExistenceMap := make(map[string]bool)

	for _, item := range array {
		itemObject := item.(map[string]interface{})
		propertyValue := itemObject[propertyName].(string)

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("property %s has duplicate value %s", propertyName, propertyValue)
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
