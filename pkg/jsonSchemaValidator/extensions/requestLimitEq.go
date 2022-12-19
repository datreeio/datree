// This file defines a custom key to implement the logic for the rule:
// https://hub.datree.io/built-in-rules/ensure-memory-request-limit-equal

package jsonSchemaValidator

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type CkMemoryEqCompiler struct{}

type CkMemoryEqSchema map[string]interface{}

var CkMemoryEq = jsonschema.MustCompileString("ckMemoryEq.json", `{
	"properties" : {
		"ckMemoryEq": {
			"type": "string"
		}
	}
}`)

func (CkMemoryEqCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	// convert ctx to json and print it

	if ckMemoryEq, ok := m["ckMemoryEq"]; ok {
		ckMemoryEqStr, validStr := ckMemoryEq.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("ckMemoryEq must be a string")
		}
		return CkMemoryEqSchema(ckMemoryEqStr), nil
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

func (s CkMemoryEqSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "ckMemoryEq"

	b, _ := json.Marshal(dataValue)
	var resources Resources
	err := json.Unmarshal(b, &resources)
	if err != nil {
		// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
		return nil
	}

	if (resources.Requests.Memory == "" || resources.Limits.Memory == "") && (resources.Requests.Memory != resources.Limits.Memory) {
		return ctx.Error(keywordPath, "one or more empty values in data value %v", dataValue)
	}

	if resources.Requests.Memory != resources.Limits.Memory {
		return ctx.Error(keywordPath, "values in data value %v do not match", dataValue)
	}

	return nil
}
