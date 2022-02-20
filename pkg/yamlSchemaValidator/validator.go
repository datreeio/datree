package yamlSchemaValidator

import (
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type YamlSchemaValidator struct {
}

func New() *YamlSchemaValidator {
	return &YamlSchemaValidator{}
}

type Result = gojsonschema.Result

func (jsv *YamlSchemaValidator) Validate(schemaContent string, yamlContent string) (*Result, error) {

	jsonSchema, _ := yaml.YAMLToJSON([]byte(schemaContent))

	json, _ := yaml.YAMLToJSON([]byte(yamlContent))

	schemaLoader := gojsonschema.NewStringLoader(string(jsonSchema))
	documentLoader := gojsonschema.NewStringLoader(string(json))

	return gojsonschema.Validate(schemaLoader, documentLoader)
}
