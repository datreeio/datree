package cmd

import (
	"github.com/datreeio/datree/bl"
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

	rootCmd.AddCommand(test.NewTestCommand(&test.TestCommandContext{
		CliVersion:  CliVersion,
		Evaluator:   app.Context.Evaluator,
		LocalConfig: app.Context.LocalConfig,
		CliClient:   app.Context.CliClient,
	}))

	rootCmd.AddCommand(version.NewVersionCommand(&version.VersionCommandContext{
		CliVersion: CliVersion,
		CliClient:  app.Context.CliClient,
	}))
}

func Execute() error {
	return rootCmd.Execute()
}

type app struct {
	Context struct {
		LocalConfig *localConfig.LocalConfiguration
		Evaluator   *bl.Evaluator
		CliClient   *cliClient.CliClient
	}
}

func startup() *app {
	app := &app{
		Context: struct {
			LocalConfig *localConfig.LocalConfiguration
			Evaluator   *bl.Evaluator
			CliClient   *cliClient.CliClient
		}{},
	}

	cliClient := cliClient.NewCliClient(deploymentConfig.URL)
	extractor := propertiesExtractor.NewPropertiesExtractor(nil)
	printer := printer.CreateNewPrinter()
	evaluator := bl.CreateNewEvaluator(extractor, cliClient, printer)

	app.Context.CliClient = cliClient
	app.Context.LocalConfig = &localConfig.LocalConfiguration{}
	app.Context.Evaluator = evaluator

	return app
}
