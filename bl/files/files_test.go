package files

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractYamlFileToUnknownStruct(t *testing.T) {
	t.Run("valid yaml file, should return an unknown struct and no error", func(t *testing.T) {
		actualResult, actualErr := ExtractYamlFileToUnknownStruct("../../internal/fixtures/policyAsCode/valid-schema.yaml")
		assert.NotEqual(t, nil, actualResult)
		assert.Equal(t, nil, actualErr)
	})

	t.Run("invalid yaml file, should return an error", func(t *testing.T) {
		actualResult, actualErr := ExtractYamlFileToUnknownStruct("../../internal/fixtures/policyAsCode/invalid-yaml.yaml")
		assert.Equal(t, UnknownStruct(nil), actualResult)
		assert.NotEqual(t, nil, actualErr)
		assert.Equal(t, errors.New("yaml: line 2: did not find expected key"), actualErr)
	})
}
