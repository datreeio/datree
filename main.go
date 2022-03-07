package main

import (
	"errors"
	"os"

	"github.com/datreeio/datree/bl/errorReporter"
	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/cmd/test"
)

const DEFAULT_ERR_EXIT_CODE = 1
const VIOLATIONS_FOUND_EXIT_CODE = 2

func main() {
	// global error handling
	defer func() {
		if panicErr := recover(); panicErr != nil {
			errorReporter.ReportCliError(panicErr)
		}
	}()

	if err := cmd.Execute(); err != nil {
		if errors.Is(err, test.ViolationsFoundError) {
			os.Exit(VIOLATIONS_FOUND_EXIT_CODE)
		}
		os.Exit(DEFAULT_ERR_EXIT_CODE)
	}
}
