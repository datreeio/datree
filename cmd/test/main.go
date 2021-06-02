package test

import (
	"fmt"
	"time"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/spf13/cobra"
)

type Evaluator interface {
	Evaluate(validFilesPathsChan chan string, invalidFilesPaths chan *validation.InvalidFile, evaluationId int) (*evaluation.EvaluationResults, []*validation.InvalidFile, []*extractor.FileConfiguration, []*evaluation.Error, error)
	CreateEvaluation(cliId string, cliVersion string, k8sVersion string) (*cliClient.CreateEvaluationResponse, error)
}

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type K8sValidator interface {
	ValidateResources(paths []string) (chan string, chan *validation.InvalidFile, chan error)
	InitClient(k8sVersion string)
}

type TestCommandFlags struct {
	Output     string
	K8sVersion string
}

type EvaluationPrinter interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintMessage(messageText string, messageColor string)
	PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string)
}

type Reader interface {
	FilterFiles(paths []string) ([]string, error)
}

type TestCommandContext struct {
	CliVersion   string
	LocalConfig  *localConfig.LocalConfiguration
	Evaluator    Evaluator
	Messager     Messager
	K8sValidator K8sValidator
	Printer      EvaluationPrinter
	Reader       Reader
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

			k8sVersion, err := cmd.Flags().GetString("schema-version")
			if err != nil {
				fmt.Println(err)
				return err
			}

			testCommandFlags := TestCommandFlags{Output: outputFlag, K8sVersion: k8sVersion}
			return test(ctx, args, testCommandFlags)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	testCommand.Flags().StringP("output", "o", "", "Define output format")
	testCommand.Flags().StringP("schema-version", "s", "", "Set kubernetes version to validate against. Defaults to 1.18.0")
	return testCommand
}

func test(ctx *TestCommandContext, paths []string, flags TestCommandFlags) error {
	isInteractiveMode := (flags.Output != "json") && (flags.Output != "yaml")

	if isInteractiveMode == true {
		messages := make(chan *messager.VersionMessage, 1)
		go ctx.Messager.LoadVersionMessages(messages, ctx.CliVersion)
		defer func() {
			msg, ok := <-messages
			if ok {
				ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
			}
		}()
	}

	filePaths, err := ctx.Reader.FilterFiles(paths)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if len(filePaths) == 0 {
		noFilesErr := fmt.Errorf("No files detected")
		fmt.Println(noFilesErr.Error())
		return noFilesErr
	}

	var _spinner *spinner.Spinner
	if isInteractiveMode == true {
		_spinner = createSpinner(" Loading...", "cyan")
		_spinner.Start()
	}

	createEvaluationResponse, err := ctx.Evaluator.CreateEvaluation(ctx.LocalConfig.CliId, ctx.CliVersion, flags.K8sVersion)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ctx.K8sValidator.InitClient(createEvaluationResponse.K8sVersion)
	validFilesPathsChan, invalidFilesPathsChan, errorsChan := ctx.K8sValidator.ValidateResources(paths)
	go func() {
		for err := range errorsChan {
			fmt.Println(err.Error())
		}
	}()

	results, invalidFiles, filesConfigurations, errors, err := ctx.Evaluator.Evaluate(validFilesPathsChan, invalidFilesPathsChan, createEvaluationResponse.EvaluationId)

	if _spinner != nil && isInteractiveMode == true {
		_spinner.Stop()
	}

	passedPolicyCheckCount := 0
	if results != nil {
		passedPolicyCheckCount = results.Summary.TotalPassedCount
	}

	configsCount := countConfigurations(filesConfigurations)
	evaluationSummary := printer.EvaluationSummary{
		ConfigsCount:              configsCount,
		FilesCount:                len(paths),
		PassedYamlValidationCount: len(paths),
		PassedK8sValidationCount:  len(filesConfigurations),
		PassedPolicyCheckCount:    passedPolicyCheckCount,
	}

	err = evaluation.PrintResults(results, invalidFiles, evaluationSummary, fmt.Sprintf("https://app.datree.io/login?cliId=%s", ctx.LocalConfig.CliId), flags.Output, ctx.Printer, createEvaluationResponse.K8sVersion)

	var invocationFailedErr error = nil

	if err != nil {
		fmt.Println(err.Error())
		invocationFailedErr = err
	} else if len(errors) > 0 {
		printEvaluationErrors(errors)
		invocationFailedErr = fmt.Errorf("Invocation failed")
	} else if len(invalidFiles) > 0 || results.Summary.TotalFailedRules > 0 {
		invocationFailedErr = fmt.Errorf("Invocation failed")
	}

	return invocationFailedErr
}

func countConfigurations(filesConfigurations []*extractor.FileConfiguration) int {
	totalConfigs := 0

	for _, fileConfiguration := range filesConfigurations {
		totalConfigs += len(fileConfiguration.Configurations)
	}

	return totalConfigs
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
