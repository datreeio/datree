package files

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAbsolutePaths(t *testing.T) {
	test_existed_file(t)
	test_directory_file(t)
	test_not_existed_file(t)
}

func test_existed_file(t *testing.T) {
	paths := []string{"../../internal/fixtures/kube/pass-all.yaml"}
	pathsChan, errors := ToAbsolutePaths(paths)
	assert.Equal(t, "/Users/noaabarki/Dev/datree/internal/fixtures/kube/pass-all.yaml", <-pathsChan)
	assert.Equal(t, nil, <-errors)
}

func test_not_existed_file(t *testing.T) {
	paths := []string{"../../internal/fixtures/kube/bla.yaml"}
	pathsChan, errors := ToAbsolutePaths(paths)
	assert.Equal(t, "", <-pathsChan)
	err := <-errors
	assert.Equal(t, fmt.Errorf("failed parsing absolute path ../../internal/fixtures/kube/bla.yaml"), err)
}

func test_directory_file(t *testing.T) {
	paths := []string{"../../internal/fixtures/kube"}
	pathsChan, errors := ToAbsolutePaths(paths)
	assert.Equal(t, "", <-pathsChan)
	err := <-errors
	assert.Equal(t, fmt.Errorf("failed parsing absolute path ../../internal/fixtures/kube"), err)
}
