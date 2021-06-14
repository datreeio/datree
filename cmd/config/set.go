package config

import (
	"fmt"

	"github.com/datreeio/datree/bl/messager"
	"github.com/spf13/cobra"
)

func NewSetCommand(ctx *ConfigCommandContext) *cobra.Command {
	setCommand := &cobra.Command{
		Use:   "set",
		Short: "Set configuration value",
		Long:  `Apply value for specific key in datree config.yaml file. Defaults to $HOME/.datree/config.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			messages := make(chan *messager.VersionMessage, 1)
			go ctx.Messager.LoadVersionMessages(messages, ctx.CliVersion)

			err := set(ctx, args[0], args[1])
			if err != nil {
				fmt.Printf("Failed setting %s with value %s. Error: %s", args[0], args[1], err)
			}

			msg, ok := <-messages
			if ok {
				ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
			}
		},
		Args: cobra.ExactArgs(2),
	}

	return setCommand
}

func set(ctx *ConfigCommandContext, key string, value string) error {
	_, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return err
	}

	err = ctx.LocalConfig.Set(key, value)
	return err
}
