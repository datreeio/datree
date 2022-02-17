package files

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/datreeio/datree/pkg/extractor"

	"github.com/stretchr/testify/assert"
)
const FIXTURES_PATH string = "../../fixtures/kube";
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
	filesPath := fmt.Sprintf("%s/pass-all.yaml", FIXTURES_PATH)

	p, _ := filepath.Abs(filesPath)
	return &toAbsolutePathsTestCase{
		name: "existed file, should return abs path with no errors",
		args: struct{ path string }{
			path: filesPath,
		},
		expected: struct {
			path string
		}{
			path: p,
		},
	}
} 

func test_not_existed_file() *toAbsolutePathsTestCase {
	filesPath := fmt.Sprintf("%s/bla.yaml", FIXTURES_PATH)

	return &toAbsolutePathsTestCase{
		name: "test not existed file, should return an error",
		args: struct{ path string }{
			path: filesPath,
		},
		expected: struct {
			path string
		}{
			path: "",
		},
	}
}

func test_directory_file() *toAbsolutePathsTestCase {
	
	return &toAbsolutePathsTestCase{
		args: struct{ path string }{
			path: FIXTURES_PATH,
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
