package jsonSchemaValidator

import (
	"encoding/json"
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

func (jsv *JSONSchemaValidator) ValidateYamlSchema(schemaContent string, yamlContent string) ([]jsonschema.Detailed, error) {
	jsonSchema, _ := yaml.YAMLToJSON([]byte(schemaContent))
	jsonYamlContent, _ := yaml.YAMLToJSON([]byte(yamlContent))
	return jsv.Validate(string(jsonSchema), jsonYamlContent)
}

func (jsv *JSONSchemaValidator) Validate(schemaContent string, yamlContent []byte) ([]jsonschema.Detailed, error) {
	var jsonYamlContent interface{}
	if err := json.Unmarshal(yamlContent, &jsonYamlContent); err != nil {
		return nil, err
	}

	compiler := jsonschema.NewCompiler()
	//format is treated as annotation in draft-2019 onwards. it needs to be explicitly enabled by compiler.AssertFormat = true.
	//see reference: https://github.com/santhosh-tekuri/jsonschema/issues/43
	compiler.AssertFormat = true

	if err := compiler.AddResource("schema.json", strings.NewReader(schemaContent)); err != nil {
		return nil, err
	}

	compiler.RegisterExtension("resourceMinimum", resourceMinimum, resourceMinimumCompiler{})
	compiler.RegisterExtension("resourceMaximum", resourceMaximum, resourceMaximumCompiler{})

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, err
	}

	err = schema.Validate(jsonYamlContent)

	if err != nil {
		if validationError, ok := err.(*jsonschema.ValidationError); ok {
			return getOnlyRelevantErrors(validationError.DetailedOutput()), nil
		} else {
			fmt.Fprintf(os.Stderr, "validation failed: %v\n", err)
			return nil, err
		}
	}
	return nil, nil
}

func (resourceMinimumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMinimum, ok := m["resourceMinimum"]; ok {
		resourceMinimumStr, validStr := resourceMinimum.(string)
		if !validStr {
			return nil, fmt.Errorf("resourceMinimum must be a string")
		}
		return resourceMinimumSchema(resourceMinimumStr), nil
	}
	return nil, nil
}

func (resourceMaximumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMaximum, ok := m["resourceMaximum"]; ok {
		resourceMaximumStr, validStr := resourceMaximum.(string)
		if !validStr {
			return nil, fmt.Errorf("resourceMaximum must be a string")
		}
		return resourceMaximumSchema(resourceMaximumStr), nil
	}
	return nil, nil
}

func (s resourceMinimumSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "resourceMinimum"
	schemaResourceMinStr := string(s)

	dataValueStr, validStr := dataValue.(string)
	if !validStr {
		return ctx.Error(keywordPath, "%s must be a string", dataValueStr)
	}

	dataValueParsedQty, err := resource.ParseQuantity(dataValueStr)
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing data value %v", dataValue)
	}

	schemaResourceMinParsedQty, err := resource.ParseQuantity(schemaResourceMinStr)
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing schema value %v", schemaResourceMinStr)
	}

	dataValueParsedQtyDecimal := dataValueParsedQty.AsDec().String()
	schemaResourceMinParsedQtyDecimal := schemaResourceMinParsedQty.AsDec().String()

	dataValueFloat, err := strconv.ParseFloat(dataValueParsedQtyDecimal, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", dataValueFloat)
	}

	schemaMinValueFloat, err := strconv.ParseFloat(schemaResourceMinParsedQtyDecimal, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", schemaMinValueFloat)
	}

	if schemaMinValueFloat > dataValueFloat {
		return ctx.Error(keywordPath, "%v is lower then resourceMinimum %v", dataValue, schemaResourceMinStr)
	}
	return nil
}

func (s resourceMaximumSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "resourceMaximum"
	schemaResourceMaxStr := string(s)

	dataValueStr, validStr := dataValue.(string)
	if !validStr {
		return ctx.Error(keywordPath, "%s must be a string", dataValueStr)
	}

	dataValueParsedQty, err := resource.ParseQuantity(dataValue.(string))
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing data value %v", dataValue)
	}

	schemaResourceMaxParsedQty, err := resource.ParseQuantity(schemaResourceMaxStr)
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing schema value %v", schemaResourceMaxStr)
	}

	dataValueParsedQtyDecimal := dataValueParsedQty.AsDec().String()
	schemaResourceMaxParsedQtyDecimal := schemaResourceMaxParsedQty.AsDec().String()

	dataValueFloat, err := strconv.ParseFloat(dataValueParsedQtyDecimal, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", dataValueFloat)
	}

	schemaMaxValueFloat, err := strconv.ParseFloat(schemaResourceMaxParsedQtyDecimal, 64)
	if err != nil {
		return ctx.Error(keywordPath, "failed float parsing value %v", schemaMaxValueFloat)
	}

	if schemaMaxValueFloat < dataValueFloat {
		return ctx.Error(keywordPath, "%v is greater then resourceMaximum %v", dataValue, schemaResourceMaxStr)
	}
	return nil
}

func getOnlyRelevantErrors(rootError jsonschema.Detailed) []jsonschema.Detailed {
	return getLeafErrors(rootError)
}

func getLeafErrors(error jsonschema.Detailed) []jsonschema.Detailed {
	if error.Errors == nil {
		// if no more child errors, I am a leaf, return me!
		return []jsonschema.Detailed{error}
	} else if strings.HasSuffix(error.KeywordLocation, "anyOf") {
		// if I am an anyOf node, return as if I am a leaf
		return []jsonschema.Detailed{error}
	} else {
		// if I'm not a leaf, return the errors from all my children!
		var errorsFromChildren []jsonschema.Detailed
		for _, childError := range error.Errors {
			errorsFromChildren = append(errorsFromChildren, getLeafErrors(childError)...)
		}
		return errorsFromChildren
	}
}
