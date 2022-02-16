package schema_validator

import (
	"fmt"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"

)

type Result = gojsonschema.Result

type JsonSchemaValidator interface {
	Validate(jsonSchema string, json string) (*Result, error)
}

type JsonSchemaValidationPrinter interface {
	PrintJsonSchemaResults(result *Result, error error)
}

type SchemaValidatorCommandContext struct {
	JsonSchemaValidator JsonSchemaValidator
	Printer             JsonSchemaValidationPrinter
}

func New(ctx *SchemaValidatorCommandContext) *cobra.Command {
	schemaValidator := &cobra.Command{
		Use:   "schema-validator",
		Short: "Execute schema validation to yaml files for given json",
		Long:  "Execute schema validation to yaml files for given json. Input should be glob or 1 yaml file and 1 json file",
		Hidden: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2{
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
				ctx.Printer.PrintJsonSchemaResults(nil ,invalidYamlFile.ValidationErrors[0])
				return invalidYamlFile.ValidationErrors[0]
			}

			schemaContent, _ := extractor.ReadFileContent(schemaPath)

			yamlContent, _ := extractor.ReadFileContent(yamlPath)

			result , err := ctx.JsonSchemaValidator.Validate(schemaContent, yamlContent)
			ctx.Printer.PrintJsonSchemaResults(result ,err)

			if err != nil {
				return err
			}

			return nil
		},
	}
	return schemaValidator
}