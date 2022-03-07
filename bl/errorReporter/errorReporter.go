package errorReporter

import (
	"fmt"
	"runtime/debug"

	"github.com/datreeio/datree/pkg/printer"

	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/localConfig"
)

func ReportCliError(panicErr interface{}) {
	reporter := NewErrorReporter(cliClient.NewCliClient(deploymentConfig.URL), localConfig.NewLocalConfig(), printer.CreateNewPrinter())
	reporter.ReportCliError(panicErr)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
}

type CliClient interface {
	ReportCliError(reportCliErrorRequest cliClient.ReportCliErrorRequest) (StatusCode int, Error error)
}

type Printer interface {
	PrintMessage(messageText string, messageColor string)
}

type ErrorReporter struct {
	config  LocalConfig
	client  CliClient
	printer Printer
}

func NewErrorReporter(client CliClient, localConfig LocalConfig, printer Printer) *ErrorReporter {
	return &ErrorReporter{
		client:  client,
		config:  localConfig,
		printer: printer,
	}
}

func (reporter *ErrorReporter) ReportCliError(panicErr interface{}) {
	errorMessage := parsePanicError(panicErr)
	cliId := reporter.getCliId()
	_, err := reporter.client.ReportCliError(cliClient.ReportCliErrorRequest{
		ClientId:     cliId,
		Token:        cliId,
		CliVersion:   cmd.CliVersion,
		ErrorMessage: errorMessage,
		StackTrace:   string(debug.Stack()),
	})
	if err != nil {
		// do nothing
	}
	reporter.printer.PrintMessage(fmt.Sprintf("Unexpected error: %s\n", errorMessage), "error")
}

func (reporter *ErrorReporter) getCliId() (cliId string) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			cliId = "unknown"
		}
	}()

	config, err := reporter.config.GetLocalConfiguration()
	if err != nil {
		return "unknown"
	} else {
		return config.CliId
	}
}

func parsePanicError(panicErr interface{}) string {
	switch panicErr := panicErr.(type) {
	case string:
		return panicErr
	case error:
		return panicErr.Error()
	default:
		return fmt.Sprintf("%v", panicErr)
	}
}
