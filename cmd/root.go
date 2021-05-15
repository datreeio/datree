package cmd

import (
	"github.com/datreeio/datree/bl/evaluator"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validator"
	"github.com/datreeio/datree/cmd/test"
	"github.com/datreeio/datree/cmd/version"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/datreeio/datree/pkg/propertiesExtractor"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "datree",
	Short: "Datree is a static code analysis tool for kubernetes files",
	Long:  `Datree is a static code analysis tool for kubernetes files. Full code can be found at https://github.com/datreeio/datree`,
}

var CliVersion string

func init() {
	app := startup()

	rootCmd.AddCommand(test.NewCommand(&test.TestCommandContext{
		CliVersion:  CliVersion,
		Evaluator:   app.context.Evaluator,
		LocalConfig: app.context.LocalConfig,
		Messager:    app.context.Messager,
	}))

	rootCmd.AddCommand(version.NewCommand(&version.VersionCommandContext{
		CliVersion: CliVersion,
		Messager:   app.context.Messager,
	}))
}

func Execute() error {
	return rootCmd.Execute()
}

type context struct {
	LocalConfig *localConfig.LocalConfiguration
	Evaluator   *evaluator.Evaluator
	CliClient   *cliClient.CliClient
	Messager    *messager.Messager
}
type app struct {
	context *context
}

func startup() *app {
	cliClient := cliClient.NewCliClient(deploymentConfig.URL)
	extractor := propertiesExtractor.NewPropertiesExtractor(nil)
	printer := printer.CreateNewPrinter()
	validator := validator.New("1.18.0")

	return &app{
		context: &context{
			LocalConfig: &localConfig.LocalConfiguration{},
			Evaluator:   evaluator.New(extractor, cliClient, printer, validator),
			CliClient:   cliClient,
			Messager:    messager.New(cliClient, printer),
		},
	}
}
