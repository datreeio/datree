package test

import (
	"fmt"
	"github.com/datreeio/datree/pkg/cliClient"
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/spf13/cobra"
)

type Evaluator interface {
	Evaluate(validFilesPathsChan chan string, invalidFilesPaths chan *validation.InvalidFile, evaluationId int) (*evaluation.EvaluationResults, []*validation.InvalidFile, []*cliClient.FileConfiguration, []*evaluation.Error, error)
	CreateEvaluation(cliId string, cliVersion string) (int, error)
}

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type K8sValidator interface {
	ValidateResources(paths []string) (chan string, chan *validation.InvalidFile, chan error)
}

type TestCommandFlags struct {
	Output string
}

type EvaluationPrinter interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintMessage(messageText string, messageColor string)
	PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary)
}
type TestCommandContext struct {
	CliVersion   string
	LocalConfig  *localConfig.LocalConfiguration
	Evaluator    Evaluator
	Messager     Messager
	K8sValidator K8sValidator
	Printer      EvaluationPrinter
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
	spinner := createSpinner(" Loading...", "cyan")
	spinner.Start()

	messages := make(chan *messager.VersionMessage, 1)
	go ctx.Messager.LoadVersionMessages(messages, ctx.CliVersion)

	validFilesPaths, invalidFilesPathsChan, errorsChan := ctx.K8sValidator.ValidateResources(paths)

	go func() {
		for err := range errorsChan {
			fmt.Println(err)
		}
	}()

	evaluationId, err := ctx.Evaluator.CreateEvaluation(ctx.LocalConfig.CliId, ctx.CliVersion)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	results, invalidFiles, filesConfigurations, errors, err := ctx.Evaluator.Evaluate(validFilesPaths, invalidFilesPathsChan, evaluationId)

	spinner.Stop()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(errors) > 0 {
		printEvaluationErrors(errors)
	}

	passedPolicyCheckCount := 0
	if results != nil {
		passedPolicyCheckCount = results.Summary.TotalPassedCount
	}

	evaluationSummary := printer.EvaluationSummary{
		FilesCount:                len(paths),
		PassedYamlValidationCount: len(paths),
		PassedK8sValidationCount:  len(filesConfigurations),
		PassedPolicyCheckCount:    passedPolicyCheckCount,
	}

	evaluation.PrintResults(results, invalidFiles, evaluationSummary, fmt.Sprintf("https://app.datree.io/login?cliId=%s", ctx.LocalConfig.CliId), flags.Output, ctx.Printer)
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
