package schema_validator

import (
	"fmt"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

type Result = gojsonschema.Result

type YamlSchemaValidator interface {
	Validate(yamlSchema string, yaml string) (*Result, error)
}

type YamlSchemaValidationPrinter interface {
	PrintYamlSchemaResults(result *Result, error error)
}

type YamlSchemaValidatorCommandContext struct {
	YamlSchemaValidator YamlSchemaValidator
	Printer             YamlSchemaValidationPrinter
}

func New(ctx *YamlSchemaValidatorCommandContext) *cobra.Command {
	schemaValidator := &cobra.Command{
		Use:    "schema-validator",
		Short:  "Execute schema validation to yaml files for given yaml schema",
		Long:   "Execute schema validation to yaml files for given yaml schema. Input should be glob or 1 yaml schema file and 1 yaml file",
		Hidden: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				errMessage := "Requires 2 args\n"
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaPath := args[0]
			yamlPath := args[1]

			_, _, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(yamlPath)
			if invalidYamlFile != nil {
				ctx.Printer.PrintYamlSchemaResults(nil, invalidYamlFile.ValidationErrors[0])
				return invalidYamlFile.ValidationErrors[0]
			}

			schemaContent, err := extractor.ReadFileContent(schemaPath)
			if err != nil {
				ctx.Printer.PrintYamlSchemaResults(nil, err)
				return err
			}
			yamlContent, err := extractor.ReadFileContent(yamlPath)
			if err != nil {
				ctx.Printer.PrintYamlSchemaResults(nil, err)
				return err
			}

			result, err := ctx.YamlSchemaValidator.Validate(schemaContent, yamlContent)
			ctx.Printer.PrintYamlSchemaResults(result, err)

			if err != nil {
				return err
			}

			return nil
		},
	}
	return schemaValidator
}
