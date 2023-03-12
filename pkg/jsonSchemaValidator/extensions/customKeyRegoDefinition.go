// This file defines a custom key to implement the logic for the rule:
// https://hub.datree.io/built-in-rules/ensure-memory-request-limit-equal

package jsonSchemaValidator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

var regoCodeToEval = RegoDefinition{}
var mainRegoPackage = ""

const regoDefinitionCustomKey = "regoDefinition"

type CustomKeyRegoDefinitionCompiler struct{}

type CustomKeyRegoDefinitionSchema map[string]interface{}

var CustomKeyRegoRule = jsonschema.MustCompileString("customKeyRegoDefinition.json", `{
	"properties" : {
		"regoCode": {
			"type": "string"
		}
	}
}`)

func (CustomKeyRegoDefinitionCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyRegoRule, ok := m[regoDefinitionCustomKey]; ok {
		customKeyRegoRuleStr, validStr := customKeyRegoRule.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("regoCode must be a string")
		}

		b, _ := json.Marshal(customKeyRegoRule)
		var regoDefinition RegoDefinition
		err := json.Unmarshal(b, &regoDefinition)
		if err != nil {
			// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
			return nil, fmt.Errorf("regoCode must be a string")
		}

		regoCodeToEval = regoDefinition
		mainRegoPackage = getPackageFromRegoCode(regoDefinition.Code)
		return CustomKeyRegoDefinitionSchema(customKeyRegoRuleStr), nil
	}
	return nil, nil
}

type RegoDefinition struct {
	Libs []string `json:"libs"`
	Code string   `json:"code"`
}

func (s CustomKeyRegoDefinitionSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	mainModuleFileName := "main.rego"
	regoFunctionEntryPoint := "violation"
	libsModuleGeneralName := "lib.rego"

	regoCtx := context.Background()
	var regoObjectParts []func(r *rego.Rego)
	regoObjectParts = append(regoObjectParts, rego.Query("data."+mainRegoPackage+"."+regoFunctionEntryPoint))

	regoObjectParts = append(regoObjectParts, rego.Module(mainModuleFileName, regoCodeToEval.Code))

	for _, lib := range regoCodeToEval.Libs {
		rnd := rand.Intn(100)
		regoObjectParts = append(regoObjectParts, rego.Module(libsModuleGeneralName+strconv.Itoa(rnd), lib))
	}
	r := rego.New(regoObjectParts...)

	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(regoCtx)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the prepared query.
	rs, err := query.Eval(regoCtx, rego.EvalInput(dataValue))

	if err != nil || len(rs) != 1 || len(rs[0].Expressions) != 1 {
		// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
		return nil
	} else {
		resultsValue := (rs[0].Expressions[0].Value).([]interface{})
		if value, ok := resultsValue[0].(bool); ok {
			if value {
				return ctx.Error(regoDefinitionCustomKey, "values in data value %v do not match", rs[0].Expressions[0].Value)
			}
			return nil
		}
	}

	return nil
}

func getPackageFromRegoCode(regoCode string) string {
	const PACKAGE = "package"
	// find the index of string "package"
	index := strings.Index(regoCode, PACKAGE)
	// get next single word after "package"
	packageStr := strings.Fields(regoCode[index:])
	return packageStr[1]
}
