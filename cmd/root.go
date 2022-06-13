package cmd

import (
	"time"

	"github.com/datreeio/datree/pkg/evaluation"

	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/cmd/completion"
	"github.com/datreeio/datree/cmd/config"
	"github.com/datreeio/datree/cmd/kustomize"
	"github.com/datreeio/datree/cmd/publish"
	schemaValidator "github.com/datreeio/datree/cmd/schema-validator"
	"github.com/datreeio/datree/cmd/test"
	"github.com/datreeio/datree/cmd/version"
	"github.com/datreeio/datree/pkg/ciContext"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/executor"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "datree",
	Short: "Datree is a static code analysis tool for kubernetes files",
	Long:  `Datree is a static code analysis tool for kubernetes files. Full code can be found at https://github.com/datreeio/datree`,
}

var CliVersion string

func NewRootCommand(app *App) *cobra.Command {
	startTime := time.Now()

	rootCmd.AddCommand(test.New(&test.TestCommandContext{
		CliVersion:     CliVersion,
		Evaluator:      app.Context.Evaluator,
		LocalConfig:    app.Context.LocalConfig,
		Messager:       app.Context.Messager,
		Printer:        app.Context.Printer,
		Reader:         app.Context.Reader,
		K8sValidator:   app.Context.K8sValidator,
		CliClient:      app.Context.CliClient,
		FilesExtractor: app.Context.FilesExtractor,
		CiContext:      app.Context.CiContext,
		StartTime:      startTime,
	}))

	rootCmd.AddCommand(kustomize.New(&test.TestCommandContext{
		CliVersion:     CliVersion,
		Evaluator:      app.Context.Evaluator,
		LocalConfig:    app.Context.LocalConfig,
		Messager:       app.Context.Messager,
		Printer:        app.Context.Printer,
		Reader:         app.Context.Reader,
		K8sValidator:   app.Context.K8sValidator,
		CliClient:      app.Context.CliClient,
		FilesExtractor: app.Context.FilesExtractor,
		StartTime:      startTime,
	}, &kustomize.KustomizeContext{CommandRunner: app.Context.CommandRunner}))

	rootCmd.AddCommand(version.New(&version.VersionCommandContext{
		CliVersion: CliVersion,
		Messager:   app.Context.Messager,
		Printer:    app.Context.Printer,
	}))

	rootCmd.AddCommand(config.New(&config.ConfigCommandContext{
		CliVersion:  CliVersion,
		Messager:    app.Context.Messager,
		Printer:     app.Context.Printer,
		LocalConfig: app.Context.LocalConfig,
	}))

	rootCmd.AddCommand(publish.New(&publish.PublishCommandContext{
		CliVersion:       CliVersion,
		LocalConfig:      app.Context.LocalConfig,
		Messager:         app.Context.Messager,
		Printer:          app.Context.Printer,
		PublishCliClient: app.Context.CliClient,
		FilesExtractor:   app.Context.FilesExtractor,
	}))

	rootCmd.AddCommand(completion.New())

	rootCmd.AddCommand(schemaValidator.New(&schemaValidator.JSONSchemaValidatorCommandContext{
		JSONSchemaValidator: app.Context.JSONSchemaValidator,
		Printer:             app.Context.Printer,
	}))

	return rootCmd
}

type Context struct {
	LocalConfig         *localConfig.LocalConfigClient
	Evaluator           *evaluation.Evaluator
	CiContext           *ciContext.CIContext
	CliClient           *cliClient.CliClient
	Messager            *messager.Messager
	Printer             *printer.Printer
	Reader              *fileReader.FileReader
	K8sValidator        *validation.K8sValidator
	JSONSchemaValidator *jsonSchemaValidator.JSONSchemaValidator
	CommandRunner       *executor.CommandRunner
	FilesExtractor      *files.FilesExtractor
}

type App struct {
	Context *Context
}
