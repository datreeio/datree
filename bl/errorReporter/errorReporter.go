package errorReporter

import (
	"fmt"
	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/localConfig"
	"runtime/debug"
)

func ReportCliError(panicErr interface{}) {
	reporter := NewErrorReporter(cliClient.NewCliClient(deploymentConfig.URL), localConfig.NewLocalConfig())
	reporter.ReportCliError(panicErr)
}

type ErrorReporter struct {
	client *cliClient.CliClient
	config *localConfig.LocalConfig
}

func NewErrorReporter(client *cliClient.CliClient, localConfig *localConfig.LocalConfig) *ErrorReporter {
	return &ErrorReporter{
		client: client,
		config: localConfig,
	}
}

func (reporter *ErrorReporter) ReportCliError(panicErr interface{}) {
	status, err := reporter.client.ReportCliError(cliClient.ReportCliErrorRequest{
		CliId:        reporter.getCliId(),
		CliVersion:   cmd.CliVersion,
		ErrorMessage: parsePanicError(panicErr),
		StackTrace:   string(debug.Stack()),
	})
	if err != nil {
		fmt.Println("failed to report unexpected panic:\nstatus code:", status, "\nerror:", err.Error())
	} else {
		fmt.Println("panic reported")
	}
}

func (reporter *ErrorReporter) getCliId() string {
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
