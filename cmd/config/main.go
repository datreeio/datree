package config

import (
	"fmt"
	"reflect"

	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/spf13/cobra"
)

type Messager interface {
	LoadVersionMessages(cliVersion string) chan *messager.VersionMessage
}

type Printer interface {
	PrintMessage(messageText string, messageColor string)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.LocalConfig, error)
	Set(key string, value string) error
	Get(key string) string
}

type ConfigCommandContext struct {
	Messager    Messager
	CliVersion  string
	Printer     Printer
	LocalConfig LocalConfig
}

func New(ctx *ConfigCommandContext) *cobra.Command {
	configCommand := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  `Internal configuration management for datree config file`,
		Example: utils.Example(`
		# Change the token in the datree config.yaml file
		datree config set token <MY_EXAMPLE_TOKEN>

		# Get the token from datree config.yaml file
		datree config get token
		`),
	}

	configCommand.AddCommand(NewSetCommand(ctx))
	configCommand.AddCommand(NewGetCommand(ctx))

	return configCommand
}

func validateKey(key string) error {
	validKeys := make(map[string]bool)
	validKeys["token"] = true

	if val, ok := validKeys[key]; !ok || !val {
		return fmt.Errorf("key must be one of: %s", reflect.ValueOf(validKeys).MapKeys())
	}
	return nil
}
