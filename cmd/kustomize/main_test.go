package kustomize

import (
	"testing"

	"github.com/datreeio/datree/cmd/test"
)

func TestKustomizeTestCommand(t *testing.T) {
	t.Skip("Skipping test")

}

func test_kustomize_run_method(t *testing.T) {

}

func invokeKustomizeRunMethod(ctx *test.TestCommandContext, args []string) error {
	cmd := New(ctx, &KustomizeContext{})
	cmd.SetArgs(args)
	err := cmd.RunE(cmd, args)
	return err
}
