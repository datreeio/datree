package config

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func NewGetCommand(ctx *ConfigCommandContext) *cobra.Command {
	getCommand := &cobra.Command{
		Use:   "get",
		Short: "Get configuration value",
		Long:  `Get value for specific key from datree config.yaml file. Defaults to $HOME/.datree/config.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			messages := ctx.Messager.LoadVersionMessages(ctx.CliVersion)

			err := get(ctx, args[0])
			if err != nil {
				fmt.Printf("Failed getting %s . Error: %s", args[0], err)
			}

			msg, ok := <-messages
			if ok {
				ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires exactly 1 argument")
			}

			return validateKey(args[0])
		},
	}

	return getCommand
}

func get(ctx *ConfigCommandContext, key string) error {
	_, err := ctx.LocalConfig.GetLocalConfiguration()
	if err != nil {
		return err
	}

	val := ctx.LocalConfig.Get(key)
	fmt.Println(val)
	return nil
}
