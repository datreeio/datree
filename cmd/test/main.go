package test

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/spf13/cobra"
)

type Evaluator interface {
	Evaluate(paths <-chan string, cliId string, cliVersion string) (*evaluation.EvaluationResults, []*evaluation.Error, error)
}

type Messager interface {
	LoadVersionMessages(cliVersion string, messages chan *messager.VersionMessage)
	HandleVersionMessage(message *messager.VersionMessage)
}

type Validator interface {
	Validate(paths []string) (<-chan string, <-chan string, <-chan error)
}

type TestCommandFlags struct {
	Output string
}

type EvaluationPrinter interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintMessage(messageText string, messageColor string)
}
type TestCommandContext struct {
	CliVersion  string
	LocalConfig *localConfig.LocalConfiguration
	Evaluator   Evaluator
	Messager    Messager
	Validator   Validator
	Printer     EvaluationPrinter
}

func New(ctx *TestCommandContext) *cobra.Command {
	testCommand := &cobra.Command{
		Use:   "test",
		Short: "Execute static analysis for pattern",
		Long:  `Execute static analysis for pattern. Input should be glob`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFlag, err := cmd.Flags().GetString("output")
			if err != nil {
				fmt.Println(err)
				return err
			}

			testCommandFlags := TestCommandFlags{Output: outputFlag}
			return test(ctx, args, testCommandFlags)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	testCommand.Flags().StringP("output", "o", "", "Define output format")
	return testCommand
}

func test(ctx *TestCommandContext, paths []string, flags TestCommandFlags) error {
	s := createSpinner(" Loading...", "cyan")
	s.Start()

	messages := make(chan *messager.VersionMessage, 1)
	go ctx.Messager.LoadVersionMessages(ctx.CliVersion, messages)

	validPaths, _, _ := ctx.Validator.Validate(paths)
	results, errors, err := ctx.Evaluator.Evaluate(validPaths, ctx.LocalConfig.CliId, ctx.CliVersion)

	s.Stop()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if results == nil {
		return fmt.Errorf("no response received")
	}

	if len(errors) > 0 {
		printEvaluationErrors(errors)
	}

	evaluation.PrintResults(results, fmt.Sprintf("https://app.datree.io/login?cliId=%s", ctx.LocalConfig.CliId), flags.Output, ctx.Printer)
	msg, ok := <-messages
	if ok {
		ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
	}

	return err
}

func createSpinner(text string, color string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = text
	s.Color(color)
	return s
}

func printEvaluationErrors(errors []*evaluation.Error) {
	fmt.Println("The following files failed:")
	for _, fileError := range errors {
		fmt.Printf("\n\tFilename: %s\n\tError: %s\n\t---------------------", fileError.Filename, fileError.Message)
	}
}
