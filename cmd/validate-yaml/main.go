package validate_yaml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/utils"
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
	ValidateResources(filesConfigurations chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile, chan *validation.FileWithWarning)
	InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string)
	GetK8sFiles(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.FileConfigurations)
}

type IPrinter interface {
	PrintTitle(title string)
	PrintYamlValidationErrors(validationErrors []error)
	PrintYamlValidationSummary(passedFiles int, allFiles int)
	PrintMessage(messageText string, messageColor string)
}

type Reader interface {
	FilterFiles(paths []string) ([]string, error)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.LocalConfig, error)
}

type CliClient interface {
	RequestEvaluationPrerunData(token string, isCi bool) (*cliClient.EvaluationPrerunDataResponse, error)
	AddFlags(flags map[string]interface{})
}

type ValidateYamlCommandContext struct {
	CliVersion     string
	CiContext      *ciContext.CIContext
	LocalConfig    LocalConfig
	Evaluator      Evaluator
	Messager       Messager
	K8sValidator   K8sValidator
	Printer        IPrinter
	Reader         Reader
	CliClient      CliClient
	FilesExtractor files.FilesExtractorInterface
}

var YamlNotValidError = errors.New("")

func SetSilentMode(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
}

func New(ctx *ValidateYamlCommandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "validate-yaml <files>",
		Short: "validates yaml files structure",
		Long:  "Validates yaml files <files> structure.",
		Example: utils.Example(`
		# Validate yaml files using file path
		datree validate-yaml kube-prod/deployment.yaml

		# Validate yaml files using glob pattern
		datree validate-yaml kube-*/*.yaml

		# Validate yaml files by sending manifests through stdin
		cat kube-prod/deployment.yaml | datree validate-yaml -
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			SetSilentMode(cmd)
			var err error = nil

			defer func() {
				if err != nil {
					ctx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()
			var invalidYamlFiles []*extractor.InvalidFile

			filesPaths, err := ctx.Reader.FilterFiles(args)
			if err != nil {
				return err
			}

			filesCount := len(filesPaths)
			if filesCount == 0 {
				err = fmt.Errorf("No files detected")
				return err
			}

			// validate files
			for _, filePath := range filesPaths {
				_, _, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(filePath)
				if invalidYamlFile != nil {
					invalidYamlFiles = append(invalidYamlFiles, invalidYamlFile)
				}
			}

			// print files with errors
			for _, invalidYamlFile := range invalidYamlFiles {
				ctx.Printer.PrintTitle(invalidYamlFile.Path)
				ctx.Printer.PrintYamlValidationErrors(invalidYamlFile.ValidationErrors)
			}

			// print summary
			validFilesCount := filesCount - len(invalidYamlFiles)
			ctx.Printer.PrintYamlValidationSummary(validFilesCount, filesCount)

			if len(invalidYamlFiles) > 0 {
				return YamlNotValidError
			}

			return nil
		},
	}
}
