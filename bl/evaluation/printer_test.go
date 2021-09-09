package evaluation

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/datreeio/datree/bl/validation"

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

type expectedOutputs struct {
	json string
	xml  string
	yaml string
}

// TODO: fill missing call assertions
func TestPrintResults(t *testing.T) {
	tests := []*printResultsTestCase{
		print_resultst(""),
		print_resultst("json"),
		print_resultst("yaml"),
		print_resultst("xml"),
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
			} else if tt.args.outputFormat == "xml" {
				mockedPrinter.AssertNotCalled(t, "PrintWarnings")
			} else {
				pwd, _ := os.Getwd()
				warnings, _ := parseToPrinterWarnings(tt.args.results, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, pwd, "1.18.0")
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
	}

	writer.Close()
	return <-out
}

func createFormattedOutput() FormattedOutput {
	evaluationResults := &EvaluationResults{
		Summary: struct {
			TotalFailedRules int
			FilesCount       int
			TotalPassedCount int
		}{
			TotalFailedRules: 4,
			FilesCount:       1,
			TotalPassedCount: 0,
		},
		FileNameRuleMapper: FileNameRuleMapper{"File1": map[int]*Rule{
			1: &Rule{
				ID:             1,
				Name:           "Ensure each container image has a pinned (tag) version",
				FailSuggestion: "Incorrect value for key `image` - specify an image version to avoid unpleasant \"version surprises\" in the future",
				OccurrencesDetails: []OccurrenceDetails{{
					MetadataName: "rss-site",
					Kind:         "Deployment",
				}},
			},
			4: &Rule{
				ID:             4,
				Name:           "Ensure each container has a configured memory limit",
				FailSuggestion: "Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization",
				OccurrencesDetails: []OccurrenceDetails{{
					MetadataName: "rss-site",
					Kind:         "Deployment",
				}},
			},
			9: &Rule{
				ID:             9,
				Name:           "Ensure workload has valid label values",
				FailSuggestion: "Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it",
				OccurrencesDetails: []OccurrenceDetails{{
					MetadataName: "rss-site",
					Kind:         "Deployment",
				}},
			},
			11: &Rule{
				ID:             11,
				Name:           "Ensure each container has a configured liveness probe",
				FailSuggestion: "Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks",
				OccurrencesDetails: []OccurrenceDetails{{
					MetadataName: "rss-site",
					Kind:         "Deployment",
				}},
			},
		}},
	}

	evaluationSummary := &printer.EvaluationSummary{
		ConfigsCount:              1,
		RulesCount:                21,
		FilesCount:                1,
		PassedYamlValidationCount: 1,
		PassedK8sValidationCount:  1,
		PassedPolicyCheckCount:    0,
	}
	return FormattedOutput{
		EvaluationResults: evaluationResults,
		EvaluationSummary: *evaluationSummary,
	}
}

func getExpectedOutputs() expectedOutputs {
	return expectedOutputs{
		json: "{\"EvaluationResults\":{\"FileNameRuleMapper\":{\"File1\":{\"1\":{\"ID\":1,\"Name\":\"Ensure each container image has a pinned (tag) version\",\"FailSuggestion\":\"Incorrect value for key `image` - specify an image version to avoid unpleasant \\\"version surprises\\\" in the future\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]},\"11\":{\"ID\":11,\"Name\":\"Ensure each container has a configured liveness probe\",\"FailSuggestion\":\"Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]},\"4\":{\"ID\":4,\"Name\":\"Ensure each container has a configured memory limit\",\"FailSuggestion\":\"Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]},\"9\":{\"ID\":9,\"Name\":\"Ensure workload has valid label values\",\"FailSuggestion\":\"Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]}}},\"Summary\":{\"TotalFailedRules\":4,\"FilesCount\":1,\"TotalPassedCount\":0}},\"EvaluationSummary\":{\"ConfigsCount\":1,\"RulesCount\":21,\"FilesCount\":1,\"PassedYamlValidationCount\":1,\"PassedK8sValidationCount\":1,\"PassedPolicyCheckCount\":0},\"InvalidYamlFiles\":null,\"InvalidK8sFiles\":null}\n",
		yaml: "evaluationresults:\n  filenamerulemapper:\n    File1:\n      1:\n        id: 1\n        name: Ensure each container image has a pinned (tag) version\n        failsuggestion: Incorrect value for key `image` - specify an image version\n          to avoid unpleasant \"version surprises\" in the future\n        occurrencesdetails:\n        - metadataname: rss-site\n          kind: Deployment\n      4:\n        id: 4\n        name: Ensure each container has a configured memory limit\n        failsuggestion: Missing property object `limits.memory` - value should be\n          within the accepted boundaries recommended by the organization\n        occurrencesdetails:\n        - metadataname: rss-site\n          kind: Deployment\n      9:\n        id: 9\n        name: Ensure workload has valid label values\n        failsuggestion: Incorrect value for key(s) under `labels` - the vales syntax\n          is not valid so the Kubernetes engine will not accept it\n        occurrencesdetails:\n        - metadataname: rss-site\n          kind: Deployment\n      11:\n        id: 11\n        name: Ensure each container has a configured liveness probe\n        failsuggestion: Missing property object `livenessProbe` - add a properly configured\n          livenessProbe to catch possible deadlocks\n        occurrencesdetails:\n        - metadataname: rss-site\n          kind: Deployment\n  summary:\n    totalfailedrules: 4\n    filescount: 1\n    totalpassedcount: 0\nevaluationsummary:\n  configscount: 1\n  rulescount: 21\n  filescount: 1\n  passedyamlvalidationcount: 1\n  passedk8svalidationcount: 1\n  passedpolicycheckcount: 0\ninvalidyamlfiles: []\ninvalidk8sfiles: []\n\n",
		xml:  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<FormattedOutput>\n\t<EvaluationResults>\n\t\t<FileNameRuleMapper>\n\t\t\t<File filename=\"File1\">\n\t\t\t\t<Rule>\n\t\t\t\t\t<ID>1</ID>\n\t\t\t\t\t<Name>Ensure each container image has a pinned (tag) version</Name>\n\t\t\t\t\t<FailSuggestion>Incorrect value for key `image` - specify an image version to avoid unpleasant &#34;version surprises&#34; in the future</FailSuggestion>\n\t\t\t\t\t<OccurrencesDetails>\n\t\t\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t\t\t</OccurrencesDetails>\n\t\t\t\t</Rule>\n\t\t\t</File>\n\t\t\t<File filename=\"File1\">\n\t\t\t\t<Rule>\n\t\t\t\t\t<ID>4</ID>\n\t\t\t\t\t<Name>Ensure each container has a configured memory limit</Name>\n\t\t\t\t\t<FailSuggestion>Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization</FailSuggestion>\n\t\t\t\t\t<OccurrencesDetails>\n\t\t\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t\t\t</OccurrencesDetails>\n\t\t\t\t</Rule>\n\t\t\t</File>\n\t\t\t<File filename=\"File1\">\n\t\t\t\t<Rule>\n\t\t\t\t\t<ID>9</ID>\n\t\t\t\t\t<Name>Ensure workload has valid label values</Name>\n\t\t\t\t\t<FailSuggestion>Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it</FailSuggestion>\n\t\t\t\t\t<OccurrencesDetails>\n\t\t\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t\t\t</OccurrencesDetails>\n\t\t\t\t</Rule>\n\t\t\t</File>\n\t\t\t<File filename=\"File1\">\n\t\t\t\t<Rule>\n\t\t\t\t\t<ID>11</ID>\n\t\t\t\t\t<Name>Ensure each container has a configured liveness probe</Name>\n\t\t\t\t\t<FailSuggestion>Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks</FailSuggestion>\n\t\t\t\t\t<OccurrencesDetails>\n\t\t\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t\t\t</OccurrencesDetails>\n\t\t\t\t</Rule>\n\t\t\t</File>\n\t\t</FileNameRuleMapper>\n\t\t<Summary>\n\t\t\t<TotalFailedRules>4</TotalFailedRules>\n\t\t\t<FilesCount>1</FilesCount>\n\t\t\t<TotalPassedCount>0</TotalPassedCount>\n\t\t</Summary>\n\t</EvaluationResults>\n\t<EvaluationSummary>\n\t\t<ConfigsCount>1</ConfigsCount>\n\t\t<RulesCount>21</RulesCount>\n\t\t<FilesCount>1</FilesCount>\n\t\t<PassedYamlValidationCount>1</PassedYamlValidationCount>\n\t\t<PassedK8sValidationCount>1</PassedK8sValidationCount>\n\t\t<PassedPolicyCheckCount>0</PassedPolicyCheckCount>\n\t</EvaluationSummary>\n</FormattedOutput>\n",
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
