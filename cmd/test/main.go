package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
	CreateEvaluation(cliId string, cliVersion string, k8sVersion string, policyName string) (*cliClient.CreateEvaluationResponse, error)
	UpdateFailedYamlValidation(invalidFiles []*validation.InvalidYamlFile, evaluationId int, stopEvaluation bool) error
	UpdateFailedK8sValidation(invalidFiles []*validation.InvalidK8sFile, evaluationId int, stopEvaluation bool) error
}

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type K8sValidator interface {
	ValidateResources(filesConfigurations chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *validation.InvalidK8sFile)
	InitClient(k8sVersion string, ignoreMissingSchemas bool)
}

type TestCommandFlags struct {
	Output               string
	K8sVersion           string
	IgnoreMissingSchemas bool
	PolicyName           string
}

type EvaluationPrinter interface {
	SimplePrintWarnings(warnings []printer.Warning)
	PrintWarnings(warnings []printer.Warning)
	SimplePrintSummaryTable(summary printer.Summary)
	PrintSummaryTable(summary printer.Summary)
	PrintMessage(messageText string, messageColor string)
	SimplePrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string)
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
		Use:   "test <pattern>",
		Short: "Execute static analysis for given <pattern>",
		Long:  "Execute static analysis for given <pattern>. Input should be glob or `-` for stdin",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				errMessage := "Requires at least 1 arg\n"
				cmd.Usage()
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error = nil
			defer func() {
				if err != nil {
					ctx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()

			outputFlag, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			k8sVersion, err := cmd.Flags().GetString("schema-version")
			if err != nil {
				return err
			}

			ignoreMissingSchemas, err := cmd.Flags().GetBool("ignore-missing-schemas")
			if err != nil {
				return err
			}

			policy, err := cmd.Flags().GetString("policy")
			if err != nil {
				return err
			}

			testCommandFlags := TestCommandFlags{Output: outputFlag, K8sVersion: k8sVersion, IgnoreMissingSchemas: ignoreMissingSchemas, PolicyName: policy}

			err = test(ctx, args, testCommandFlags)
			if err != nil {
				return err
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	testCommand.Flags().StringP("output", "o", "", "Define output format")
	testCommand.Flags().StringP("schema-version", "s", "", "Set kubernetes version to validate against. Defaults to 1.18.0")
	testCommand.Flags().StringP("policy", "p", "", "Policy name to run against")
	testCommand.Flags().BoolP("ignore-missing-schemas", "", false, "Ignore missing schemas when executing schema validation step")
	return testCommand
}

func test(ctx *TestCommandContext, paths []string, flags TestCommandFlags) error {
	localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return err
	}

	if paths[0] == "-" {
		if len(paths) > 1 {
			return fmt.Errorf(fmt.Sprintf("Unexpected args: [%s]", strings.Join(paths[1:], ",")))
		}
		tempFile, err := ioutil.TempFile("", "datree_temp_*.yaml")
		if err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())

		if _, err := io.Copy(tempFile, os.Stdin); err != nil {
			return err
		}
		paths = []string{tempFile.Name()}
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
		return err
	}
	filesPathsLen := len(filesPaths)
	if filesPathsLen == 0 {
		noFilesErr := fmt.Errorf("No files detected")
		return noFilesErr
	}

	var _spinner *spinner.Spinner
	if isInteractiveMode == true {
		color := "cyan"
		if flags.Output == "simple" {
			color = ""
		}
		_spinner = createSpinner(" Loading...", color)
		_spinner.Start()
	}

	concurrency := 100

	validYamlFilesConfigurationsChan, invalidYamlFilesChan := files.ExtractFilesConfigurations(filesPaths, concurrency)

	createEvaluationResponse, err := ctx.Evaluator.CreateEvaluation(localConfigContent.CliId, ctx.CliVersion, flags.K8sVersion, flags.PolicyName)
	if err != nil {
		return err
	}

	ctx.K8sValidator.InitClient(createEvaluationResponse.K8sVersion, flags.IgnoreMissingSchemas)
	validK8sFilesConfigurationsChan, invalidK8sFilesChan := ctx.K8sValidator.ValidateResources(validYamlFilesConfigurationsChan, concurrency)

	invalidYamlFiles := aggregateInvalidYamlFiles(invalidYamlFilesChan)

	invalidYamlFilesLen := len(invalidYamlFiles)

	stopEvaluation := invalidYamlFilesLen == filesPathsLen
	err = ctx.Evaluator.UpdateFailedYamlValidation(invalidYamlFiles, createEvaluationResponse.EvaluationId, stopEvaluation)
	if err != nil {
		return err
	}

	invalidK8sFiles := aggregateInvalidK8sFiles(invalidK8sFilesChan)

	invalidK8sFilesLen := len(invalidK8sFiles)
	stopEvaluation = invalidYamlFilesLen+invalidK8sFilesLen == filesPathsLen

	if len(invalidK8sFiles) > 0 {
		err = ctx.Evaluator.UpdateFailedK8sValidation(invalidK8sFiles, createEvaluationResponse.EvaluationId, stopEvaluation)
		if err != nil {
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

	err = evaluation.PrintResults(results, invalidYamlFiles, invalidK8sFiles, evaluationSummary, fmt.Sprintf("https://app.datree.io/login?cliId=%s", localConfigContent.CliId), flags.Output, ctx.Printer, createEvaluationResponse.K8sVersion, createEvaluationResponse.PolicyName)

	var invocationFailedErr error = nil

	if err != nil {
		invocationFailedErr = err
	} else if len(invalidYamlFiles) > 0 || len(invalidK8sFiles) > 0 || results.Summary.TotalFailedRules > 0 {
		invocationFailedErr = fmt.Errorf("")
	}

	return invocationFailedErr
}
