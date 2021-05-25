package evaluation

import (
	"os"
	"testing"

	"github.com/datreeio/datree/bl/validation"

	"github.com/datreeio/datree/pkg/printer"
	"github.com/stretchr/testify/mock"
)

type mockPrinter struct {
	mock.Mock
}

func (m *mockPrinter) PrintWarnings(warnings []printer.Warning) {
	m.Called(warnings)
}

func (c *mockPrinter) PrintSummaryTable(summary printer.Summary) {
	c.Called(summary)
}

type printResultsTestCaseArgs struct {
	results           *EvaluationResults
	invalidFiles      []*validation.InvalidFile
	evaluationSummary printer.EvaluationSummary
	loginURL          string
	outputFormat      string
	printer           Printer
}

type printResultsTestCase struct {
	name     string
	args     *printResultsTestCaseArgs
	expected error
}

// TODO: write actual tests, what is this even?
func TestPrintResults(t *testing.T) {
	tests := []*printResultsTestCase{}
	for _, tt := range tests {
		mockedPrinter := &mockPrinter{}
		mockedPrinter.On("PrintWarnings", mock.Anything)
		mockedPrinter.On("PrintSummaryTable", mock.Anything)
		t.Run(tt.name, func(t *testing.T) {
			PrintResults(tt.args.results, tt.args.invalidFiles, tt.args.evaluationSummary, tt.args.loginURL, tt.args.outputFormat, tt.args.printer, "1.18.0")

			if tt.args.outputFormat == "json" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
				mockedPrinter.AssertCalled(t, "PrintSummaryTable")
			} else if tt.args.outputFormat == "yaml" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
				mockedPrinter.AssertCalled(t, "PrintSummaryTable")

			} else {
				pwd, _ := os.Getwd()
				warnings, _ := parseToPrinterWarnings(tt.args.results, tt.args.invalidFiles, pwd, "1.18.0")
				mockedPrinter.AssertCalled(t, "PrintWarnings", warnings)
				mockedPrinter.AssertCalled(t, "PrintSummaryTable", tt.args.results, tt.args.loginURL)
			}
		})
	}
}
