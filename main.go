package main

import (
	"errors"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"os"

	"github.com/datreeio/datree/bl/errorReporter"
	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/cmd/test"
)

const DEFAULT_ERR_EXIT_CODE = 1
const VIOLATIONS_FOUND_EXIT_CODE = 2

func main() {

	reporter := errorReporter.NewErrorReporter(cliClient.NewCliClient(deploymentConfig.URL), localConfig.NewLocalConfig(), printer.CreateNewPrinter())

	defer func() {
		if panicErr := recover(); panicErr != nil {
			reporter.ReportPanicError(panicErr)
			os.Exit(DEFAULT_ERR_EXIT_CODE)
		}
	}()

	if err := cmd.Execute(); err != nil {
		if errors.Is(err, test.ViolationsFoundError) {
			os.Exit(VIOLATIONS_FOUND_EXIT_CODE)
		}
		reporter.ReportUnexpectedError(err)
		os.Exit(DEFAULT_ERR_EXIT_CODE)
	}
}
