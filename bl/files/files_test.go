package files

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/datreeio/datree/pkg/extractor"

	"github.com/stretchr/testify/assert"
)

type toAbsolutePathsTestCase struct {
	name string
	args struct {
		path string
	}
	expected struct {
		path string
	}
}

func TestToAbsolutePath(t *testing.T) {
	tests := []*toAbsolutePathsTestCase{
		test_existed_file(),
		test_directory_file(),
		test_not_existed_file(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			absolutePath, _ := extractor.ToAbsolutePath(tt.args.path)
			assert.Equal(t, tt.expected.path, absolutePath)
		})
	}
}

func test_existed_file() *toAbsolutePathsTestCase {
	const filePathName = "../../fixtures/kube/pass-all.yaml"
	p, _ := filepath.Abs(filePathName)
	return &toAbsolutePathsTestCase{
		name: "existed file, should return abs path with no errors",
		args: struct{ path string }{
			path: filePathName ,
		},
		expected: struct {
			path string
		}{
			path: p,
		},
	}
}

func test_not_existed_file() *toAbsolutePathsTestCase {
	const fileNamePath = "../../fixtures/kube/bla.yaml"
	return &toAbsolutePathsTestCase{
		name: "test not existed file, should return an error",
		args: struct{ path string }{
			path: fileNamePath ,
		},
		expected: struct {
			path string
		}{
			path: "",
		},
	}
}

func test_directory_file() *toAbsolutePathsTestCase {
	const fileNamePath = "../../fixtures/kube"
	return &toAbsolutePathsTestCase{
		args: struct{ path string }{
			path: fileNamePath,
		},
		expected: struct {
			path string
		}{
			path: "",
		},
	}
}

func TestExtractYamlFileToUnknownStruct(t *testing.T) {
	const fileNamePath = "../../fixtures/policyAsCode/valid-schema.yaml"
	t.Run("valid yaml file, should return an unknown struct and no error", func(t *testing.T) {
		actualResult, actualErr := ExtractYamlFileToUnknownStruct(fileNamePath)
		assert.NotEqual(t, nil, actualResult)
		assert.Equal(t, nil, actualErr)
	})

	t.Run("invalid yaml file, should return an error", func(t *testing.T) {
		const fileNamePath = "../../fixtures/policyAsCode/invalid-yaml.yaml"
		actualResult, actualErr := ExtractYamlFileToUnknownStruct(fileNamePath)
		assert.Equal(t, UnknownStruct(nil), actualResult)
		assert.NotEqual(t, nil, actualErr)
		assert.Equal(t, errors.New("yaml: line 2: did not find expected key"), actualErr)
	})
}
