package jsonSchemaValidator

import (
	"testing"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

const yamlFilesPath = "../../internal/fixtures/policyAsCode/custom-keys"

func TestValidateCustomKeysFail(t *testing.T) {
	fileReader := fileReader.CreateFileReader(nil)

	resourceJson, err := fileReader.ReadFileContent(yamlFilesPath + "/fail-yaml-file.yaml")
	if err != nil {
		panic(err)
	}

	customRuleSchemaJson, err := fileReader.ReadFileContent(yamlFilesPath + "/custom-rule.yaml")
	if err != nil {
		panic(err)
	}

	jsonSchemaValidator := New()

	resourceYaml, _ := yaml.JSONToYAML([]byte(resourceJson))
	customRuleYaml, _ := yaml.JSONToYAML([]byte(customRuleSchemaJson))

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(string(customRuleYaml), string(resourceYaml))

	assert.GreaterOrEqual(t, len(errorsResult), 1)
}

func TestValidateCustomKeysPass(t *testing.T) {
	fileReader := fileReader.CreateFileReader(nil)

	resourceJson, err := fileReader.ReadFileContent(yamlFilesPath + "/pass-yaml-file.yaml")
	if err != nil {
		panic(err)
	}

	customRuleSchemaJson, err := fileReader.ReadFileContent(yamlFilesPath + "/custom-rule.yaml")
	if err != nil {
		panic(err)
	}

	jsonSchemaValidator := New()

	resourceYaml, _ := yaml.JSONToYAML([]byte(resourceJson))
	customRuleYaml, _ := yaml.JSONToYAML([]byte(customRuleSchemaJson))

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(string(customRuleYaml), string(resourceYaml))

	assert.Empty(t, errorsResult)
}
