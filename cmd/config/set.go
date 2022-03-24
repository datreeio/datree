package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

func NewSetCommand(ctx *ConfigCommandContext) *cobra.Command {
	setCommand := &cobra.Command{
		Use:   "set",
		Short: "Set configuration value",
		Long:  `Apply value for specific key in datree config.yaml file. Defaults to $HOME/.datree/config.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			messages := ctx.Messager.LoadVersionMessages(ctx.CliVersion)

			err := set(ctx, args[0], args[1])
			if err != nil {
				fmt.Printf("Failed setting %s with value %s. Error: %s", args[0], args[1], err)
			}

			msg, ok := <-messages
			if ok {
				ctx.Printer.PrintMessage(msg.MessageText+"\n", msg.MessageColor)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("requires exactly 2 arguments")
			}

			validKeys := make(map[string]bool)
			validKeys["token"] = true
			validKeys["offline"] = true

			if val, ok := validKeys[args[0]]; !ok || !val {
				return fmt.Errorf("key must be one of: %s", reflect.ValueOf(validKeys).MapKeys())
			}

			return nil
		},
	}

	return setCommand
}

func set(ctx *ConfigCommandContext, key string, value string) error {
	err := ctx.LocalConfig.Set(key, value)
	return err
}
