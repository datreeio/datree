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

func (jsv *JSONSchemaValidator) ValidateYamlSchema(schemaContent string, yamlContent string) ([]jsonschema.Detailed, error) {
	jsonSchema, _ := yaml.YAMLToJSON([]byte(schemaContent))
	jsonYamlContent, _ := yaml.YAMLToJSON([]byte(yamlContent))
	return jsv.Validate(string(jsonSchema), jsonYamlContent)
}

func (jsv *JSONSchemaValidator) Validate(schemaContent string, yamlContent []byte) ([]jsonschema.Detailed, error) {
	var jsonYamlContent interface{}
	//todo what happens if unmarshal fails?
	if err := json.Unmarshal(yamlContent, &jsonYamlContent); err != nil {
		panic(err)
	}

	compiler := jsonschema.NewCompiler()
	//format is treated as annotation in draft-2019 onwards. it needs to be explicitly enabled by compiler.AssertFormat = true.
	//see reference: https://github.com/santhosh-tekuri/jsonschema/issues/43
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
		resourceMinimumStr := resourceMinimum.(string)
		return resourceMinimumSchema(resourceMinimumStr), nil
	}
	return nil, nil
}

func (resourceMaximumCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if resourceMaximum, ok := m["resourceMaximum"]; ok {
		resourceMaximumStr := resourceMaximum.(string)
		return resourceMaximumSchema(resourceMaximumStr), nil
	}
	return nil, nil
}

//todo check type convertions and add error handling
func (s resourceMinimumSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "resourceMinimum"
	ruleResourceMinimumStr := string(s)
	dataValueParsedQty, err := resource.ParseQuantity(dataValue.(string))
	if err != nil {
		//todo add on which type of shiteness it happened
		return ctx.Error(keywordPath, "failed parsing value %v", dataValue)
	}
	//todo rename
	rmSchemaValueParsedQ, err := resource.ParseQuantity(ruleResourceMinimumStr)
	if err != nil {
		return ctx.Error(keywordPath, "failed parsing value %v", ruleResourceMinimumStr)
	}

	//todo rename?
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

//todo think it we need 2 functions / just one
func getOnlyRelevantErrors(rootError jsonschema.Detailed) []jsonschema.Detailed {
	return getLeafErrors(rootError)
}

func getLeafErrors(error jsonschema.Detailed) []jsonschema.Detailed {
	if error.Errors == nil {
		// if no more child errors, I am a leaf, return me!
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
