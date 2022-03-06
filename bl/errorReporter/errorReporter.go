package errorReporter

import (
	"fmt"
	"runtime/debug"

	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/localConfig"
)

func ReportCliError(panicErr interface{}) {
	reporter := NewErrorReporter(cliClient.NewCliClient(deploymentConfig.URL), localConfig.NewLocalConfig())
	reporter.ReportCliError(panicErr)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
}

type CliClient interface {
	ReportCliError(reportCliErrorRequest cliClient.ReportCliErrorRequest) (StatusCode int, Error error)
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

func (reporter *ErrorReporter) ReportCliError(panicErr interface{}) {
	status, err := reporter.client.ReportCliError(cliClient.ReportCliErrorRequest{
		Token:        reporter.getCliId(),
		CliVersion:   cmd.CliVersion,
		ErrorMessage: parsePanicError(panicErr),
		StackTrace:   string(debug.Stack()),
	})
	if err != nil {
		fmt.Println("failed to report unexpected error:\nstatus code:", status, "\nerror:", err.Error())
	} else {
		fmt.Println("unexpected error reported")
	}
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
