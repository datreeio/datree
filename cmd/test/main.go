package test

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/spf13/cobra"
)

type Evaluator interface {
	Evaluate(filesConfigurationsChan []*extractor.FileConfigurations, evaluationId int) (*evaluation.EvaluationResults, error)
	CreateEvaluation(cliId string, cliVersion string, k8sVersion string) (*cliClient.CreateEvaluationResponse, error)
	UpdateFailedYamlValidation(invalidFiles []*validation.InvalidYamlFile, evaluationId int, stopEvaluation bool) error
	UpdateFailedK8sValidation(invalidFiles []*validation.InvalidK8sFile, evaluationId int, stopEvaluation bool) error
}

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type K8sValidator interface {
	ValidateResources(filesConfigurations chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *validation.InvalidK8sFile)
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

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
}

type TestCommandContext struct {
	CliVersion   string
	LocalConfig  LocalConfig
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
	localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return err
	}

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

	filesPaths, err := ctx.Reader.FilterFiles(paths)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	filesPathsLen := len(filesPaths)
	if filesPathsLen == 0 {
		noFilesErr := fmt.Errorf("No files detected")
		fmt.Println(noFilesErr.Error())
		return noFilesErr
	}

	var _spinner *spinner.Spinner
	if isInteractiveMode == true {
		_spinner = createSpinner(" Loading...", "cyan")
		_spinner.Start()
	}

	concurrency := 100

	validYamlFilesConfigurationsChan, invalidYamlFilesChan := files.ExtractFilesConfigurations(filesPaths, concurrency)

	createEvaluationResponse, err := ctx.Evaluator.CreateEvaluation(localConfigContent.CliId, ctx.CliVersion, flags.K8sVersion)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ctx.K8sValidator.InitClient(createEvaluationResponse.K8sVersion)
	validK8sFilesConfigurationsChan, invalidK8sFilesChan := ctx.K8sValidator.ValidateResources(validYamlFilesConfigurationsChan, concurrency)

	invalidYamlFiles := aggregateInvalidYamlFiles(invalidYamlFilesChan)

	invalidYamlFilesLen := len(invalidYamlFiles)

	stopEvaluation := invalidYamlFilesLen == filesPathsLen
	err = ctx.Evaluator.UpdateFailedYamlValidation(invalidYamlFiles, createEvaluationResponse.EvaluationId, stopEvaluation)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	invalidK8sFiles := aggregateInvalidK8sFiles(invalidK8sFilesChan)

	invalidK8sFilesLen := len(invalidK8sFiles)
	stopEvaluation = invalidYamlFilesLen+invalidK8sFilesLen == filesPathsLen

	if len(invalidK8sFiles) > 0 {
		err = ctx.Evaluator.UpdateFailedK8sValidation(invalidK8sFiles, createEvaluationResponse.EvaluationId, stopEvaluation)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	var validK8sFilesConfigurations []*extractor.FileConfigurations
	for fileConfigurations := range validK8sFilesConfigurationsChan {
		validK8sFilesConfigurations = append(validK8sFilesConfigurations, fileConfigurations)
	}

	results, err := ctx.Evaluator.Evaluate(validK8sFilesConfigurations, createEvaluationResponse.EvaluationId)

	if _spinner != nil && isInteractiveMode == true {
		_spinner.Stop()
	}

	passedPolicyCheckCount := 0
	if results != nil {
		passedPolicyCheckCount = results.Summary.TotalPassedCount
	}

	passedYamlValidationCount := filesPathsLen - invalidYamlFilesLen
	passedK8sValidationCount := passedYamlValidationCount - invalidK8sFilesLen

	configsCount := countConfigurations(validK8sFilesConfigurations)

	evaluationSummary := printer.EvaluationSummary{
		FilesCount:                filesPathsLen,
		RulesCount:                createEvaluationResponse.RulesCount,
		PassedYamlValidationCount: passedYamlValidationCount,
		PassedK8sValidationCount:  passedK8sValidationCount,
		ConfigsCount:              configsCount,
		PassedPolicyCheckCount:    passedPolicyCheckCount,
	}

	err = evaluation.PrintResults(results, invalidYamlFiles, invalidK8sFiles, evaluationSummary, fmt.Sprintf("https://app.datree.io/login?cliId=%s", localConfigContent.CliId), flags.Output, ctx.Printer, createEvaluationResponse.K8sVersion)

	var invocationFailedErr error = nil

	if err != nil {
		fmt.Println(err.Error())
		invocationFailedErr = err
	} else if len(invalidYamlFiles) > 0 || len(invalidK8sFiles) > 0 || results.Summary.TotalFailedRules > 0 {
		invocationFailedErr = fmt.Errorf("Evaluation failed")
	}

	return invocationFailedErr
}
