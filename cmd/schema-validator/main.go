package schema_validator

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

type Result = gojsonschema.Result

type JSONSchemaValidator interface {
	ValidateYamlSchemaNew(yamlSchema string, yaml string) ([]jsonschema.Detailed, error)
}

type JSONSchemaValidationPrinter interface {
	PrintYamlSchemaResults(result []jsonschema.Detailed, error error)
}

type JSONSchemaValidatorCommandContext struct {
	JSONSchemaValidator JSONSchemaValidator
	Printer             JSONSchemaValidationPrinter
}

func ExtractYamlFilesContent(schemaPath string, yamlPath string) (string, string, error) {
	_, _, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(yamlPath)
	if invalidYamlFile != nil {
		return "", "", invalidYamlFile.ValidationErrors[0]
	}
	schemaContent, err := extractor.ReadFileContent(schemaPath)
	if err != nil {
		return "", "", err
	}
	yamlContent, err := extractor.ReadFileContent(yamlPath)
	if err != nil {
		return "", "", err
	}

	return schemaContent, yamlContent, nil
}

func New(ctx *JSONSchemaValidatorCommandContext) *cobra.Command {
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

			schemaContent, yamlContent, err := ExtractYamlFilesContent(schemaPath, yamlPath)
			if err != nil {
				ctx.Printer.PrintYamlSchemaResults(nil, err)
				return err
			}
			result, err := ctx.JSONSchemaValidator.ValidateYamlSchemaNew(schemaContent, yamlContent)
			ctx.Printer.PrintYamlSchemaResults(result, err)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return schemaValidator
}
