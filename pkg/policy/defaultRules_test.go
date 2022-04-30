package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

var defaultRulesFileContent string
var defaultRulesSchemaFileContent string

const defaultRulesYamlPath = "./defaultRules.yaml"
const defaultRulesJSONSchemaPath = "./defaultRulesSchema.json"

type DefaultRules struct {
	ApiVersion string                  `json:"apiVersion"`
	Rules      []DefaultRuleDefinition `json:"rules"`
}

func TestMain(m *testing.M) {
	defaultRulesFile, err := getFileFromPath(defaultRulesYamlPath)
	if err != nil {
		panic(err)
	}

	defaultRulesFileContent = defaultRulesFile

	defaultRulesSchemaFile, err := getFileFromPath(defaultRulesJSONSchemaPath)
	if err != nil {
		panic(err)
	}

	defaultRulesSchemaFileContent = defaultRulesSchemaFile

	testsRunResult := m.Run()
	os.Exit(testsRunResult)
}

func TestDefaultRulesFileExists(t *testing.T) {
	_, fileReadError := getFileFromPath(defaultRulesYamlPath)

	assert.Nil(t, fileReadError)
}

func TestDefaultRulesFileFitsJSONSchema(t *testing.T) {
	validationError := validateYamlUsingJSONSchema(defaultRulesYamlPath, defaultRulesSchemaFileContent)

	assert.Nil(t, validationError)
}

func TestDefaultRulesHasUniqueNamesInRules(t *testing.T) {
	defaultRulesMap, conversionToMapError := convertYamlFileToMap(defaultRulesFileContent)
	assert.Nil(t, conversionToMapError)

	uniquenessValidationError := validateUniqueStringValuesInRulesForProperty("UniqueName", defaultRulesMap.Rules)
	assert.Nil(t, uniquenessValidationError)
}

func TestDefaultRulesHasUniqueIDsInRules(t *testing.T) {
	defaultRulesMap, conversionToMapError := convertYamlFileToMap(defaultRulesFileContent)
	assert.Nil(t, conversionToMapError)

	uniquenessValidationError := validateUniqueIntValuesInRulesForProperty("ID", defaultRulesMap.Rules)

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
		panic(errors.New("error occurred while validating yaml file against schema"))
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
		propertyValue := getStringValueOfRuleProperty(item, propertyName)

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("property %s has duplicate value %s", propertyName, propertyValue)
		}

		propertyValuesExistenceMap[propertyValue] = true
	}

	return nil
}

func getStringValueOfRuleProperty(rule DefaultRuleDefinition, propertyName string) string {
	itemValue := reflect.ValueOf(rule)
	propertyValue := reflect.Indirect(itemValue).FieldByName(propertyName).String()

	return propertyValue
}

func validateUniqueIntValuesInRulesForProperty(propertyName string, rules []DefaultRuleDefinition) error {
	propertyValuesExistenceMap := make(map[int64]bool)

	for _, item := range rules {
		propertyValue := getIntValueOfRuleProperty(item, propertyName)

		if propertyValuesExistenceMap[propertyValue] {
			return fmt.Errorf("property %s has duplicate value %d", propertyName, propertyValue)
		}

		propertyValuesExistenceMap[propertyValue] = true
	}

	return nil
}

func getIntValueOfRuleProperty(rule DefaultRuleDefinition, propertyName string) int64 {
	itemValue := reflect.ValueOf(rule)
	propertyValue := reflect.Indirect(itemValue).FieldByName(propertyName).Int()

	return propertyValue
}
