package extractor

import (
	"path/filepath"
	"testing"

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

func test_existed_file() *toAbsolutePathsTestCase {
	p, _ := filepath.Abs("../../internal/fixtures/kube/pass-all.yaml")
	return &toAbsolutePathsTestCase{
		name: "existed file, should return abs path with no errors",
		args: struct{ path string }{
			path: "../../internal/fixtures/kube/pass-all.yaml",
		},
		expected: struct {
			path string
		}{
			path: p,
		},
	}
}

func test_not_existed_file() *toAbsolutePathsTestCase {
	return &toAbsolutePathsTestCase{
		name: "test not existed file, should return an error",
		args: struct{ path string }{
			path: "../../internal/fixtures/kube/bla.yaml",
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
			path: "../../internal/fixtures/kube",
		},
		expected: struct {
			path string
		}{
			path: "",
		},
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
			absolutePath, _ := ToAbsolutePath(tt.args.path)
			assert.Equal(t, tt.expected.path, absolutePath)
		})
	}
}

func TestExtractConfigurationsFromYamlFile(t *testing.T) {
	t.Run("invalid yaml path, should return invalid file with returned ToAbsolutePath error", func(t *testing.T) {
		configurations, absolutePath, err := ExtractConfigurationsFromYamlFile("")

		assert.Empty(t, configurations)
		assert.Empty(t, absolutePath)
		assert.NotNil(t, err)
	})

	t.Run("invalid yaml file, should return an error", func(t *testing.T) {
		path := "../../internal/fixtures/policyAsCode/invalid-yaml.yaml"
		configurations, absolutePath, err := ExtractConfigurationsFromYamlFile(path)

		assert.Empty(t, configurations)
		assert.Empty(t, absolutePath)
		assert.Equal(t, "yaml validation error: yaml: line 2: did not find expected key\n", err.ValidationErrors[0].Error())
	})

	t.Run("valid yaml file, should return a configuration and no error", func(t *testing.T) {
		path := "../../internal/fixtures/jsonSchema/yamlSchema.yaml"
		configurations, absolutePath, err := ExtractConfigurationsFromYamlFile(path)

		assert.NotEmpty(t, configurations)
		assert.NotEmpty(t, absolutePath)
		assert.Nil(t, err)
	})
	t.Run("should not panic even though annotations is an array instead of object", func(t *testing.T) {
		path := "./extractorTestFiles/annotationsAreArray.yaml"
		configurations, _, _ := ExtractConfigurationsFromYamlFile(path)

		assert.Nil(t, (*configurations)[0].Annotations)
		assert.Equal(t, "rss-site", (*configurations)[0].MetadataName)
		assert.Equal(t, "Deployment", (*configurations)[0].Kind)
		assert.Equal(t, "apps/v1", (*configurations)[0].ApiVersion)
	})
	t.Run("should not panic even though metadata does not exist", func(t *testing.T) {
		path := "./extractorTestFiles/noMetadata.yaml"
		configurations, _, _ := ExtractConfigurationsFromYamlFile(path)

		assert.Nil(t, (*configurations)[0].Annotations)
		assert.Equal(t, "", (*configurations)[0].MetadataName)
		assert.Equal(t, "Deployment", (*configurations)[0].Kind)
		assert.Equal(t, "apps/v1", (*configurations)[0].ApiVersion)
	})
}
