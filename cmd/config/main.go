package config

import (
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

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
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
	}

	configCommand.AddCommand(NewSetCommand(ctx))

	return configCommand
}
