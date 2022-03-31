package cmd

import (
	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/cmd/completion"
	"github.com/datreeio/datree/cmd/config"
	"github.com/datreeio/datree/cmd/kustomize"
	"github.com/datreeio/datree/cmd/publish"
	schema_validator "github.com/datreeio/datree/cmd/schema-validator"
	"github.com/datreeio/datree/cmd/test"
	"github.com/datreeio/datree/cmd/version"
	"github.com/datreeio/datree/internal/deploymentConfig"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/executor"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/networkValidator"
	"github.com/datreeio/datree/pkg/printer"
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

	rootCmd.AddCommand(test.New(&test.TestCommandContext{
		CliVersion:   CliVersion,
		Evaluator:    app.context.Evaluator,
		LocalConfig:  app.context.LocalConfig,
		Messager:     app.context.Messager,
		Printer:      app.context.Printer,
		Reader:       app.context.Reader,
		K8sValidator: app.context.K8sValidator,
		CliClient:    app.context.CliClient,
	}))

	rootCmd.AddCommand(kustomize.New(&test.TestCommandContext{
		CliVersion:   CliVersion,
		Evaluator:    app.context.Evaluator,
		LocalConfig:  app.context.LocalConfig,
		Messager:     app.context.Messager,
		Printer:      app.context.Printer,
		Reader:       app.context.Reader,
		K8sValidator: app.context.K8sValidator,
		CliClient:    app.context.CliClient,
	}, &kustomize.KustomizeContext{CommandRunner: app.context.CommandRunner}))

	rootCmd.AddCommand(version.New(&version.VersionCommandContext{
		CliVersion: CliVersion,
		Messager:   app.context.Messager,
		Printer:    app.context.Printer,
	}))

	rootCmd.AddCommand(config.New(&config.ConfigCommandContext{
		CliVersion:  CliVersion,
		Messager:    app.context.Messager,
		Printer:     app.context.Printer,
		LocalConfig: app.context.LocalConfig,
	}))

	rootCmd.AddCommand(publish.New(&publish.PublishCommandContext{
		CliVersion:       CliVersion,
		LocalConfig:      app.context.LocalConfig,
		Messager:         app.context.Messager,
		Printer:          app.context.Printer,
		PublishCliClient: app.context.CliClient,
	}))

	rootCmd.AddCommand(completion.New())

	rootCmd.AddCommand(schema_validator.New(&schema_validator.JSONSchemaValidatorCommandContext{
		JSONSchemaValidator: app.context.JSONSchemaValidator,
		Printer:             app.context.Printer,
	}))
}

func Execute() error {
	return rootCmd.Execute()
}

type context struct {
	LocalConfig         *localConfig.LocalConfigClient
	Evaluator           *evaluation.Evaluator
	CliClient           *cliClient.CliClient
	Messager            *messager.Messager
	Printer             *printer.Printer
	Reader              *fileReader.FileReader
	K8sValidator        *validation.K8sValidator
	JSONSchemaValidator *jsonSchemaValidator.JSONSchemaValidator
	CommandRunner       *executor.CommandRunner
}

type app struct {
	context *context
}

func startup() *app {
	validator := networkValidator.NewNetworkValidator()
	cliClient := cliClient.NewCliClient(deploymentConfig.URL, validator)
	localConfig := localConfig.NewLocalConfigClient(cliClient, validator)
	printer := printer.CreateNewPrinter()

	return &app{
		context: &context{
			LocalConfig:         localConfig,
			Evaluator:           evaluation.New(cliClient),
			CliClient:           cliClient,
			Messager:            messager.New(cliClient),
			Printer:             printer,
			Reader:              fileReader.CreateFileReader(nil),
			K8sValidator:        validation.New(),
			JSONSchemaValidator: jsonSchemaValidator.New(),
			CommandRunner:       executor.CreateNewCommandRunner(),
		},
	}
}
