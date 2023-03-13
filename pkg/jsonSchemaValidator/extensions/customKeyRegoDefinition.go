// This file defines a custom key to implement the logic for the rule:
// https://hub.datree.io/built-in-rules/ensure-memory-request-limit-equal

package jsonSchemaValidator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

var regoCodeToEval = RegoDefinition{}
var mainRegoPackage = ""
var mainModuleFileName = "main.rego"
var regoFunctionEntryPoint = "violation"
var libsModuleGeneralName = "lib.rego"

const RegoDefinitionCustomKey = "regoDefinition"

type CustomKeyRegoDefinitionCompiler struct{}

type CustomKeyRegoDefinitionSchema map[string]interface{}

var CustomKeyRegoRule = jsonschema.MustCompileString("customKeyRegoDefinition.json", `{
	"properties" : {
		"regoDefinition": {
			"type": "string"
		}
	}
}`)

func (CustomKeyRegoDefinitionCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyRegoRule, ok := m[RegoDefinitionCustomKey]; ok {
		customKeyRegoRuleStr, validStr := customKeyRegoRule.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("regoDefinition must be a string")
		}

		b, _ := json.Marshal(customKeyRegoRule)
		var regoDefinition RegoDefinition
		err := json.Unmarshal(b, &regoDefinition)
		if err != nil {
			return nil, fmt.Errorf("regoDefinition must be a string")
		}

		regoCodeToEval = regoDefinition
		return CustomKeyRegoDefinitionSchema(customKeyRegoRuleStr), nil
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
		return ctx.Error(RegoDefinitionCustomKey, "can't compile rego code")
	}

	// Execute the prepared query.
	rs, err := query.Eval(regoCtx, rego.EvalInput(dataValue))

	if err != nil || len(rs) != 1 || len(rs[0].Expressions) != 1 {
		// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
		return nil
	} else {
		resultsValue := (rs[0].Expressions[0].Value).([]interface{})
		value, ok := resultsValue[0].(bool)
		if ok {
			if value {
				return ctx.Error(RegoDefinitionCustomKey, "values in data value %v do not match", rs[0].Expressions[0].Value)
			}
			return nil
		} else {
			return ctx.Error(RegoDefinitionCustomKey, "violation needs to return boolean")
		}
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
	mainRegoPackage = getPackageFromRegoCode(regoCode)

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
