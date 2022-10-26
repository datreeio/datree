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
	RequestEvaluationPrerunData(token string, isCi bool) (*cliClient.EvaluationPrerunDataResponse, error)
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
		Args: func(cmd *cobra.Command, args []string) error {
			err := utils.ValidateStdinPathArgument(args)
			if err != nil {
				return err
			}
			return testCommandFlags.Validate()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return test.LoadVersionMessages(testCtx, args, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// fmt.Println(cmd.Flags().Changed("save-rendered"))
			test.SetSilentMode(cmd)
			var err error = nil
			defer func() {
				if err != nil {
					testCtx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()

			out, err := kustomizeCtx.CommandRunner.ExecuteKustomizeBin(args)
			if err != nil {
				return err
			}

			var tempFilename string
			if testCommandFlags.SaveRendered != "" {
				tempFilename, err = test.SaveRenderedFile(testCommandFlags.SaveRendered, out)
				if err != nil {
					return err
				}
			} else {
				tempFilename, err = kustomizeCtx.CommandRunner.CreateTempFile("datree_kustomize", out)
				if err != nil {
					return err
				}
				defer os.Remove(tempFilename)
			}

			err = test.TestWrapper(testCtx, []string{tempFilename}, testCommandFlags)
			if err != nil {
				return err
			}
			return nil
		},
	}
	testCommandFlags.AddFlags(kustomizeTestCommand)

	kustomizeCommand := &cobra.Command{
		Use:   "kustomize",
		Short: "Render resources defined in a kustomization.yaml file and run a policy check against them",
	}

	kustomizeCommand.AddCommand(kustomizeTestCommand)

	return kustomizeCommand
}
