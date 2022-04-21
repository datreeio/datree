package evaluation

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"

	"github.com/datreeio/datree/pkg/printer"
	"github.com/stretchr/testify/assert"
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
	results           FormattedResults
	invalidYamlFiles  []*extractor.InvalidFile
	invalidK8sFiles   []*extractor.InvalidFile
	evaluationSummary printer.EvaluationSummary
	loginURL          string
	outputFormat      string
}

type printResultsTestCase struct {
	name     string
	args     *printResultsTestCaseArgs
	expected error
}

type expectedOutputs struct {
	json  string
	xml   string
	yaml  string
	JUnit string
}

// TODO: fill missing call assertions
func TestPrintResults(t *testing.T) {
	tests := []*printResultsTestCase{
		print_resultst(""),
		print_resultst("json"),
		print_resultst("yaml"),
		print_resultst("xml"),
		print_resultst("JUnit"),
	}
	for _, tt := range tests {
		mockedPrinter := &mockPrinter{}
		mockedPrinter.On("PrintWarnings", mock.Anything, mock.Anything, mock.Anything)
		mockedPrinter.On("PrintSummaryTable", mock.Anything)
		mockedPrinter.On("PrintEvaluationSummary", mock.Anything, mock.Anything)

		t.Run(tt.name, func(t *testing.T) {
			_ = PrintResults(&PrintResultsData{tt.args.results, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, tt.args.evaluationSummary, tt.args.loginURL, tt.args.outputFormat, mockedPrinter, "1.18.0", false, "Default", validation.K8sValidationWarningPerValidFile{}})

			if tt.args.outputFormat == "json" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
			} else if tt.args.outputFormat == "yaml" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
			} else if tt.args.outputFormat == "xml" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
			} else if tt.args.outputFormat == "JUnit" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
			} else {
				pwd, _ := os.Getwd()
				warnings, _ := parseToPrinterWarnings(tt.args.results.EvaluationResults, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, pwd, "1.18.0", validation.K8sValidationWarningPerValidFile{}, false)
				mockedPrinter.AssertCalled(t, "PrintWarnings", warnings)
			}
		})
	}
}

func TestCustomOutputs(t *testing.T) {
	formattedOutput := createFormattedOutput()
	expectedOutputs := getExpectedOutputs()

	jsonStdout := readOutput("json", formattedOutput)
	assert.Equal(t, expectedOutputs.json, jsonStdout)

	yamlStdout := readOutput("yaml", formattedOutput)
	assert.Equal(t, expectedOutputs.yaml, yamlStdout)

	xmlStdout := readOutput("xml", formattedOutput)
	assert.Equal(t, expectedOutputs.xml, xmlStdout)

	JUnitStdout := readOutput("JUnit", formattedOutput)
	assert.Equal(t, expectedOutputs.JUnit, JUnitStdout)
}

func readOutput(outputFormat string, formattedOutput FormattedOutput) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		out <- buf.String()
	}()

	switch {
	case outputFormat == "json":
		jsonOutput(&formattedOutput)
	case outputFormat == "yaml":
		yamlOutput(&formattedOutput)
	case outputFormat == "xml":
		xmlOutput(&formattedOutput)
	case outputFormat == "JUnit":
		err := jUnitOutput(&formattedOutput)
		if err != nil {
			panic("unexpected error in printer_test: " + err.Error())
		}
	}

	writer.Close()
	return <-out
}

func createFormattedOutput() FormattedOutput {
	evaluationResults := &NonInteractiveEvaluationResults{
		PolicySummary: &PolicySummary{
			PolicyName:         "Default",
			TotalRulesInPolicy: 21,
			TotalRulesFailed:   4,
			TotalPassedCount:   0,
		},
		FormattedEvaluationResults: []*FormattedEvaluationResults{
			{
				FileName: "File1",
				RuleResults: []*RuleResult{
					{
						Identifier:       "CONTAINERS_MISSING_IMAGE_VALUE_VERSION",
						Name:             "Ensure each container image has a pinned (tag) version",
						MessageOnFailure: "Incorrect value for key `image` - specify an image version to avoid unpleasant \"version surprises\" in the future",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
							Occurrences:  1,
						}},
					},
					{
						Identifier:       "CONTAINERS_MISSING_MEMORY_LIMIT_KEY",
						Name:             "Ensure each container has a configured memory limit",
						MessageOnFailure: "Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
							Occurrences:  1,
						}},
					},
					{
						Identifier:       "WORKLOAD_INVALID_LABELS_VALUE",
						Name:             "Ensure workload has valid label values",
						MessageOnFailure: "Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
							Occurrences:  1,
						}},
					},
					{
						Identifier:       "CONTAINERS_MISSING_LIVENESSPROBE_KEY",
						Name:             "Ensure each container has a configured liveness probe",
						MessageOnFailure: "Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
							Occurrences:  1,
						}},
					},
				},
			},
		},
	}

	return FormattedOutput{
		PolicyValidationResults: evaluationResults.FormattedEvaluationResults,
		PolicySummary:           evaluationResults.PolicySummary,
		EvaluationSummary: NonInteractiveEvaluationSummary{
			ConfigsCount:                1,
			FilesCount:                  1,
			PassedYamlValidationCount:   1,
			K8sValidation:               "1/1",
			PassedPolicyValidationCount: 0,
		},
	}
}

func getExpectedOutputs() expectedOutputs {
	jsonOutput, _ := os.ReadFile("./printer_test_expected_outputs/json_output.json")
	yamlOutput, _ := os.ReadFile("./printer_test_expected_outputs/yaml_output.yaml")
	xmlOutput, _ := os.ReadFile("./printer_test_expected_outputs/xml_output.xml")
	jUnitOutput, _ := os.ReadFile("./printer_test_expected_outputs/JUnit_output.xml")
	return expectedOutputs{
		json:  string(jsonOutput),
		yaml:  string(yamlOutput),
		xml:   string(xmlOutput),
		JUnit: string(jUnitOutput),
	}
}

func print_resultst(outputFormat string) *printResultsTestCase {
	return &printResultsTestCase{
		name: "Print Results Text",
		args: &printResultsTestCaseArgs{
			results: FormattedResults{
				EvaluationResults: &EvaluationResults{
					FileNameRuleMapper: map[string]map[string]*Rule{},
					Summary: EvaluationResultsSummery{
						TotalFailedRules:  0,
						TotalSkippedRules: 0,
						FilesCount:        0,
						FilesPassedCount:  0,
					},
				},
				NonInteractiveEvaluationResults: &NonInteractiveEvaluationResults{
					PolicySummary: &PolicySummary{
						PolicyName:         "Default",
						TotalRulesInPolicy: 0,
						TotalRulesFailed:   0,
						TotalPassedCount:   0,
					},
					FormattedEvaluationResults: []*FormattedEvaluationResults{},
				},
			},
			invalidYamlFiles:  []*extractor.InvalidFile{},
			invalidK8sFiles:   []*extractor.InvalidFile{},
			evaluationSummary: printer.EvaluationSummary{},
			loginURL:          "login/url",
			outputFormat:      outputFormat,
		},
		expected: nil,
	}
}
