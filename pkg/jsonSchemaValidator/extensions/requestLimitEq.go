// This file defines 2 custom keys in order to implement the logic for these rules:
// Link to docs for both rules

package jsonSchemaValidator

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type CkCpuEqCompiler struct{}
type CkMemoryEqCompiler struct{}

type CkCpuEqSchema map[string]interface{}
type CkMemoryEqSchema map[string]interface{}

var CkCpuEq = jsonschema.MustCompileString("ckCpuEq.json", `{
	"properties" : {
		"ckCpuEq": {
			"type": "string"
		}
	}
}`)

var CkMemoryEq = jsonschema.MustCompileString("ckMemoryEq.json", `{
	"properties" : {
		"ckMemoryEq": {
			"type": "string"
		}
	}
}`)

func (CkCpuEqCompiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	// convert ctx to json and print it

	if ckCpuEq, ok := m["ckCpuEq"]; ok {
		ckCpuEqStr, validStr := ckCpuEq.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("ckCpuEq must be a string")
		}
		return CkCpuEqSchema(ckCpuEqStr), nil
	}
	return nil, nil
}

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

func (s CkCpuEqSchema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	keywordPath := "ckCpuEq"

	b, _ := json.Marshal(dataValue)
	var resources Resources
	err := json.Unmarshal(b, &resources)
	if err != nil {
		// We expect a certain format, if this fails that means the format is wrong and we can ignore this validation and move on
		return nil
	}

	if (resources.Requests.CPU == "" || resources.Limits.CPU == "") && (resources.Requests.CPU != resources.Limits.CPU) {
		return ctx.Error(keywordPath, "one or more empty values in data value %v", dataValue)
	}

	if resources.Requests.CPU != resources.Limits.CPU {
		return ctx.Error(keywordPath, "values in data value %v do not match", dataValue)
	}

	return nil
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
