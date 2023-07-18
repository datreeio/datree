package jsonSchemaValidator

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	extensions "github.com/datreeio/datree/pkg/jsonSchemaValidator/extensions"
	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/stretchr/testify/assert"
)

const yamlFilesPath = "../../internal/fixtures/policyAsCode/custom-keys"
const regoYamlFilesPath = "../jsonSchemaValidator/test_fixtures"

func TestValidateResourceMinMaxCustomKeysFail(t *testing.T) {
	failResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			yamlFilesPath+"/fail-yaml-file.yaml",
			yamlFilesPath+"/schema-with-resource-quotas.yaml",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, failResourceYamlFileContent)

	assert.GreaterOrEqual(t, len(errorsResult), 1)
	assert.Equal(t, errorsResult[0].Error, "1G is greater then resourceMaximum 500m")
}

func TestValidateResourceMinMaxCustomKeysPass(t *testing.T) {
	passResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			yamlFilesPath+"/pass-yaml-file.yaml",
			yamlFilesPath+"/schema-with-resource-quotas.yaml",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, passResourceYamlFileContent)

	assert.Empty(t, errorsResult)
}

//go:embed test_fixtures/invalid-rego-definition.json
var invalidRegoDefinitionJson string

//go:embed test_fixtures/invalid-cel-definition.json
var invalidCELDefinitionJson string

//go:embed test_fixtures/valid-rego-definition.json
var validRegoDefinitionJson string

//go:embed test_fixtures/valid-cel-definition.json
var validCELDefinitionJson string

//go:embed test_fixtures/rule-fail.yaml
var ruleFail string

//go:embed test_fixtures/rule-pass.yaml
var rulePass string

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
			jsonYamlContent, err := getInterfaceFromYamlContext(rulePass)
			if err != nil {
				t.Fatal(err)
			}

			if err := schema.Validate(jsonYamlContent); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("invalidInstance", func(t *testing.T) {
			jsonYamlContent, err := getInterfaceFromYamlContext(ruleFail)
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

func TestValidateRegoDefinitionCustomKeyPass(t *testing.T) {
	passResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			regoYamlFilesPath+"/rule-pass.yaml",
			regoYamlFilesPath+"/valid-rego-definition.json",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, passResourceYamlFileContent)

	assert.Empty(t, errorsResult)
}

func TestValidateRegoDefinitionCustomKeyPassDueToResourceNotInConstraint(t *testing.T) {
	passResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			regoYamlFilesPath+"/rego-rule-pass-due-to-not-it-constraint.yaml",
			regoYamlFilesPath+"/valid-rego-definition.json",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, passResourceYamlFileContent)

	assert.Empty(t, errorsResult)
}

func TestValidateRegoDefinitionCustomKeyFail(t *testing.T) {
	failResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			regoYamlFilesPath+"/rule-fail.yaml",
			regoYamlFilesPath+"/valid-rego-definition.json",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, failResourceYamlFileContent)

	assert.GreaterOrEqual(t, len(errorsResult), 1)
	assert.Contains(t, errorsResult[0].Error, "do not match")
}

func TestValidateRegoDefinitionCustomKeyFailDueToRegoCompile(t *testing.T) {
	failResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			regoYamlFilesPath+"/rule-pass.yaml",
			regoYamlFilesPath+"/invalid-rego-definition-code.json",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, failResourceYamlFileContent)

	assert.GreaterOrEqual(t, len(errorsResult), 1)
	assert.Contains(t, errorsResult[0].Error, "can't compile rego code, rego code must have a package")
}

func TestCELDefinitionCustomKey(t *testing.T) {
	t.Run("invalidSchema", func(t *testing.T) {
		c := jsonschema.NewCompiler()
		c.RegisterExtension(extensions.CELDefinitionCustomKey, extensions.CustomKeyCELRule, extensions.CustomKeyCELDefinitionCompiler{})
		err := c.AddResource("test.json", strings.NewReader(invalidCELDefinitionJson))
		if err != nil {
			t.Fatal(err)
		}
		_, err = c.Compile("test.json")
		if err == nil {
			t.Fatal("error expected")
		}
		assert.Contains(t, err.Error(), "CELExpression must be an object of type CELExpression json: cannot unmarshal number into Go struct field CELExpression.expression of type string")
	})
	t.Run("validSchema", func(t *testing.T) {
		c := jsonschema.NewCompiler()
		c.RegisterExtension(extensions.CELDefinitionCustomKey, extensions.CustomKeyCELRule, extensions.CustomKeyCELDefinitionCompiler{})
		if err := c.AddResource("test.json", strings.NewReader(validCELDefinitionJson)); err != nil {
			t.Fatal(err)
		}
		schema, err := c.Compile("test.json")
		if err != nil {
			t.Fatal(err)
		}
		t.Run("validInstance", func(t *testing.T) {
			jsonYamlContent, err := getInterfaceFromYamlContext(rulePass)
			if err != nil {
				t.Fatal(err)
			}

			if err := schema.Validate(jsonYamlContent); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("invalidInstance", func(t *testing.T) {
			jsonYamlContent, err := getInterfaceFromYamlContext(ruleFail)
			if err != nil {
				t.Fatal(err)
			}
			if err := schema.Validate(jsonYamlContent); err == nil {
				t.Fatal("validation must fail")
			} else {
				fmt.Println(err.(*jsonschema.ValidationError).GoString())
				t.Logf("%#v", err)
				if !strings.Contains(err.(*jsonschema.ValidationError).GoString(), "doesn't validate") {
					t.Fatal("validation error expected to contain CELDefinition message")
				}
			}
		})
	})
}

func TestValidateCELDefinitionCustomKeyPass(t *testing.T) {
	passResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			regoYamlFilesPath+"/rule-pass.yaml",
			regoYamlFilesPath+"/valid-cel-definition.json",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, passResourceYamlFileContent)

	assert.Empty(t, errorsResult)
}

func TestValidateCELDefinitionCustomKeyFailDueToCELCompile(t *testing.T) {
	failResourceYamlFileContent, customRuleSchemaYamlFileContent :=
		getResourceAndSchemaYamlContentsAsString(
			regoYamlFilesPath+"/rule-pass.yaml",
			regoYamlFilesPath+"/invalid-cel-definition-expression.json",
		)

	jsonSchemaValidator := New()

	errorsResult, _ := jsonSchemaValidator.ValidateYamlSchema(customRuleSchemaYamlFileContent, failResourceYamlFileContent)

	assert.GreaterOrEqual(t, len(errorsResult), 1)
	assert.Contains(t, errorsResult[0].Error, "cel expression compile error")
}

func getResourceAndSchemaYamlContentsAsString(resourceToValidatePath string, schemaPath string) (string, string) {
	fileReader := fileReader.CreateFileReader(nil)

	resourceYamlFileContent, err := fileReader.ReadFileContent(resourceToValidatePath)
	if err != nil {
		panic(err)
	}

	customRuleSchemaYamlFileContent, err := fileReader.ReadFileContent(schemaPath)
	if err != nil {
		panic(err)
	}

	return resourceYamlFileContent, customRuleSchemaYamlFileContent
}

func getInterfaceFromYamlContext(yamlContent string) (interface{}, error) {
	var jsonYamlContent interface{}
	ruleFailsYamlBytes, _ := yaml.YAMLToJSON([]byte(yamlContent))
	err := json.Unmarshal(ruleFailsYamlBytes, &jsonYamlContent)
	if err != nil {
		return nil, err
	}
	return jsonYamlContent, nil
}
