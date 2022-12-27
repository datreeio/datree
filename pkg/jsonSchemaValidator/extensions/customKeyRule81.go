// This file defines a custom key to implement the logic for the rule:
// https://hub.datree.io/built-in-rules/ensure-memory-request-limit-equal

package jsonSchemaValidator

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type CustomKeyRule81Compiler struct{}

type CustomKeyRule81Schema map[string]interface{}

var CustomKeyRule81 = jsonschema.MustCompileString("customKeyRule81.json", `{
	"properties" : {
		"customKeyRule81": {
			"type": "string"
		}
	}
}`)

func (CustomKeyRule81Compiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyRule81, ok := m["customKeyRule81"]; ok {
		customKeyRule81Str, validStr := customKeyRule81.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("customKeyRule81 must be a string")
		}
		return CustomKeyRule81Schema(customKeyRule81Str), nil
	}
	return nil, nil
}

type Resources struct {
	Requests Requests `json:"requests"`
	Limits   Limits   `json:"limits"`
}
type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Limits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

func (s CustomKeyRule81Schema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "customKeyRule81"

	b, _ := json.Marshal(dataValue)
	var resources Resources
	err := json.Unmarshal(b, &resources)
	if err != nil {
		// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
		return nil
	}

	if (resources.Requests.Memory == "" || resources.Limits.Memory == "") && (resources.Requests.Memory != resources.Limits.Memory) {
		return ctx.Error(keywordPath, "empty value in %v", dataValue)
	}

	if resources.Requests.Memory != resources.Limits.Memory {
		return ctx.Error(keywordPath, "values in data value %v do not match", dataValue)
	}

	return nil
}
