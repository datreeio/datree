package validate_yaml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/datreeio/datree/pkg/cliClient"
	pkgExtractor "github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/datreeio/datree/pkg/yamlValidator"
	"github.com/spf13/cobra"
)

const (
	STATUS_PASSED = "passed"
	STATUS_FAILED = "failed"
)

type IReader interface {
	FilterFiles(paths []string) ([]string, error)
}

type IPrinter interface {
	GetFileNameText(title string) string
	GetYamlValidationErrorsText(validationErrors []error) string
	GetYamlValidationSummaryText(passedFiles int, allFiles int) string
	PrintMessage(messageText string, messageColor string)
}

type ICliClient interface {
	SendValidateYamlResult(request *cliClient.ValidatedYamlResult)
}

type ILocalConfig interface {
	GetLocalConfiguration() (*localConfig.LocalConfig, error)
}

type ValidateYamlCommandContext struct {
	Printer     IPrinter
	Reader      IReader
	Extractor   yamlValidator.IExtractor
	CliClient   ICliClient
	LocalConfig ILocalConfig
	CliVersion  string
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
		datree validate-yaml kube-*/*.*
		`),
		Hidden: true,
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

			filesPaths, err := ctx.Reader.FilterFiles(args)
			if err != nil {
				return err
			}

			filesCount := len(filesPaths)
			if filesCount == 0 {
				err = fmt.Errorf("No files detected")
				return err
			}

			newYamlValidator := yamlValidator.New(&yamlValidator.YamlValidatorOptions{
				Extractor: ctx.Extractor,
			})

			invalidYamlFiles := newYamlValidator.ValidateFiles(filesPaths)
			PrintValidationResults(ctx.Printer, invalidYamlFiles, filesCount)

			isValid := len(invalidYamlFiles) == 0

			SendResults(ctx.LocalConfig, ctx.CliClient, ctx.CliVersion, isValid, invalidYamlFiles, filesPaths)

			if !isValid {
				return YamlNotValidError
			}

			return nil
		},
	}
}

func PrintValidationResults(printer IPrinter, invalidFiles []*pkgExtractor.InvalidFile, filesCount int) {
	for _, invalidFile := range invalidFiles {
		printer.GetFileNameText(invalidFile.Path)
		printer.GetYamlValidationErrorsText(invalidFile.ValidationErrors)
	}

	validFilesCount := filesCount - len(invalidFiles)
	printer.GetYamlValidationSummaryText(validFilesCount, filesCount)
}

func SendResults(localConfig ILocalConfig, client ICliClient, cliVersion string, isValid bool, invalidYamlFiles []*pkgExtractor.InvalidFile, filesPaths []string) {
	osInfo := utils.NewOSInfo()
	resultFiles := prepareValidationResults(invalidYamlFiles, filesPaths)
	configData, err := localConfig.GetLocalConfiguration()

	if err != nil {
		return
	}

	status := STATUS_PASSED
	if !isValid {
		status = STATUS_FAILED
	}
	result := &cliClient.ValidatedYamlResult{
		Token:    configData.Token,
		ClientId: configData.ClientId,
		Files:    resultFiles,
		Status:   status,
		Metadata: &cliClient.Metadata{
			CliVersion:      cliVersion,
			Os:              osInfo.OS,
			KernelVersion:   osInfo.KernelVersion,
			PlatformVersion: osInfo.PlatformVersion,
		},
	}

	client.SendValidateYamlResult(result)
}

func prepareValidationResults(invalidFiles []*pkgExtractor.InvalidFile, filesPaths []string) []*cliClient.ValidatedFile {
	var validationResults []*cliClient.ValidatedFile
	filesMap := make(map[string]bool)

	for _, filename := range filesPaths {
		absoluteFilePath, _ := pkgExtractor.ToAbsolutePath(filename)
		filesMap[absoluteFilePath] = true
	}
	for _, invalidFile := range invalidFiles {
		filesMap[invalidFile.Path] = false
	}
	for filename, isValid := range filesMap {
		validationResults = append(validationResults, &cliClient.ValidatedFile{
			Path:    filename,
			IsValid: isValid,
		})
	}
	return validationResults
}
