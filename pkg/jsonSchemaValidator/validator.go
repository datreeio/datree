package jsonSchemaValidator

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/xeipuuv/gojsonschema"

	"k8s.io/apimachinery/pkg/api/resource"
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

func (jsv *JSONSchemaValidator) ValidateYamlSchemaNew(schemaContent string, yamlContent string) ([]jsonschema.Detailed, error) {
	jsonSchema, _ := yaml.YAMLToJSON([]byte(schemaContent))
	return jsv.NewValidate(string(jsonSchema), yamlContent)
}

func (jsv *JSONSchemaValidator) Validate(schemaContent string, yamlContent string) (*Result, error) {
	jsonContent, _ := yaml.YAMLToJSON([]byte(yamlContent))

	schemaLoader := gojsonschema.NewStringLoader(schemaContent)
	documentLoader := gojsonschema.NewStringLoader(string(jsonContent))

	return gojsonschema.Validate(schemaLoader, documentLoader)
}

func (jsv *JSONSchemaValidator) NewValidate(schemaContent string, yamlContent string) ([]jsonschema.Detailed, error) {
	var m interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &m)
	if err != nil {
		panic(err)
	}
	m, err = toStringKeys(m)
	if err != nil {
		panic(err)
	}

	compiler := jsonschema.NewCompiler()

	if err := compiler.AddResource("schema.json", strings.NewReader(schemaContent)); err != nil {
		panic(err)
	}

	var resourceMinimum = jsonschema.MustCompileString("resourceMinimum.json", `{
	"properties" : {
		"resourceMinimum": {
			"type": "string"
		}
	}
}`)

	var resourceMaximum = jsonschema.MustCompileString("resourceMaximum.json", `{
	"properties" : {
		"resourceMaximum": {
			"type": "string"
		}
	}
}`)

	compiler.RegisterExtension("resourceMinimum", resourceMinimum, resourceMinimumCompiler{})
	compiler.RegisterExtension("resourceMaximum", resourceMaximum, resourceMaximumCompiler{})

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		panic(err)
	}

	err = schema.Validate(m)

	if err != nil {
		if ve, ok := err.(*jsonschema.ValidationError); ok {
			out := ve.DetailedOutput()
			res := getErrors(out.Errors)
			return res, nil
		} else {
			fmt.Fprintf(os.Stderr, "validation failed: %v\n", err)
			return nil, err
		}
	}
	return nil, err
}

/*
This package accepts only map[string]interface{}, so we need to manually convert them to
map[string]interface{}
*/
func toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
			}
			m[k], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}

type resourceMinimumCompiler struct{}

type resourceMaximumCompiler struct{}

func (resourceMinimumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMinimum, ok := m["resourceMinimum"]; ok {
		n := resourceMinimum.(string)
		return resourceMinimumSchema(n), nil
	}

	// nothing to compile, return nil
	return nil, nil
}
func (resourceMaximumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMinimum, ok := m["resourceMaximum"]; ok {
		n := resourceMinimum.(string)
		return resourceMaximumSchema(n), nil
	}

	// nothing to compile, return nil
	return nil, nil
}

type resourceMinimumSchema string
type resourceMaximumSchema string

func (s resourceMinimumSchema) Validate(ctx jsonschema.ValidationContext, v interface{}) error {
	switch v.(type) {
	case string:
		resourceMinimumStr := string(s)
		rmDataValueParsedQ, err := resource.ParseQuantity(v.(string))
		if err != nil {
			if err != nil {
				return ctx.Error("resourceMinimum", "failed parsing value %v", v)
			}
		}

		rmSchemaValueParsedQ, err := resource.ParseQuantity(resourceMinimumStr)
		if err != nil {
			return ctx.Error("resourceMinimum", "failed parsing value %v", resourceMinimumStr)
		}

		rmDecStr := rmDataValueParsedQ.AsDec().String()
		rmRfDecStr := rmSchemaValueParsedQ.AsDec().String()

		resourceMinimumSchemaVal, _ := strconv.ParseFloat(rmDecStr, 64)
		resourceMinimumDataVal, _ := strconv.ParseFloat(rmRfDecStr, 64)

		if resourceMinimumDataVal > resourceMinimumSchemaVal {
			return ctx.Error("resourceMinimum", "%v is lower then resourceMinimum %v", v, resourceMinimumStr)
		}
		return nil
	default:
		return nil
	}
}

func (s resourceMaximumSchema) Validate(ctx jsonschema.ValidationContext, v interface{}) error {
	switch v.(type) {
	case string:
		resourceMaximumStr := string(s)
		rmDataValueParsedQ, err := resource.ParseQuantity(v.(string))
		if err != nil {
			fmt.Println(err.Error())
			return ctx.Error("resourceMaximum", "failed parsing value %v", v)
		}

		rmSchemaValueParsedQ, err := resource.ParseQuantity(resourceMaximumStr)
		if err != nil {
			return ctx.Error("resourceMaximum", "failed parsing value %v", resourceMaximumStr)
		}

		rmDecStr := rmDataValueParsedQ.AsDec().String()
		rmRfDecStr := rmSchemaValueParsedQ.AsDec().String()

		resourceMaximumSchemaVal, _ := strconv.ParseFloat(rmDecStr, 64)
		resourceMaximumDataVal, _ := strconv.ParseFloat(rmRfDecStr, 64)

		if resourceMaximumDataVal < resourceMaximumSchemaVal {
			return ctx.Error("resourceMaximum", "%v is greater then resourceMaximum %v", v, resourceMaximumStr)
		}
		return nil
	default:
		return nil
	}
}

func getErrors(errors []jsonschema.Detailed) []jsonschema.Detailed {
	if len(errors) > 0 {
		for _, err := range errors {
			if len(err.Errors) > 0 {
				return getErrors(err.Errors)
			} else {
				return errors
			}
		}
	} else {
		return errors
	}
	return nil
}
