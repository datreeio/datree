package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNetworkError(t *testing.T) {
	assert.Equal(t, IsNetworkError(errors.New("blocked")), true)
	assert.Equal(t, IsNetworkError(errors.New("Blocked")), true)
	assert.Equal(t, IsNetworkError(errors.New("network Error")), true)
	assert.Equal(t, IsNetworkError(errors.New("CONNECTION REFUSED")), true)
	assert.Equal(t, IsNetworkError(errors.New("no such host")), true)
	assert.Equal(t, IsNetworkError(errors.New("i/o timeout")), true)
	assert.Equal(t, IsNetworkError(errors.New("server misbehaving")), true)
	assert.Equal(t, IsNetworkError(errors.New("this is not an error")), false)
}
