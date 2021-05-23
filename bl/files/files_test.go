package files

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type toAbsolutePathsTestCase struct {
	name string
	args struct {
		paths []string
	}
	expected struct {
		path string
	}
}

func TestToAbsolutePaths(t *testing.T) {
	tests := []*toAbsolutePathsTestCase{
		test_existed_file(),
		test_directory_file(),
		test_not_existed_file(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathsChan := ToAbsolutePaths(tt.args.paths)
			assert.Equal(t, tt.expected.path, <-pathsChan)
		})
	}

}

func test_existed_file() *toAbsolutePathsTestCase {
	p, _ := filepath.Abs("../../internal/fixtures/kube/pass-all.yaml")
	return &toAbsolutePathsTestCase{
		name: "existed file, should return abs path with no errors",
		args: struct{ paths []string }{
			paths: []string{"../../internal/fixtures/kube/pass-all.yaml"},
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
		args: struct{ paths []string }{
			paths: []string{"../../internal/fixtures/kube/bla.yaml"},
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
		args: struct{ paths []string }{
			paths: []string{"../../internal/fixtures/kube"},
		},
		expected: struct {
			path string
		}{
			path: "",
		},
	}
}
