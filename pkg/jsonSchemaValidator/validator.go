package jsonSchemaValidator

import (
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type JSONSchemaValidator struct {
}

func New() *JSONSchemaValidator {
	return &JSONSchemaValidator{}
}

type Result = gojsonschema.Result

func (jsv *JSONSchemaValidator) ValidateYamlSchema(schemaContent string, yamlContent string) (*Result, error) {
	jsonSchema, _ := yaml.YAMLToJSON([]byte(schemaContent))
	return jsv.Validate(string(jsonSchema), yamlContent)
}

func (jsv *JSONSchemaValidator) Validate(schemaContent string, yamlContent string) (*Result, error) {
	jsonContent, _ := yaml.YAMLToJSON([]byte(yamlContent))

	schemaLoader := gojsonschema.NewStringLoader(schemaContent)
	documentLoader := gojsonschema.NewStringLoader(string(jsonContent))

	return gojsonschema.Validate(schemaLoader, documentLoader)
}
