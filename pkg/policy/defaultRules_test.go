package policy

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

//go:embed defaultRulesSchema.json
var defaultRulesSchemaContent string

const defaultRulesYamlPath = "./defaultRules.yaml"

type DefaultRules struct {
	ApiVersion string                  `json:"apiVersion"`
	Rules      []DefaultRuleDefinition `json:"rules"`
}

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

	uniquenessValidationError := validateUniqueStringValuesInRulesForProperty("UniqueName", defaultRulesFileContentJSON.Rules)
	assert.Nil(t, uniquenessValidationError)
}

func TestDefaultRulesHasUniqueIDsInRules(t *testing.T) {
	defaultRulesFileContent, _ := getFileFromPath(defaultRulesYamlPath)

	defaultRulesFileContentJSON, conversionToJSONError := convertYamlFileToMap(defaultRulesFileContent)
	assert.Nil(t, conversionToJSONError)

	uniquenessValidationError := validateUniqueFloatValuesInRulesForProperty("ID", defaultRulesFileContentJSON.Rules)

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

func convertYamlFileToMap(yamlFileContent string) (DefaultRules, error) {
	yamlFileContentRawJSON, yamlParseError := yaml.YAMLToJSON([]byte(yamlFileContent))

	if yamlParseError != nil {
		return DefaultRules{}, yamlParseError
	}

	var yamlFileContentJSON DefaultRules
	jsonUnmarshallingError := json.Unmarshal(yamlFileContentRawJSON, &yamlFileContentJSON)

	if jsonUnmarshallingError != nil {
		return DefaultRules{}, jsonUnmarshallingError
	}

	return yamlFileContentJSON, nil
}

func validateUniqueStringValuesInRulesForProperty(propertyName string, rules []DefaultRuleDefinition) error {
	propertyValuesExistenceMap := make(map[string]bool)

	for _, item := range rules {
		r := reflect.ValueOf(item)
		propertyValue := reflect.Indirect(r).FieldByName(propertyName).String()

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("property %s has duplicate value %s", propertyName, propertyValue)
		}

		propertyValuesExistenceMap[propertyValue] = true
	}

	return nil
}

func validateUniqueFloatValuesInRulesForProperty(propertyName string, rules []DefaultRuleDefinition) error {
	propertyValuesExistenceMap := make(map[int64]bool)

	for _, item := range rules {
		r := reflect.ValueOf(item)
		propertyValue := reflect.Indirect(r).FieldByName(propertyName).Int()

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("property %s has duplicate value %d", propertyName, propertyValue)
		}

		propertyValuesExistenceMap[propertyValue] = true
	}

	return nil
}
