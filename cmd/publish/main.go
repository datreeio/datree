package publish

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	jsonSchemaValidator "github.com/datreeio/datree/pkg/jsonSchemaValidator/extensions"
	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/evaluation"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

type Messager interface {
	LoadVersionMessages(cliVersion string) chan *messager.VersionMessage
}

type Printer interface {
	PrintMessage(messageText string, messageColor string)
	PrintYamlSchemaResults(errorsResult []jsonschema.Detailed, error error)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.LocalConfig, error)
}

type CliClient interface {
	PublishPolicies(policiesConfiguration files.UnknownStruct, token string) (*cliClient.PublishFailedResponse, error)
}

type JSONSchemaValidator interface {
	ValidateYamlSchema(yamlSchema string, yaml string) ([]jsonschema.Detailed, error)
}

type PublishCommandContext struct {
	CliVersion          string
	LocalConfig         LocalConfig
	Messager            Messager
	Printer             Printer
	PublishCliClient    CliClient
	FilesExtractor      files.FilesExtractorInterface
	JSONSchemaValidator JSONSchemaValidator
}

func New(ctx *PublishCommandContext) *cobra.Command {
	var localConfigContent *localConfig.LocalConfig

	publishCommand := &cobra.Command{
		Use:   "publish <fileName>",
		Short: "Publish policies configuration for given <fileName>.",
		Long:  "Publish policies configuration for given <fileName>. Input should be the path to the Policy-as-Code yaml configuration file",
		Example: utils.Example(`
		# Publish the policies configuration YAML file
		datree publish policies.yaml

		# Note You need to first enable Policy-as-Code (PaC) on the settings page in the dashboard
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				errMessage := "requires 1 arg"
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			outputFlag, _ := cmd.Flags().GetString("output")
			if !evaluation.IsFormattedOutputOption(outputFlag) {

				messages := ctx.Messager.LoadVersionMessages(ctx.CliVersion)
				for msg := range messages {
					ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
				}
			}
			var err error
			localConfigContent, err = ctx.LocalConfig.GetLocalConfiguration()
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			publishFailedResponse, err := publish(ctx, args[0], localConfigContent)
			if publishFailedResponse != nil {
				ctx.Printer.PrintMessage("Publish failed: "+publishFailedResponse.Message+"\n", "error")
				for _, message := range publishFailedResponse.Payload {
					ctx.Printer.PrintMessage("\t"+message+"\n", "error")
				}
			} else if err != nil {
				ctx.Printer.PrintMessage("Publish failed: \n"+err.Error()+"\n", "error")
			} else {
				ctx.Printer.PrintMessage("Published successfully\n", "green")
			}

			return err
		},
	}

	return publishCommand
}

type MessagesContext struct {
	CliVersion  string
	LocalConfig LocalConfig
	Messager    Messager
	Printer     Printer
	CliClient   *cliClient.CliClient
}

//go:embed regoDefinitionSchema.yaml
var regoDefinitionSchemaStr string

func publish(ctx *PublishCommandContext, path string, localConfigContent *localConfig.LocalConfig) (*cliClient.PublishFailedResponse, error) {
	policiesConfiguration, err := ctx.FilesExtractor.ExtractYamlFileToUnknownStruct(path)
	if err != nil {
		return nil, err
	}

	if customRules, ok := policiesConfiguration["customRules"]; ok {
		valid, err :=
			validateIfRegoRules(ctx, customRules.([]interface{}))
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, fmt.Errorf("failed to validate rego custom rules")
		}
	}

	return ctx.PublishCliClient.PublishPolicies(policiesConfiguration, localConfigContent.Token)
}

func validateIfRegoRules(ctx *PublishCommandContext, customRules []interface{}) (bool, error) {
	for _, regoRule := range customRules {
		b, _ := json.Marshal(regoRule)
		regoRuleStr := string(b)

		isRegoRule := strings.Contains(regoRuleStr, jsonSchemaValidator.RegoDefinitionCustomKey)
		if isRegoRule {
			errorsResult, err := ctx.JSONSchemaValidator.ValidateYamlSchema(regoDefinitionSchemaStr, regoRuleStr)

			if err != nil {
				return false, err
			}
			if len(errorsResult) > 0 {
				ctx.Printer.PrintYamlSchemaResults(errorsResult, err)
				return false, nil
			}
		}
	}
	return true, nil
}
