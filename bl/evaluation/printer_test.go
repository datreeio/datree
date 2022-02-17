package evaluation

import (
	"bytes"
	"github.com/datreeio/datree/pkg/extractor"
	"io"
	"log"
	"os"
	"testing"

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
	results           ResultType
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
				warnings, _ := parseToPrinterWarnings(tt.args.results.EvaluationResults, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, pwd, "1.18.0")
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
						}},
					},
					{
						Identifier:       "CONTAINERS_MISSING_MEMORY_LIMIT_KEY",
						Name:             "Ensure each container has a configured memory limit",
						MessageOnFailure: "Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
						}},
					},
					{
						Identifier:       "WORKLOAD_INVALID_LABELS_VALUE",
						Name:             "Ensure workload has valid label values",
						MessageOnFailure: "Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
						}},
					},
					{
						Identifier:       "CONTAINERS_MISSING_LIVENESSPROBE_KEY",
						Name:             "Ensure each container has a configured liveness probe",
						MessageOnFailure: "Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks",
						OccurrencesDetails: []OccurrenceDetails{{
							MetadataName: "rss-site",
							Kind:         "Deployment",
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
			PassedK8sValidationCount:    1,
			PassedPolicyValidationCount: 0,
		},
	}
}

func getExpectedOutputs() expectedOutputs {
	return expectedOutputs{
		json: "{\"policyValidationResults\":[{\"fileName\":\"File1\",\"ruleResults\":[{\"identifier\":\"CONTAINERS_MISSING_IMAGE_VALUE_VERSION\",\"name\":\"Ensure each container image has a pinned (tag) version\",\"messageOnFailure\":\"Incorrect value for key `image` - specify an image version to avoid unpleasant \\\"version surprises\\\" in the future\",\"occurrencesDetails\":[{\"metadataName\":\"rss-site\",\"kind\":\"Deployment\"}]},{\"identifier\":\"CONTAINERS_MISSING_MEMORY_LIMIT_KEY\",\"name\":\"Ensure each container has a configured memory limit\",\"messageOnFailure\":\"Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization\",\"occurrencesDetails\":[{\"metadataName\":\"rss-site\",\"kind\":\"Deployment\"}]},{\"identifier\":\"WORKLOAD_INVALID_LABELS_VALUE\",\"name\":\"Ensure workload has valid label values\",\"messageOnFailure\":\"Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it\",\"occurrencesDetails\":[{\"metadataName\":\"rss-site\",\"kind\":\"Deployment\"}]},{\"identifier\":\"CONTAINERS_MISSING_LIVENESSPROBE_KEY\",\"name\":\"Ensure each container has a configured liveness probe\",\"messageOnFailure\":\"Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks\",\"occurrencesDetails\":[{\"metadataName\":\"rss-site\",\"kind\":\"Deployment\"}]}]}],\"policySummary\":{\"policyName\":\"Default\",\"totalRulesInPolicy\":21,\"totalRulesFailed\":4,\"totalPassedCount\":0},\"evaluationSummary\":{\"configsCount\":1,\"filesCount\":1,\"passedYamlValidationCount\":1,\"passedK8sValidationCount\":1,\"passedPolicyValidationCount\":0},\"yamlValidationResults\":null,\"k8sValidationResults\":null}\n",
		yaml: "policyValidationResults:\n- fileName: File1\n  ruleResults:\n  - identifier: CONTAINERS_MISSING_IMAGE_VALUE_VERSION\n    name: Ensure each container image has a pinned (tag) version\n    messageOnFailure: Incorrect value for key `image` - specify an image version to\n      avoid unpleasant \"version surprises\" in the future\n    occurrencesDetails:\n    - metadataName: rss-site\n      kind: Deployment\n  - identifier: CONTAINERS_MISSING_MEMORY_LIMIT_KEY\n    name: Ensure each container has a configured memory limit\n    messageOnFailure: Missing property object `limits.memory` - value should be within\n      the accepted boundaries recommended by the organization\n    occurrencesDetails:\n    - metadataName: rss-site\n      kind: Deployment\n  - identifier: WORKLOAD_INVALID_LABELS_VALUE\n    name: Ensure workload has valid label values\n    messageOnFailure: Incorrect value for key(s) under `labels` - the vales syntax\n      is not valid so the Kubernetes engine will not accept it\n    occurrencesDetails:\n    - metadataName: rss-site\n      kind: Deployment\n  - identifier: CONTAINERS_MISSING_LIVENESSPROBE_KEY\n    name: Ensure each container has a configured liveness probe\n    messageOnFailure: Missing property object `livenessProbe` - add a properly configured\n      livenessProbe to catch possible deadlocks\n    occurrencesDetails:\n    - metadataName: rss-site\n      kind: Deployment\npolicySummary:\n  policyName: Default\n  totalRulesInPolicy: 21\n  totalRulesFailed: 4\n  totalPassedCount: 0\nevaluationSummary:\n  configsCount: 1\n  filesCount: 1\n  passedYamlValidationCount: 1\n  passedK8sValidationCount: 1\n  passedPolicyValidationCount: 0\nyamlValidationResults: []\nk8sValidationResults: []\n\n",
		xml:  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<FormattedOutput>\n\t<policyValidationResults>\n\t\t<fileName>File1</fileName>\n\t\t<ruleResults>\n\t\t\t<identifier>CONTAINERS_MISSING_IMAGE_VALUE_VERSION</identifier>\n\t\t\t<name>Ensure each container image has a pinned (tag) version</name>\n\t\t\t<messageOnFailure>Incorrect value for key `image` - specify an image version to avoid unpleasant &#34;version surprises&#34; in the future</messageOnFailure>\n\t\t\t<occurrencesDetails>\n\t\t\t\t<metadataName>rss-site</metadataName>\n\t\t\t\t<kind>Deployment</kind>\n\t\t\t</occurrencesDetails>\n\t\t</ruleResults>\n\t\t<ruleResults>\n\t\t\t<identifier>CONTAINERS_MISSING_MEMORY_LIMIT_KEY</identifier>\n\t\t\t<name>Ensure each container has a configured memory limit</name>\n\t\t\t<messageOnFailure>Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization</messageOnFailure>\n\t\t\t<occurrencesDetails>\n\t\t\t\t<metadataName>rss-site</metadataName>\n\t\t\t\t<kind>Deployment</kind>\n\t\t\t</occurrencesDetails>\n\t\t</ruleResults>\n\t\t<ruleResults>\n\t\t\t<identifier>WORKLOAD_INVALID_LABELS_VALUE</identifier>\n\t\t\t<name>Ensure workload has valid label values</name>\n\t\t\t<messageOnFailure>Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it</messageOnFailure>\n\t\t\t<occurrencesDetails>\n\t\t\t\t<metadataName>rss-site</metadataName>\n\t\t\t\t<kind>Deployment</kind>\n\t\t\t</occurrencesDetails>\n\t\t</ruleResults>\n\t\t<ruleResults>\n\t\t\t<identifier>CONTAINERS_MISSING_LIVENESSPROBE_KEY</identifier>\n\t\t\t<name>Ensure each container has a configured liveness probe</name>\n\t\t\t<messageOnFailure>Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks</messageOnFailure>\n\t\t\t<occurrencesDetails>\n\t\t\t\t<metadataName>rss-site</metadataName>\n\t\t\t\t<kind>Deployment</kind>\n\t\t\t</occurrencesDetails>\n\t\t</ruleResults>\n\t</policyValidationResults>\n\t<policySummary>\n\t\t<policyName>Default</policyName>\n\t\t<totalRulesInPolicy>21</totalRulesInPolicy>\n\t\t<totalRulesFailed>4</totalRulesFailed>\n\t\t<totalPassedCount>0</totalPassedCount>\n\t</policySummary>\n\t<evaluationSummary>\n\t\t<configsCount>1</configsCount>\n\t\t<filesCount>1</filesCount>\n\t\t<passedYamlValidationCount>1</passedYamlValidationCount>\n\t\t<passedK8sValidationCount>1</passedK8sValidationCount>\n\t\t<passedPolicyValidationCount>0</passedPolicyValidationCount>\n\t</evaluationSummary>\n</FormattedOutput>\n",
	}
}

func print_resultst(outputFormat string) *printResultsTestCase {
	return &printResultsTestCase{
		name: "Print Results Text",
		args: &printResultsTestCaseArgs{
			results: ResultType{
				EvaluationResults: &EvaluationResults{
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
