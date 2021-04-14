package cmd

import (
	"os"

	"github.com/datreeio/datree/bl"
	"github.com/datreeio/datree/cmd/test"
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

func init() {
	app := startup()

	rootCmd.AddCommand(test.CreateTestCommand(&test.TestCommandContext{
		Evaluator:   app.Context.Evaluator,
		LocalConfig: app.Context.LocalConfig,
	}))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type app struct {
	Context struct {
		LocalConfig *localConfig.LocalConfiguration
		Evaluator   *bl.Evaluator
	}
}

func startup() *app {
	app := &app{
		Context: struct {
			LocalConfig *localConfig.LocalConfiguration
			Evaluator   *bl.Evaluator
		}{},
	}

	client := cliClient.NewCliClient(deploymentConfig.URL)
	extractor := propertiesExtractor.NewPropertiesExtractor(nil)
	printer := printer.CreateNewPrinter()
	evaluator := bl.CreateNewEvaluator(extractor, client, printer)

	app.Context.LocalConfig = &localConfig.LocalConfiguration{}
	app.Context.Evaluator = evaluator

	return app
}
