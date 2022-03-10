package kustomize

import (
	"os"
	"strings"

	"github.com/datreeio/datree/pkg/cliClient"

	"github.com/datreeio/datree/cmd/test"
	"github.com/datreeio/datree/pkg/executor"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

type CliClient interface {
	RequestPrerunDataForEvaluation(token string) (*cliClient.PrerunDataForEvaluationResponse, int, error)
}

type KustomizeCommandRunner interface {
	BuildCommandDescription(dir string, name string, args []string) string
	RunCommand(name string, args []string) (executor.CommandOutput, error)
	ExecuteKustomizeBin(args []string) ([]byte, error)
	CreateTempFile(tempFilePrefix string, content []byte) (string, error)
}

type KustomizeContext struct {
	CommandRunner KustomizeCommandRunner
}

const notFoundStatusCode = 404

func New(testCtx *test.TestCommandContext, kustomizeCtx *KustomizeContext) *cobra.Command {
	testCommandFlags := test.NewTestCommandFlags()
	kustomizeTestCommand := &cobra.Command{
		Use:   "test <args>",
		Short: "Execute datree test for kustomize build <args>",
		Long:  "Execute datree test for kustomize build <args>. Input should be a kustomize build directory or file.",
		Example: utils.Example(`
		# Test the kustomize build for the current directory
		datree kustomize test .

		# Test the kustomize build for some shared configuration directory
		datree kustomize test /home/config/production

		# Test the kustomize build from github
		datree kustomize test https://github.com/kubernetes-sigs/kustomize.git/examples/helloWorld?ref=v1.0.6
		`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return test.LoadVersionMessages(testCtx, args, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			test.SetSilentMode(cmd)
			var err error = nil
			defer func() {
				if err != nil {
					testCtx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()

			err = testCommandFlags.Validate()
			if err != nil {
				return err
			}

			localConfigContent, err := testCtx.LocalConfig.GetLocalConfiguration()
			if err != nil {
				return err
			}

			prerunDataForEvaluation, statusCode, err := testCtx.CliClient.RequestPrerunDataForEvaluation(localConfigContent.CliId)

			// getting prerun data can return 404 if user has a valid token but he didn't sign up yet - suppose to be ok
			// getting prerun data can return 400 if user has invalid token - we suppose to return an error
			if err != nil && statusCode != notFoundStatusCode {
				return err
			}

			testCommandOptions, err := test.GenerateTestCommandData(testCommandFlags, localConfigContent, prerunDataForEvaluation)

			out, err := kustomizeCtx.CommandRunner.ExecuteKustomizeBin(args)
			if err != nil {
				return err
			}

			tempFilename, err := kustomizeCtx.CommandRunner.CreateTempFile("datree_kustomize", out)
			if err != nil {
				return err
			}
			defer os.Remove(tempFilename)

			err = test.Test(testCtx, []string{tempFilename}, testCommandOptions)
			if err != nil {
				return err
			}
			return nil
		},
	}
	testCommandFlags.AddFlags(kustomizeTestCommand)

	kustomizeCommand := &cobra.Command{
		Use: "kustomize",
	}

	kustomizeCommand.AddCommand(kustomizeTestCommand)

	return kustomizeCommand
}
