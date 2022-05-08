package validate_yaml

import (
	"errors"
	"fmt"
	"strings"

	pkgExtractor "github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

type IPrinter interface {
	PrintFilename(title string)
	PrintYamlValidationErrors(validationErrors []error)
	PrintYamlValidationSummary(passedFiles int, allFiles int)
	PrintMessage(messageText string, messageColor string)
}

type IReader interface {
	FilterFiles(paths []string) ([]string, error)
}

type IExtractor interface {
	ExtractConfigurationsFromYamlFile(path string) (*[]pkgExtractor.Configuration, string, *pkgExtractor.InvalidFile)
}

type ValidateYamlCommandContext struct {
	Printer   IPrinter
	Reader    IReader
	Extractor IExtractor
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

			filesPaths, err := ctx.Reader.FilterFiles(args)
			if err != nil {
				return err
			}

			filesCount := len(filesPaths)
			if filesCount == 0 {
				err = fmt.Errorf("No files detected")
				return err
			}

			invalidYamlFiles := ValidateFiles(ctx.Extractor, filesPaths)
			PrintValidationResults(ctx.Printer, invalidYamlFiles, filesCount)

			if len(invalidYamlFiles) > 0 {
				return YamlNotValidError
			}

			return nil
		},
	}
}

func ValidateFiles(extractor IExtractor, filesPaths []string) []*pkgExtractor.InvalidFile {
	var invalidYamlFiles []*pkgExtractor.InvalidFile
	for _, filePath := range filesPaths {
		_, _, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(filePath)
		if invalidYamlFile != nil {
			invalidYamlFiles = append(invalidYamlFiles, invalidYamlFile)
		}
	}

	return invalidYamlFiles
}

func PrintValidationResults(printer IPrinter, invalidFiles []*pkgExtractor.InvalidFile, filesCount int) {
	for _, invalidFile := range invalidFiles {
		printer.PrintFilename(invalidFile.Path)
		printer.PrintYamlValidationErrors(invalidFile.ValidationErrors)
	}

	// print summary
	validFilesCount := filesCount - len(invalidFiles)
	printer.PrintYamlValidationSummary(validFilesCount, filesCount)
}
