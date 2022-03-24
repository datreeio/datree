package config

import (
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
		`),
	}

	configCommand.AddCommand(NewSetCommand(ctx))

	return configCommand
}
