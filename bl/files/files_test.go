package files

import (
	"fmt"
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
		err  error
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
			pathsChan, errors := ToAbsolutePaths(tt.args.paths)
			assert.Equal(t, tt.expected.path, <-pathsChan)
			assert.Equal(t, tt.expected.err, <-errors)
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
			err  error
		}{
			path: p,
			err:  nil,
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
			err  error
		}{
			path: "",
			err:  fmt.Errorf("failed parsing absolute path ../../internal/fixtures/kube/bla.yaml"),
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
			err  error
		}{
			path: "",
			err:  fmt.Errorf("failed parsing absolute path ../../internal/fixtures/kube"),
		},
	}
}
