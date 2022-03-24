package errorReporter

import (
	"runtime/debug"

	"github.com/datreeio/datree/pkg/utils"

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

type ErrorReporter struct {
	config LocalConfig
	client CliClient
}

func NewErrorReporter(client CliClient, localConfig LocalConfig) *ErrorReporter {
	return &ErrorReporter{
		client: client,
		config: localConfig,
	}
}

func (reporter *ErrorReporter) ReportPanicError(panicErr interface{}) {
	reporter.ReportError(panicErr, "/report-cli-panic-error")
}

func (reporter *ErrorReporter) ReportUnexpectedError(unexpectedError error) {
	reporter.ReportError(unexpectedError, "/report-cli-unexpected-error")
}

func (reporter *ErrorReporter) ReportError(error interface{}, uri string) {
	errorMessage := utils.ParseErrorToString(error)
	localConfig := reporter.getLocalConfig()
	_, _ = reporter.client.ReportCliError(cliClient.ReportCliErrorRequest{
		ClientId:     localConfig.ClientId,
		Token:        localConfig.Token,
		CliVersion:   cmd.CliVersion,
		ErrorMessage: errorMessage,
		StackTrace:   string(debug.Stack()),
	}, uri)

}

func (reporter *ErrorReporter) getLocalConfig() (unknownLocalConfig *localConfig.LocalConfig) {
	unknownLocalConfig = &localConfig.LocalConfig{ClientId: "unknown", Token: "unknown"}
	defer func() {
		_ = recover()

	}()

	config, err := reporter.config.GetLocalConfiguration()
	if err != nil {
		return unknownLocalConfig
	} else {
		return config
	}

}
