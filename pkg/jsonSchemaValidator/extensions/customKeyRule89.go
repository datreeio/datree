// This file defines a custom key to implement the logic for the rule:
// https://hub.datree.io/built-in-rules/ensure-hostpath-mounts-readonly

package jsonSchemaValidator

import (
	"fmt"

	"github.com/itchyny/gojq"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type CustomKeyRule89Compiler struct{}

type CustomKeyRule89Schema map[string]interface{}

var CustomKeyRule89 = jsonschema.MustCompileString("customKeyRule89.json", `{
	"properties" : {
		"customKeyRule89": {
			"type": "string"
		}
	}
}`)

func (CustomKeyRule89Compiler) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (jsonschema.ExtSchema, error) {
	if customKeyRule89, ok := m["customKeyRule89"]; ok {
		customKeyRule89Str, validStr := customKeyRule89.(map[string]interface{})
		if !validStr {
			return nil, fmt.Errorf("customKeyRule89 must be a string")
		}
		return CustomKeyRule89Schema(customKeyRule89Str), nil
	}
	return nil, nil
}

func (s CustomKeyRule89Schema) Validate(ctx jsonschema.ValidationContext, dataValue interface{}) error {
	namesOfVolumesWithHostPath := getNamesOfVolumesWithHostPath(dataValue)
	namesOfVolumeMountsWithoutReadonly := getNamesOfVolumeMountsWithoutReadonly(dataValue)

	for _, nameOfVolumeMountWithoutReadonly := range namesOfVolumeMountsWithoutReadonly {
		for _, nameOfVolumeWithHostPath := range namesOfVolumesWithHostPath {
			if nameOfVolumeMountWithoutReadonly == nameOfVolumeWithHostPath {
				return ctx.Error(CustomKeyValidationErrorKeyPath, "volumeMounts: a container is using a hostPath volume without setting it to read-only")
			}
		}
	}
	return nil
}

func getNamesOfVolumesWithHostPath(dataValue interface{}) []string {
	var hostPathVolumes []string

	// Get all volumes with hostPath
	hostPathVolumeNameQuery, _ := gojq.Parse(".volumes[] | select(.hostPath != null) | .name")
	hostPathVolumeNameIter := hostPathVolumeNameQuery.Run(dataValue)
	for {
		volumeName, ok := hostPathVolumeNameIter.Next()
		if !ok {
			break
		}
		if _, ok := volumeName.(error); ok {
			break
		}

		hostPathVolumes = append(hostPathVolumes, volumeName.(string))
	}
	return hostPathVolumes
}

func getNamesOfVolumeMountsWithoutReadonly(dataValue interface{}) []string {
	var namesOfVolumeMountsWithoutReadonly []string

	// Get all volumeMounts without readOnly
	volumeMountsNameQuery, _ := gojq.Parse(".containers[] | .volumeMounts[] | select((.readOnly == null) or (.readOnly == false)) | .name")
	volumeMountsNameIter := volumeMountsNameQuery.Run(dataValue)
	for {
		volumeMountName, ok := volumeMountsNameIter.Next()
		if !ok || volumeMountName == nil {
			break
		}
		if _, ok := volumeMountName.(error); ok {
			break
		}

		namesOfVolumeMountsWithoutReadonly = append(namesOfVolumeMountsWithoutReadonly, volumeMountName.(string))
	}

	return namesOfVolumeMountsWithoutReadonly
}
