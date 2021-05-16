package evaluation

import (
	"github.com/datreeio/datree/bl/output"
	"github.com/datreeio/datree/pkg/printer"
)

type Printer interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
}

// url := "https://app.datree.io/login?cliId=" + cliId
func PrintResults(results *EvaluationResults, loginURL string, outputFormat string, printer Printer) error {
	switch {
	case outputFormat == "json":
		return output.JSONOutput(results)
	case outputFormat == "yaml":
		return output.YAMLOutput(results)
	default:
		return output.TextOutput(results, loginURL, printer)
	}
}
