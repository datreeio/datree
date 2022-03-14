package test

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/policy"
	policy_factory "github.com/datreeio/datree/bl/policy"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/pkg/errors"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/datreeio/datree/pkg/utils"

	"github.com/briandowns/spinner"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

type Evaluator interface {
	Evaluate(policyCheckData evaluation.PolicyCheckData) (evaluation.PolicyCheckResultData, error)
	SendEvaluationResult(evaluationRequestData evaluation.EvaluationRequestData) (*cliClient.SendEvaluationResultsResponse, error)
}

type Messager interface {
	LoadVersionMessages(cliVersion string) chan *messager.VersionMessage
}

type K8sValidator interface {
	ValidateResources(filesConfigurations chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile)
	InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string)
	GetK8sFiles(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.FileConfigurations)
}

type TestCommandFlags struct {
	Output               string
	K8sVersion           string
	IgnoreMissingSchemas bool
	OnlyK8sFiles         bool
	PolicyName           string
	SchemaLocations      []string
}

// TestCommandFlags constructor
func NewTestCommandFlags() *TestCommandFlags {
	return &TestCommandFlags{
		Output:               "",
		K8sVersion:           "",
		IgnoreMissingSchemas: false,
		OnlyK8sFiles:         false,
		PolicyName:           "",
		SchemaLocations:      make([]string, 0),
	}
}

//

func (flags *TestCommandFlags) Validate() error {
	outputValue := flags.Output

	if outputValue != "" {
		if (outputValue != "simple") && (outputValue != "json") && (outputValue != "yaml") && (outputValue != "xml") {

			return fmt.Errorf("Invalid --output option - %q\n"+
				"Valid output values are - simple, yaml, json, xml\n", outputValue)
		}
	}

	err := validateK8sVersionFormatIfProvided(flags.K8sVersion)

	if err != nil {
		return err
	}

	return nil

}

type EvaluationPrinter interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintMessage(messageText string, messageColor string)
	PrintPromptMessage(promptMessage string)
	PrintEvaluationSummary(evaluationSummary printer.EvaluationSummary, k8sVersion string)
	SetTheme(theme *printer.Theme)
}

type Reader interface {
	FilterFiles(paths []string) ([]string, error)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.LocalConfig, error)
}

var ViolationsFoundError = errors.New("")

type CliClient interface {
	RequestEvaluationPrerunData(token string) (*cliClient.EvaluationPrerunDataResponse, error)
}

type TestCommandData struct {
	Output                string
	K8sVersion            string
	IgnoreMissingSchemas  bool
	OnlyK8sFiles          bool
	Policy                policy.Policy
	SchemaLocations       []string
	Token                 string
	RegistrationURL       string
	PromptRegistrationURL string
}

type TestCommandContext struct {
	CliVersion   string
	LocalConfig  LocalConfig
	Evaluator    Evaluator
	Messager     Messager
	K8sValidator K8sValidator
	Printer      EvaluationPrinter
	Reader       Reader
	CliClient    CliClient
}

func LoadVersionMessages(ctx *TestCommandContext, args []string, cmd *cobra.Command) error {
	outputFlag, _ := cmd.Flags().GetString("output")
	if (outputFlag != "json") && (outputFlag != "yaml") && (outputFlag != "xml") {

		messages := ctx.Messager.LoadVersionMessages(ctx.CliVersion)
		for msg := range messages {
			ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
		}
	}
	return nil
}

func SetSilentMode(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
}

func New(ctx *TestCommandContext) *cobra.Command {
	testCommandFlags := NewTestCommandFlags()
	testCommand := &cobra.Command{
		Use:   "test <pattern>",
		Short: "Execute static analysis for given <pattern>",
		Long:  "Execute static analysis for given <pattern>. Input should be glob or `-` for stdin",
		Example: utils.Example(`
		# Test the configuration using file path
		datree test kube-prod/deployment.yaml

		# Test the configuration using glob pattern
		datree test kube-*/*.yaml

		# Test the configuration by sending manifests through stdin
		cat kube-prod/deployment.yaml | datree test -
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				errMessage := "Requires at least 1 arg\n"
				return fmt.Errorf(errMessage)
			}
			err := testCommandFlags.Validate()
			if err != nil {
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return LoadVersionMessages(ctx, args, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			SetSilentMode(cmd)
			var err error = nil
			defer func() {
				if err != nil {
					ctx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()

			err = testCommandFlags.Validate()
			if err != nil {
				return err
			}

			localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
			if err != nil {
				return err
			}

			evaluationPrerunData, err := ctx.CliClient.RequestEvaluationPrerunData(localConfigContent.Token)

			if err != nil {
				return err
			}

			testCommandOptions, err := GenerateTestCommandData(testCommandFlags, localConfigContent, evaluationPrerunData)
			if err != nil {
				return err
			}

			err = Test(ctx, args, testCommandOptions)
			if err != nil {
				return err
			}
			return nil
		},
	}
	testCommandFlags.AddFlags(testCommand)
	return testCommand
}

// AddFlags registers flags for a cli
func (flags *TestCommandFlags) AddFlags(cmd *cobra.Command) {
	if ciContext.Extract() != nil && ciContext.Extract().CIMetadata.ShouldHideEmojis {
		cmd.Flags().StringVarP(&flags.Output, "output", "o", "simple", "Define output format")
	} else {
		cmd.Flags().StringVarP(&flags.Output, "output", "o", "", "Define output format")
	}
	cmd.Flags().StringVarP(&flags.K8sVersion, "schema-version", "s", "", "Set kubernetes version to validate against. Defaults to 1.18.0")
	cmd.Flags().StringVarP(&flags.PolicyName, "policy", "p", "", "Policy name to run against")
	cmd.Flags().BoolVarP(&flags.OnlyK8sFiles, "only-k8s-files", "", false, "Evaluate only valid yaml files with the properties 'apiVersion' and 'kind'. Ignore everything else")

	// kubeconform flag
	cmd.Flags().StringArrayVarP(&flags.SchemaLocations, "schema-location", "", []string{}, "Override schemas location search path (can be specified multiple times)")
	cmd.Flags().BoolVarP(&flags.IgnoreMissingSchemas, "ignore-missing-schemas", "", false, "Ignore missing schemas when executing schema validation step")
}

func GenerateTestCommandData(testCommandFlags *TestCommandFlags, localConfigContent *localConfig.LocalConfig, evaluationPrerunDataResp *cliClient.EvaluationPrerunDataResponse) (*TestCommandData, error) {
	k8sVersion := testCommandFlags.K8sVersion
	if k8sVersion == "" {
		k8sVersion = localConfigContent.SchemaVersion
	}
	if k8sVersion == "" {
		k8sVersion = evaluationPrerunDataResp.DefaultK8sVersion
	}

	if k8sVersion == "" {
		k8sVersion = "1.18.0"
	}

	policy, err := policy_factory.CreatePolicy(evaluationPrerunDataResp.PoliciesJson, testCommandFlags.PolicyName)
	if err != nil {
		return nil, err
	}

	testCommandOptions := &TestCommandData{Output: testCommandFlags.Output,
		K8sVersion:            k8sVersion,
		IgnoreMissingSchemas:  testCommandFlags.IgnoreMissingSchemas,
		OnlyK8sFiles:          testCommandFlags.OnlyK8sFiles,
		Policy:                policy,
		SchemaLocations:       testCommandFlags.SchemaLocations,
		Token:                 localConfigContent.Token,
		RegistrationURL:       evaluationPrerunDataResp.RegistrationURL,
		PromptRegistrationURL: evaluationPrerunDataResp.PromptRegistrationURL,
	}

	return testCommandOptions, nil
}

func validateK8sVersionFormatIfProvided(k8sVersion string) error {
	if k8sVersion == "" {
		return nil
	}

	var isK8sVersionInCorrectFormat, _ = regexp.MatchString(`^[0-9]+\.[0-9]+\.[0-9]+$`, k8sVersion)
	if isK8sVersionInCorrectFormat {
		return nil
	} else {
		return fmt.Errorf("The specified schema-version %q is not in the correct format.\n"+
			"Make sure you are following the semantic versioning format <MAJOR>.<MINOR>.<PATCH>\n"+
			"Read more about kubernetes versioning: https://kubernetes.io/releases/version-skew-policy/#supported-versions", k8sVersion)
	}
}

func Test(ctx *TestCommandContext, paths []string, prerunData *TestCommandData) error {

	if paths[0] == "-" {
		if len(paths) > 1 {
			return fmt.Errorf(fmt.Sprintf("Unexpected args: [%s]", strings.Join(paths[1:], ",")))
		}
		tempFile, err := os.CreateTemp("", "datree_temp_*.yaml")
		if err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())

		if _, err := io.Copy(tempFile, os.Stdin); err != nil {
			return err
		}
		paths = []string{tempFile.Name()}
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

	if prerunData.Output == "simple" {
		ctx.Printer.SetTheme(printer.CreateSimpleTheme())
	}

	evaluationResultData, err := evaluate(ctx, filesPaths, prerunData)
	if err != nil {
		return err
	}

	results := evaluationResultData.FormattedResults
	passedPolicyCheckCount := 0
	if results.EvaluationResults != nil {
		passedPolicyCheckCount = results.EvaluationResults.Summary.TotalPassedCount
	}

	validationManager := evaluationResultData.ValidationManager

	passedYamlValidationCount := filesPathsLen - validationManager.InvalidYamlFilesCount()

	evaluationSummary := printer.EvaluationSummary{
		FilesCount:                filesPathsLen,
		RulesCount:                evaluationResultData.RulesCount,
		PassedYamlValidationCount: passedYamlValidationCount,
		PassedK8sValidationCount:  validationManager.ValidK8sFilesConfigurationsCount(),
		ConfigsCount:              validationManager.ValidK8sConfigurationsCount(),
		PassedPolicyCheckCount:    passedPolicyCheckCount,
	}

	err = evaluation.PrintResults(results, validationManager.InvalidYamlFiles(), validationManager.InvalidK8sFiles(), evaluationSummary, prerunData.RegistrationURL, prerunData.Output, ctx.Printer, prerunData.K8sVersion, prerunData.Policy.Name)

	if evaluationResultData.PromptMessage != "" {
		ctx.Printer.PrintPromptMessage(evaluationResultData.PromptMessage)
		answer, _, err := keyboard.GetSingleKey()

		if err != nil {
			fmt.Println("Failed to get prompt answer")
			return err
		}

		if strings.ToLower(string(answer)) != "n" {
			openBrowser(prerunData.PromptRegistrationURL)
		}
	}

	if err != nil {
		return err
	}

	if wereViolationsFound(validationManager, &results) {
		return ViolationsFoundError
	}

	return nil
}

type EvaluationResultData struct {
	ValidationManager *ValidationManager
	RulesCount        int
	FormattedResults  evaluation.FormattedResults
	PromptMessage     string
}

func evaluate(ctx *TestCommandContext, filesPaths []string, prerunData *TestCommandData) (EvaluationResultData, error) {
	isInteractiveMode := (prerunData.Output != "json") && (prerunData.Output != "yaml") && (prerunData.Output != "xml")

	var _spinner *spinner.Spinner
	if isInteractiveMode && prerunData.Output != "simple" {
		_spinner = createSpinner(" Loading...", "cyan")
		_spinner.Start()
	}

	defer func() {
		if _spinner != nil {
			_spinner.Stop()
		}
	}()

	validationManager := &ValidationManager{}

	ctx.K8sValidator.InitClient(prerunData.K8sVersion, prerunData.IgnoreMissingSchemas, prerunData.SchemaLocations)

	concurrency := 100

	validYamlConfigurationsChan, invalidYamlFilesChan := files.ExtractFilesConfigurations(filesPaths, concurrency)

	validationManager.AggregateInvalidYamlFiles(invalidYamlFilesChan)

	if prerunData.OnlyK8sFiles {
		var ignoredYamlFilesChan chan *extractor.FileConfigurations
		validYamlConfigurationsChan, ignoredYamlFilesChan = ctx.K8sValidator.GetK8sFiles(validYamlConfigurationsChan, concurrency)
		validationManager.AggregateIgnoredYamlFiles(ignoredYamlFilesChan)
	}

	validK8sFilesConfigurationsChan, invalidK8sFilesChan := ctx.K8sValidator.ValidateResources(validYamlConfigurationsChan, concurrency)

	validationManager.AggregateInvalidK8sFiles(invalidK8sFilesChan)

	validationManager.AggregateValidK8sFiles(validK8sFilesConfigurationsChan)

	policyName := prerunData.Policy.Name

	policyCheckData := evaluation.PolicyCheckData{
		FilesConfigurations: validationManager.ValidK8sFilesConfigurations(),
		IsInteractiveMode:   isInteractiveMode,
		PolicyName:          policyName,
		Policy:              prerunData.Policy,
	}

	emptyEvaluationResultData := EvaluationResultData{nil, 0, evaluation.FormattedResults{}, ""}

	policyCheckResultData, err := ctx.Evaluator.Evaluate(policyCheckData)
	if err != nil {
		return emptyEvaluationResultData, err
	}

	var failedYamlFiles []string
	if validationManager.InvalidYamlFilesCount() > 0 {
		for _, invalidYamlFile := range validationManager.InvalidYamlFiles() {
			failedYamlFiles = append(failedYamlFiles, invalidYamlFile.Path)
		}
	}

	var failedK8sFiles []string
	if validationManager.InvalidK8sFilesCount() > 0 {
		for _, invalidK8sFile := range validationManager.InvalidK8sFiles() {
			failedK8sFiles = append(failedK8sFiles, invalidK8sFile.Path)
		}
	}

	ciContext := ciContext.Extract()

	evaluationRequestData := evaluation.EvaluationRequestData{
		CliId:              prerunData.Token,
		CliVersion:         ctx.CliVersion,
		K8sVersion:         prerunData.K8sVersion,
		PolicyName:         policyName,
		CiContext:          ciContext,
		RulesData:          policyCheckResultData.RulesData,
		FilesData:          policyCheckResultData.FilesData,
		FailedYamlFiles:    failedYamlFiles,
		FailedK8sFiles:     failedK8sFiles,
		PolicyCheckResults: policyCheckResultData.RawResults,
	}

	sendEvaluationResultsResponse, err := ctx.Evaluator.SendEvaluationResult(evaluationRequestData)
	if err != nil {
		return emptyEvaluationResultData, err
	}

	evaluationResultData := EvaluationResultData{
		ValidationManager: validationManager,
		RulesCount:        policyCheckResultData.RulesCount,
		FormattedResults:  policyCheckResultData.FormattedResults,
		PromptMessage:     sendEvaluationResultsResponse.PromptMessage,
	}

	return evaluationResultData, err
}

func wereViolationsFound(validationManager *ValidationManager, results *evaluation.FormattedResults) bool {
	return (validationManager.InvalidYamlFilesCount() > 0 || validationManager.InvalidK8sFilesCount() > 0 || results.EvaluationResults.Summary.TotalFailedRules > 0)
}
