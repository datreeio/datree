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

func (c *mockPrinter) PrintEvaluationSummary(summary printer.EvaluationSummary, k8sVersion string) {
	c.Called(summary, k8sVersion)
}

type printResultsTestCaseArgs struct {
	results           *EvaluationResults
	invalidYamlFiles  []*validation.InvalidYamlFile
	invalidK8sFiles   []*validation.InvalidK8sFile
	evaluationSummary printer.EvaluationSummary
	loginURL          string
	outputFormat      string
}

type printResultsTestCase struct {
	name     string
	args     *printResultsTestCaseArgs
	expected error
}

// TODO: fill missing call assertions
func TestPrintResults(t *testing.T) {
	tests := []*printResultsTestCase{
		print_resultst(""),
		print_resultst("json"),
		print_resultst("yaml"),
	}
	for _, tt := range tests {
		mockedPrinter := &mockPrinter{}
		mockedPrinter.On("PrintWarnings", mock.Anything, mock.Anything, mock.Anything)
		mockedPrinter.On("PrintSummaryTable", mock.Anything)
		mockedPrinter.On("PrintEvaluationSummary", mock.Anything, mock.Anything)

		t.Run(tt.name, func(t *testing.T) {
			_ = PrintResults(tt.args.results, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, tt.args.evaluationSummary, tt.args.loginURL, tt.args.outputFormat, mockedPrinter, "1.18.0", "Default")

			if tt.args.outputFormat == "json" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
			} else if tt.args.outputFormat == "yaml" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")

			} else {
				pwd, _ := os.Getwd()
				warnings, _ := parseToPrinterWarnings(tt.args.results, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, pwd, "1.18.0")
				mockedPrinter.AssertCalled(t, "PrintWarnings", warnings)
			}
		})
	}
}

func print_resultst(outputFormat string) *printResultsTestCase {
	return &printResultsTestCase{
		name: "Print Results Text",
		args: &printResultsTestCaseArgs{
			results: &EvaluationResults{
				FileNameRuleMapper: map[string]map[int]*Rule{},
				Summary: struct {
					TotalFailedRules int
					FilesCount       int
					TotalPassedCount int
				}{
					TotalFailedRules: 0,
					FilesCount:       0,
					TotalPassedCount: 0,
				},
			},
			invalidYamlFiles:  []*validation.InvalidYamlFile{},
			invalidK8sFiles:   []*validation.InvalidK8sFile{},
			evaluationSummary: printer.EvaluationSummary{},
			loginURL:          "login/url",
			outputFormat:      outputFormat,
		},
		expected: nil,
	}
}
