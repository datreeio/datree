// This file defines a custom key to implement the logic for rego rule:

package jsonSchemaValidator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

const RegoDefinitionCustomKey = "regoDefinition"

var regoCodeToEval = RegoDefinition{}

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

		b, err := json.Marshal(customKeyRegoRule)
		if err != nil {
			return nil, fmt.Errorf("regoDefinition faild to marshal to json, %s", err.Error())
		}

		var regoDefinition RegoDefinition
		err = json.Unmarshal(b, &regoDefinition)
		if err != nil {
			return nil, fmt.Errorf("regoDefinition must be an object of type RegoDefinition")
		}

		regoCodeToEval = regoDefinition
		return CustomKeyRegoDefinitionSchema(customKeyRegoRuleObj), nil
	}
	return nil, nil
}

type RegoDefinition struct {
	Libs []string `json:"libs"`
	Code string   `json:"code"`
}

func (s CustomKeyRegoDefinitionSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	regoCtx := context.Background()

	r := retrieveRegoFromSchema(regoCodeToEval.Code)

	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(regoCtx)
	if err != nil {
		return ctx.Error(RegoDefinitionCustomKey, "can't compile rego code, %s", err.Error())
	}

	// Execute the prepared query.
	rs, err := query.Eval(regoCtx, rego.EvalInput(dataValue))

	if err != nil {
		return ctx.Error(RegoDefinitionCustomKey, "failed to evaluate rego due to %s", err.Error())
	}

	if len(rs) != 1 || len(rs[0].Expressions) != 1 {
		return ctx.Error(RegoDefinitionCustomKey, "failed to evaluate rego, unexpected results")
	}

	resultsValue := (rs[0].Expressions[0].Value).([]interface{})
	violationReturnValue, ok := resultsValue[0].(bool)
	if ok {
		if violationReturnValue {
			return ctx.Error(RegoDefinitionCustomKey, "values in data value %v do not match", rs[0].Expressions[0].Value)
		}
		return nil
	} else {
		return ctx.Error(RegoDefinitionCustomKey, "violation needs to return a boolean")
	}
}

func getPackageFromRegoCode(regoCode string) string {
	const PACKAGE = "package"
	// find the index of string "package"
	index := strings.Index(regoCode, PACKAGE)
	// get next single word after "package"
	packageStr := strings.Fields(regoCode[index:])
	return packageStr[1]
}

func retrieveRegoFromSchema(regoCode string) *rego.Rego {
	var mainModuleFileName = "main.rego"
	var regoFunctionEntryPoint = "violation"

	mainRegoPackage := getPackageFromRegoCode(regoCode)

	var regoObjectParts []func(r *rego.Rego)
	regoObjectParts = append(regoObjectParts, rego.Query("data."+mainRegoPackage+"."+regoFunctionEntryPoint))

	regoObjectParts = append(regoObjectParts, rego.Module(mainModuleFileName, regoCode))

	for _, lib := range regoCodeToEval.Libs {
		libPackageName := getPackageFromRegoCode(lib)
		regoObjectParts = append(regoObjectParts, rego.Module(libPackageName, lib))
	}
	r := rego.New(regoObjectParts...)
	return r
}
