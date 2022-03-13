package errorReporter

import (
	"fmt"
	"runtime/debug"

	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/localConfig"
)


type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.LocalConfig, error)
}

type CliClient interface {
	ReportCliError(reportCliErrorRequest cliClient.ReportCliErrorRequest, uri string) (StatusCode int, Error error)
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

func (reporter *ErrorReporter) ReportCliPanicError(panicErr interface{}) {
	reporter.ReportCliError(panicErr, "/report-cli-panic-error")
}

func (reporter *ErrorReporter) ReportCliUnexpectedError(unexpectedError error) {
	reporter.ReportCliError(unexpectedError, "/report-cli-unexpected-error")
}

func (reporter *ErrorReporter) ReportCliError(panicErr interface{}, uri string) {
	errorMessage := parsePanicError(panicErr)
	cliId := reporter.getCliId()
	_, err := reporter.client.ReportCliError(cliClient.ReportCliErrorRequest{
		ClientId:     cliId,
		Token:        cliId,
		CliVersion:   cmd.CliVersion,
		ErrorMessage: errorMessage,
		StackTrace:   string(debug.Stack()),
	}, uri)
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
		return config.Token
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
