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
	PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) (*cliClient.PublishFailedResponse, error)
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
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			messages := make(chan *messager.VersionMessage, 1)
			go ctx.Messager.LoadVersionMessages(messages, ctx.CliVersion)
			defer func() {
				msg, ok := <-messages
				if ok {
					ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
				}
			}()

			publishFailedResponse, err := publish(ctx, args[0])
			if publishFailedResponse != nil {
				ctx.Printer.PrintMessage("Publish failed:\n", "error")
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

func publish(ctx *PublishCommandContext, path string) (*cliClient.PublishFailedResponse, error) {
	localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return nil, err
	}

	policiesConfiguration, err := files.ExtractYamlFileToUnknownStruct(path)
	if err != nil {
		return nil, err
	}

	return ctx.PublishCliClient.PublishPolicies(policiesConfiguration, localConfigContent.CliId)
}
