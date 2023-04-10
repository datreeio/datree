// This file defines a custom key to implement the logic for rego rule:

package jsonSchemaValidator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"strings"
)

const RegoDefinitionCustomKey = "regoDefinition"

type CustomKeyRegoDefinitionCompiler struct{}

type CustomKeyRegoDefinitionSchema map[string]interface{}

var CustomKeyRegoRule = jsonschema.MustCompileString("customKeyRegoDefinition.json", `{
	"properties" : {
		"regoDefinition": {
			"type": "object"
		}
	}
}`)

func (CustomKeyRegoDefinitionCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyRegoRule, ok := m[RegoDefinitionCustomKey]; ok {
		customKeyRegoRuleObj, validObject := customKeyRegoRule.(map[string]interface{})
		if !validObject {
			return nil, fmt.Errorf("regoDefinition must be an object")
		}

		regoDefinitionSchema, err := convertCustomKeyRegoDefinitionSchemaToRegoDefinitionSchema(customKeyRegoRuleObj)
		if err != nil {
			return nil, err
		}

		if regoDefinitionSchema.Code == "" {
			return nil, fmt.Errorf("regoDefinition.code can't be empty")
		}

		return CustomKeyRegoDefinitionSchema(customKeyRegoRuleObj), nil
	}
	return nil, nil
}

type RegoDefinition struct {
	Libs []string `json:"libs"`
	Code string   `json:"code"`
}

func (customKeyRegoDefinitionSchema CustomKeyRegoDefinitionSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	regoDefinitionSchema, err := convertCustomKeyRegoDefinitionSchemaToRegoDefinitionSchema(customKeyRegoDefinitionSchema)
	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, err.Error())
	}

	regoCtx := context.Background()

	regoObject, err := retrieveRegoFromSchema(regoDefinitionSchema)
	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, "can't compile rego code, %s", err.Error())
	}

	// Create a prepared query that can be evaluated.
	query, err := regoObject.PrepareForEval(regoCtx)
	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, "can't compile rego code, %s", err.Error())
	}

	// Execute the prepared query.
	rs, err := query.Eval(regoCtx, rego.EvalInput(dataValue))

	if err != nil {
		return ctx.Error(CustomKeyValidationErrorKeyPath, "failed to evaluate rego due to %s", err.Error())
	}

	if len(rs) != 1 || len(rs[0].Expressions) != 1 {
		return ctx.Error(CustomKeyValidationErrorKeyPath, "failed to evaluate rego, unexpected results")
	}

	resultValues := (rs[0].Expressions[0].Value).([]interface{})
	for _, resultValue := range resultValues {
		violationReturnValue, ok := resultValue.(bool)
		if !ok {
			return ctx.Error(CustomKeyValidationErrorKeyPath, "violation needs to return a boolean")
		}
		if violationReturnValue {
			return ctx.Error(RegoDefinitionCustomKey, "values in data value %v do not match", dataValue)
		}
	}

	return nil
}

func getPackageFromRegoCode(regoCode string) (string, error) {
	const PACKAGE = "package"
	// find the index of string "package"
	index := strings.Index(regoCode, PACKAGE)
	if index == -1 {
		return "", fmt.Errorf("rego code must have a package")
	}
	// get next single word after "package"
	packageStr := strings.Fields(regoCode[index:])
	return packageStr[1], nil
}

func retrieveRegoFromSchema(regoDefinitionSchema *RegoDefinition) (*rego.Rego, error) {
	const mainModuleFileName = "main.rego"
	const regoFunctionEntryPoint = "violation"

	mainRegoPackage, err := getPackageFromRegoCode(regoDefinitionSchema.Code)
	if err != nil {
		return nil, err
	}

	var regoObjectParts []func(r *rego.Rego)
	regoObjectParts = append(regoObjectParts, rego.Query("data."+mainRegoPackage+"."+regoFunctionEntryPoint))

	regoObjectParts = append(regoObjectParts, rego.Module(mainModuleFileName, regoDefinitionSchema.Code))

	for _, lib := range regoDefinitionSchema.Libs {
		libPackageName, err := getPackageFromRegoCode(lib)
		if err != nil {
			return nil, err
		}
		regoObjectParts = append(regoObjectParts, rego.Module(libPackageName, lib))
	}
	regoObject := rego.New(regoObjectParts...)
	return regoObject, nil
}

func convertCustomKeyRegoDefinitionSchemaToRegoDefinitionSchema(regoDefinitionSchema CustomKeyRegoDefinitionSchema) (*RegoDefinition, error) {
	b, err := json.Marshal(regoDefinitionSchema)
	if err != nil {
		return nil, fmt.Errorf("regoDefinition failed to marshal to json, %s", err.Error())
	}

	var regoDefinition RegoDefinition
	err = json.Unmarshal(b, &regoDefinition)
	if err != nil {
		return nil, fmt.Errorf("regoDefinition must be an object of type RegoDefinition %s", err.Error())
	}
	return &regoDefinition, nil
}
