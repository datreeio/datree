package publish

import (
	"fmt"
	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/pkg/cliClient"
	"strings"

	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/spf13/cobra"
)

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type Printer interface {
	PrintMessage(messageText string, messageColor string)
}

type Reader interface {
	FilterFiles(paths []string) ([]string, error)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
}

type PublishCommandContext struct {
	CliVersion  string
	LocalConfig LocalConfig
	Messager    Messager
	Printer     Printer
	Reader      Reader
	CliClient   *cliClient.CliClient
}

func New(ctx *PublishCommandContext) *cobra.Command {
	publishCommand := &cobra.Command{
		Use:   "publish <fileName>",
		Short: "Publish policies configuration for given <fileName>.",
		Long:  "Publish policies configuration for given <fileName>. Input should be the path to the Policy as Code yaml configuration file ",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				errMessage := "Requires 1 arg\n"
				cmd.Usage()
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error = nil
			defer func() {
				if err != nil {
					ctx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()

			err = publish(ctx, args[0])
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

	publishResponse, err := ctx.CliClient.PublishPolicies(policiesConfiguration, localConfigContent.CliId)
	if err != nil {
		return err
	}

	if publishResponse.IsSuccessful {
		ctx.Printer.PrintMessage("Publish Successful", "green")
	} else {
		printRows := append([]string{"Publish Failed"}, publishResponse.Errors...)
		errorMessage := strings.Join(printRows, "\n") + "\n"
		ctx.Printer.PrintMessage(errorMessage, "error")
	}
	return nil
}
