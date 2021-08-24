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

type CliClient interface {
	PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) error
}

type PublishCommandContext struct {
	CliVersion       string
	LocalConfig      LocalConfig
	Messager         Messager
	Printer          Printer
	PublishCliClient CliClient
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
			messages := make(chan *messager.VersionMessage, 1)
			go ctx.Messager.LoadVersionMessages(messages, ctx.CliVersion)
			defer func() {
				msg, ok := <-messages
				if ok {
					ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
				}
			}()

			err := publish(ctx, args[0])
			if err != nil {
				ctx.Printer.PrintMessage("Publish failed:\n"+err.Error()+"\n", "error")
			} else {
				ctx.Printer.PrintMessage("Published successfully\n", "green")
			}

			return err
		},
		SilenceUsage:  true,
		SilenceErrors: true,
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

func publish(ctx *PublishCommandContext, path string) error {
	localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return err
	}

	policiesConfiguration, err := files.ExtractYamlFileToUnknownStruct(path)
	if err != nil {
		return err
	}

	err = ctx.PublishCliClient.PublishPolicies(policiesConfiguration, localConfigContent.CliId)
	return err
}
