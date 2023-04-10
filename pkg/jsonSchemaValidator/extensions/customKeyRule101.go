// This file defines a custom key to implement the logic for the rule:
// https://hub.datree.io/built-in-rules/ensure-memory-request-limit-equal

package jsonSchemaValidator

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type CustomKeyRule101Compiler struct{}

type CustomKeyRule101Schema map[string]interface{}

var CustomKeyRule101 = jsonschema.MustCompileString("customKeyRule101.json", `{
	"properties" : {
		"customKeyRule101": {
			"type": "string"
		}
	}
}`)

func (CustomKeyRule101Compiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyRule101, ok := m["customKeyRule101"]; ok {
		customKeyRule101Str, validStr := customKeyRule101.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("customKeyRule101 must be a string")
		}
		return CustomKeyRule101Schema(customKeyRule101Str), nil
	}
	return nil, nil
}

type Rules101 struct {
	ApiGroups []string `json:"apiGroups"`
	Resources []string `json:"resources"`
	Verbs     []string `json:"verbs"`
}

func (s CustomKeyRule101Schema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	b, _ := json.Marshal(dataValue)
	var ruleArray []Rules101
	err := json.Unmarshal(b, &ruleArray)
	if err != nil {
		// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
		return nil
	}

	for _, rule := range ruleArray {
		for _, resource := range rule.Resources {
			if resource == "pods" || resource == "*" {
				for _, verb := range rule.Verbs {
					if verb == "create" || verb == "*" {
						return ctx.Error(CustomKeyValidationErrorKeyPath, "invalid verb or resource")
					}
				}
			}
		}
	}

	return nil
}
