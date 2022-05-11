package validate_yaml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/datreeio/datree/pkg/utils"
	"github.com/datreeio/datree/pkg/yamlValidator"
	"github.com/spf13/cobra"
)

type IReader interface {
	FilterFiles(paths []string) ([]string, error)
}

type ValidateYamlCommandContext struct {
	Printer   yamlValidator.IPrinter
	Reader    IReader
	Extractor yamlValidator.IExtractor
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

			invalidYamlFiles := yamlValidator.ValidateFiles(ctx.Extractor, filesPaths)
			yamlValidator.PrintValidationResults(ctx.Printer, invalidYamlFiles, filesCount)

			if len(invalidYamlFiles) > 0 {
				return YamlNotValidError
			}

			return nil
		},
	}
}
