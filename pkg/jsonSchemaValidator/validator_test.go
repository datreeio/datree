package jsonSchemaValidator

import (
	_ "embed"
	"encoding/json"
	"strings"
	"testing"

	extensions "github.com/datreeio/datree/pkg/jsonSchemaValidator/extensions"
	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v5"

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

//go:embed test_fixtures/invalid-rego-definition.json
var invalidRegoDefinitionJson string

//go:embed test_fixtures/valid-rego-definition.json
var validRegoDefinitionJson string

//go:embed test_fixtures/rego-rule-fail.yaml
var regoRuleFail string

//go:embed test_fixtures/rego-rule-pass.yaml
var regoRulePass string

func TestRegoDefinitionCustomKey(t *testing.T) {
	t.Run("invalidSchema", func(t *testing.T) {
		c := jsonschema.NewCompiler()
		c.RegisterExtension(extensions.RegoDefinitionCustomKey, extensions.CustomKeyRegoRule, extensions.CustomKeyRegoDefinitionCompiler{})
		err := c.AddResource("test.json", strings.NewReader(invalidRegoDefinitionJson))
		if err != nil {
			t.Fatal(err)
		}
		_, err = c.Compile("test.json")
		if err == nil {
			t.Fatal("error expected")
		}
		assert.Contains(t, err.Error(), "regoDefinition must be an object of type RegoDefinition json: cannot unmarshal number into Go struct field RegoDefinition.code of type string")
	})
	t.Run("validSchema", func(t *testing.T) {
		c := jsonschema.NewCompiler()
		c.RegisterExtension(extensions.RegoDefinitionCustomKey, extensions.CustomKeyRegoRule, extensions.CustomKeyRegoDefinitionCompiler{})
		if err := c.AddResource("test.json", strings.NewReader(validRegoDefinitionJson)); err != nil {
			t.Fatal(err)
		}
		schema, err := c.Compile("test.json")
		if err != nil {
			t.Fatal(err)
		}
		t.Run("validInstance", func(t *testing.T) {
			jsonYamlContent, err := getInterfaceFromYamlContext(regoRulePass)
			if err != nil {
				t.Fatal(err)
			}

			if err := schema.Validate(jsonYamlContent); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("invalidInstance", func(t *testing.T) {
			jsonYamlContent, err := getInterfaceFromYamlContext(regoRuleFail)
			if err != nil {
				t.Fatal(err)
			}
			if err := schema.Validate(jsonYamlContent); err == nil {
				t.Fatal("validation must fail")
			} else {
				t.Logf("%#v", err)
				if !strings.Contains(err.(*jsonschema.ValidationError).GoString(), "doesn't validate") {
					t.Fatal("validation error expected to contain regoDefinition message")
				}
			}
		})
	})
}

func getInterfaceFromYamlContext(yamlContent string) (interface{}, error) {
	var jsonYamlContent interface{}
	regoRuleFailsYamlBytes, _ := yaml.YAMLToJSON([]byte(yamlContent))
	err := json.Unmarshal(regoRuleFailsYamlBytes, &jsonYamlContent)
	if err != nil {
		return nil, err
	}
	return jsonYamlContent, nil
}
