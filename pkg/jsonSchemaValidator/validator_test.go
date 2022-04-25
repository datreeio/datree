package jsonSchemaValidator

import (
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/stretchr/testify/assert"
)

const yamlFilesPath = "../../internal/fixtures/policyAsCode/custom-keys"

func TestValidateCustomKeysFail(t *testing.T) {
	fileReader := fileReader.CreateFileReader(nil)

	failResourceYamlFileContent, err := fileReader.ReadFileContent(yamlFilesPath + "/fail-yaml-file.yaml")
	if err != nil {
		panic(err)
	}

	customRuleSchemaYamlFileContent, err := fileReader.ReadFileContent(yamlFilesPath + "/schema-with-resource-quotas.yaml")
	if err != nil {
		panic(err)
	}

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, failResourceYamlFileContent)

	assert.GreaterOrEqual(t, len(errorsResult), 1)
	assert.Equal(t, errorsResult[0].Error, "1G is greater then resourceMaximum 500m")
}

func TestValidateCustomKeysPass(t *testing.T) {
	fileReader := fileReader.CreateFileReader(nil)

	passResourceYamlFileContent, err := fileReader.ReadFileContent(yamlFilesPath + "/pass-yaml-file.yaml")
	if err != nil {
		panic(err)
	}

	customRuleSchemaYamlFileContent, err := fileReader.ReadFileContent(yamlFilesPath + "/schema-with-resource-quotas.yaml")
	if err != nil {
		panic(err)
	}

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, passResourceYamlFileContent)

	assert.Empty(t, errorsResult)
}
