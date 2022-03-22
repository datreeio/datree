package main

import (
	"errors"
	"fmt"
	"github.com/datreeio/datree/pkg/networkValidator"
	"os"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/deploymentConfig"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/datreeio/datree/pkg/utils"

	"github.com/datreeio/datree/bl/errorReporter"
	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/cmd/test"
)

const DEFAULT_ERR_EXIT_CODE = 1
const VIOLATIONS_FOUND_EXIT_CODE = 2

func main() {
	validator := networkValidator.NewNetworkValidator()
	cliClient := cliClient.NewCliClient(deploymentConfig.URL, validator)
	localConfig := localConfig.NewLocalConfigClient(cliClient, validator)

	reporter := errorReporter.NewErrorReporter(cliClient, localConfig)
	globalPrinter := printer.CreateNewPrinter()

	defer func() {
		if panicErr := recover(); panicErr != nil {
			reporter.ReportPanicError(panicErr)

			globalPrinter.PrintMessage(fmt.Sprintf("Unexpected error: %s\n", utils.ParseErrorToString(panicErr)), "error")
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
