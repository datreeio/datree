package completion

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteCommand(t *testing.T) {

	t.Run("Test Completion Bash", func(t *testing.T) {
		cmd := New()
		out = new(bytes.Buffer)

		cmd.SetArgs([]string{"bash"})
		cmd.Execute()

		got := out.(*bytes.Buffer).Bytes()
		expected := []byte(`# bash completion for completion`)

		assert.Equal(t, string(expected), string(got)[:len(expected)])
	})

	t.Run("Test Completion Zsh", func(t *testing.T) {
		cmd := New()
		out = new(bytes.Buffer)

		cmd.SetArgs([]string{"zsh"})
		cmd.Execute()

		got := out.(*bytes.Buffer).Bytes()
		expected := []byte(`#compdef completion`)

		assert.Equal(t, string(expected), string(got)[:len(expected)])
	})

	t.Run("Test Completion Fish", func(t *testing.T) {
		cmd := New()
		out = new(bytes.Buffer)

		cmd.SetArgs([]string{"fish"})
		cmd.Execute()

		got := out.(*bytes.Buffer).Bytes()
		expected := []byte(`# fish completion for completion`)

		assert.Equal(t, string(expected), string(got)[:len(expected)])
	})

	t.Run("Test Completion Powershell", func(t *testing.T) {
		cmd := New()
		out = new(bytes.Buffer)

		cmd.SetArgs([]string{"powershell"})
		cmd.Execute()

		got := out.(*bytes.Buffer).Bytes()
		expected := []byte(`# powershell completion for completion`)

		assert.Equal(t, string(expected), string(got)[:len(expected)])
	})
}
