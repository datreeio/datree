package jsonSchemaValidator

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"k8s.io/apimachinery/pkg/api/resource"
)

type JSONSchemaValidator struct {
}

func New() *JSONSchemaValidator {
	return &JSONSchemaValidator{}
}

type resourceMinimumCompiler struct{}
type resourceMaximumCompiler struct{}

type resourceMinimumSchema string
type resourceMaximumSchema string

func (jsv *JSONSchemaValidator) ValidateYamlSchema(schemaContent string, yamlContent string) ([]jsonschema.Detailed, error) {
	jsonSchema, _ := yaml.YAMLToJSON([]byte(schemaContent))
	return jsv.Validate(string(jsonSchema), yamlContent)
}

func (jsv *JSONSchemaValidator) Validate(schemaContent string, yamlContent string) ([]jsonschema.Detailed, error) {
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
	compiler.AssertFormat = true

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
		if validationError, ok := err.(*jsonschema.ValidationError); ok {
			out := validationError.DetailedOutput()
			errorsResult := getErrors(out.Errors)
			return errorsResult, nil
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

func (resourceMinimumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMinimum, ok := m["resourceMinimum"]; ok {
		resourceMinimumStr := resourceMinimum.(string)
		return resourceMinimumSchema(resourceMinimumStr), nil
	}
	return nil, nil
}

func (resourceMaximumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMinimum, ok := m["resourceMaximum"]; ok {
		resourceMinimumStr := resourceMinimum.(string)
		return resourceMaximumSchema(resourceMinimumStr), nil
	}
	return nil, nil
}

func (s resourceMinimumSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "resourceMinimum"
	ruleResourceMinimumStr := string(s)
	dataValueParsedQty, err := resource.ParseQuantity(dataValue.(string))
	if err != nil {
		if err != nil {
			return ctx.Error(keywordPath, "failed parsing value %v", dataValue)
		}
	}

	rmSchemaValueParsedQ, err := resource.ParseQuantity(ruleResourceMinimumStr)
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing value %v", ruleResourceMinimumStr)
	}

	rmDecStr := dataValueParsedQty.AsDec().String()
	rmRfDecStr := rmSchemaValueParsedQ.AsDec().String()

	resourceMinimumSchemaVal, err := strconv.ParseFloat(rmDecStr, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", resourceMinimumSchemaVal)
	}

	resourceMinimumDataVal, err := strconv.ParseFloat(rmRfDecStr, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", resourceMinimumDataVal)
	}

	if resourceMinimumDataVal > resourceMinimumSchemaVal {
		return ctx.Error(keywordPath, "%v is lower then resourceMinimum %v", dataValue, ruleResourceMinimumStr)
	}
	return nil
}

func (s resourceMaximumSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "resourceMaximum"
	ruleResourceMaximumStr := string(s)
	dataValueParsedQty, err := resource.ParseQuantity(dataValue.(string))
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing value %v", dataValue)
	}

	rmSchemaValueParsedQ, err := resource.ParseQuantity(ruleResourceMaximumStr)
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing value %v", ruleResourceMaximumStr)
	}

	rmDecStr := dataValueParsedQty.AsDec().String()
	rmRfDecStr := rmSchemaValueParsedQ.AsDec().String()

	resourceMaximumSchemaVal, err := strconv.ParseFloat(rmDecStr, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", resourceMaximumSchemaVal)
	}

	resourceMaximumDataVal, err := strconv.ParseFloat(rmRfDecStr, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", resourceMaximumDataVal)
	}

	if resourceMaximumDataVal < resourceMaximumSchemaVal {
		return ctx.Error(keywordPath, "%v is greater then resourceMaximum %v", dataValue, ruleResourceMaximumStr)
	}
	return nil
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
