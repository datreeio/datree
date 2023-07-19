// This file defines a custom key to implement the logic for cel rule:

package jsonSchemaValidator

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

const CELDefinitionCustomKey = "CELDefinition"

type CustomKeyCELDefinitionCompiler struct{}

type CustomKeyCELDefinitionSchema []interface{}

var CustomKeyCELRule = jsonschema.MustCompileString("customKeyCELDefinition.json", `{
	"properties" : {
		"CELDefinition": {
			"type": "array"
		}
	}
}`)

func (CustomKeyCELDefinitionCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyCELRule, ok := m[CELDefinitionCustomKey]; ok {
		customKeyCELRuleObj, validObject := customKeyCELRule.([]interface{})
		if !validObject {
			return nil, fmt.Errorf("CELDefinition must be an array")
		}

		CELDefinitionSchema, err := convertCustomKeyCELDefinitionSchemaToCELDefinitionSchema(customKeyCELRuleObj)
		if err != nil {
			return nil, err
		}

		if len(CELDefinitionSchema.CELExpressions) == 0 {
			return nil, fmt.Errorf("CELDefinition can't be empty")
		}

		return CustomKeyCELDefinitionSchema(customKeyCELRuleObj), nil
	}
	return nil, nil
}

func (customKeyCELDefinitionSchema CustomKeyCELDefinitionSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	CELDefinitionSchema, err := convertCustomKeyCELDefinitionSchemaToCELDefinitionSchema(customKeyCELDefinitionSchema)
	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, err.Error())
	}
	// wrap dataValue (the resource that should be validated) inside a struct with parent object key
	resourceWithParentKey := make(map[string]interface{})
	resourceWithParentKey["object"] = dataValue

	// prepare CEL env inputs - in our case the only input is the resource that should be validated
	inputs, err := getCELEnvInputs(resourceWithParentKey)
	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, err.Error())
	}

	env, err := cel.NewEnv(inputs...)
	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, err.Error())
	}

	for _, celExpression := range CELDefinitionSchema.CELExpressions {
		ast, issues := env.Compile(celExpression.Expression)
		if issues != nil && issues.Err() != nil {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "cel expression compile error: %s", issues.Err())
		}

		prg, err := env.Program(ast)
		if err != nil {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "cel program construction error: %s", err)
		}

		res1, _, err := prg.Eval(resourceWithParentKey)
		if err != nil {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "cel evaluation error: %s", err)
		}

		if res1.Type().TypeName() != "bool" {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "cel expression needs to return a boolean")
		}

		celReturnValue, ok := res1.Value().(bool)
		if !ok {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "cel expression needs to return a boolean")
		}
		if !celReturnValue {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "cel expression failure message: %s", celExpression.Message)
		}
	}

	return nil
}

type CELExpression struct {
	Expression string `json:"expression"`
	Message    string `json:"message"`
}

type CELDefinition struct {
	CELExpressions []CELExpression
}

func convertCustomKeyCELDefinitionSchemaToCELDefinitionSchema(CELDefinitionSchema CustomKeyCELDefinitionSchema) (*CELDefinition, error) {
	var CELDefinition CELDefinition
	for _, CELExpressionFromSchema := range CELDefinitionSchema {
		var CELExpression CELExpression
		b, err := json.Marshal(CELExpressionFromSchema)
		if err != nil {
			return nil, fmt.Errorf("CELExpression failed to marshal to json, %s", err.Error())
		}
		err = json.Unmarshal(b, &CELExpression)
		if err != nil {
			return nil, fmt.Errorf("CELExpression must be an object of type CELExpression %s", err.Error())
		}
		CELDefinition.CELExpressions = append(CELDefinition.CELExpressions, CELExpression)
	}

	return &CELDefinition, nil
}

func getCELEnvInputs(dataValue map[string]interface{}) ([]cel.EnvOption, error) {
	inputVars := make([]cel.EnvOption, 0, len(dataValue))
	for input := range dataValue {
		inputVars = append(inputVars, cel.Variable(input, cel.DynType))
	}
	return inputVars, nil
}
