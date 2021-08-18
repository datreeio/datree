package publish

import (
	"fmt"
	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/spf13/cobra"
)

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type Printer interface {
	PrintMessage(messageText string, messageColor string)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
}

type PublishCommandContext struct {
	CliVersion  string
	LocalConfig LocalConfig
	Messager    Messager
	Printer     Printer
	CliClient   *cliClient.CliClient
}

func New(ctx *PublishCommandContext) *cobra.Command {
	publishCommand := &cobra.Command{
		Use:   "publish <fileName>",
		Short: "Publish policies configuration for given <fileName>.",
		Long:  "Publish policies configuration for given <fileName>. Input should be the path to the Policy-as-Code yaml configuration file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				errMessage := "Requires 1 arg\n"
				cmd.Usage()
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := publish(ctx, args[0])
			if err != nil {
				ctx.Printer.PrintMessage("Publish Failed:\n"+err.Error()+"\n", "error")
			} else {
				ctx.Printer.PrintMessage("Publish Successful", "green")
			}

			return err
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	return publishCommand
}

func publish(ctx *PublishCommandContext, path string) error {
	localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return err
	}

	policiesConfiguration, err := files.ExtractYamlFileToUnknownStruct(path)
	if err != nil {
		return err
	}

	err = ctx.CliClient.PublishPolicies(policiesConfiguration, localConfigContent.CliId)
	return err
}
