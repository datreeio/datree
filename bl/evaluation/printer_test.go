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
	results           ResultType
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
		PolicySummary: PolicySummary{
			PolicyName:         "Default",
			TotalRulesInPolicy: 21,
			TotalRulesFailed:   4,
			TotalPassedCount:   0,
		},
		FormattedEvaluationResults: []*FormattedEvaluationResults{&FormattedEvaluationResults{
			FileName: "File1",
			RuleRresults: []*RulerObject{&RulerObject{
				Identifier:       "CONTAINERS_MISSING_IMAGE_VALUE_VERSION",
				Name:             "Ensure each container image has a pinned (tag) version",
				MessageOnFailure: "Incorrect value for key `image` - specify an image version to avoid unpleasant \"version surprises\" in the future",
				OccurrencesDetails: []OccurrenceDetails{{
					MetadataName: "rss-site",
					Kind:         "Deployment",
				}},
			},
				&RulerObject{
					Identifier:       "CONTAINERS_MISSING_MEMORY_LIMIT_KEY",
					Name:             "Ensure each container has a configured memory limit",
					MessageOnFailure: "Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization",
					OccurrencesDetails: []OccurrenceDetails{{
						MetadataName: "rss-site",
						Kind:         "Deployment",
					}},
				},
				&RulerObject{
					Identifier:       "WORKLOAD_INVALID_LABELS_VALUE",
					Name:             "Ensure workload has valid label values",
					MessageOnFailure: "Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it",
					OccurrencesDetails: []OccurrenceDetails{{
						MetadataName: "rss-site",
						Kind:         "Deployment",
					}},
				},
				&RulerObject{
					Identifier:       "CONTAINERS_MISSING_LIVENESSPROBE_KEY",
					Name:             "Ensure each container has a configured liveness probe",
					MessageOnFailure: "Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks",
					OccurrencesDetails: []OccurrenceDetails{{
						MetadataName: "rss-site",
						Kind:         "Deployment",
					}},
				},
			},
		}},
	}

	return FormattedOutput{
		PolicyValidationResults: evaluationResults.FormattedEvaluationResults,
		PolicySummary:           evaluationResults.PolicySummary,
		EvaluationSummary: struct {
			ConfigsCount              int
			FilesCount                int
			PassedYamlValidationCount int
			PassedK8sValidationCount  int
			PassedPolicyCheckCount    int
		}{
			ConfigsCount:              1,
			FilesCount:                1,
			PassedYamlValidationCount: 1,
			PassedK8sValidationCount:  1,
			PassedPolicyCheckCount:    0,
		},
	}
}

func getExpectedOutputs() expectedOutputs {
	return expectedOutputs{
		json: "{\"PolicyValidationResults\":[{\"FileName\":\"File1\",\"RuleRresults\":[{\"Identifier\":\"CONTAINERS_MISSING_IMAGE_VALUE_VERSION\",\"Name\":\"Ensure each container image has a pinned (tag) version\",\"MessageOnFailure\":\"Incorrect value for key `image` - specify an image version to avoid unpleasant \\\"version surprises\\\" in the future\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]},{\"Identifier\":\"CONTAINERS_MISSING_MEMORY_LIMIT_KEY\",\"Name\":\"Ensure each container has a configured memory limit\",\"MessageOnFailure\":\"Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]},{\"Identifier\":\"WORKLOAD_INVALID_LABELS_VALUE\",\"Name\":\"Ensure workload has valid label values\",\"MessageOnFailure\":\"Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]},{\"Identifier\":\"CONTAINERS_MISSING_LIVENESSPROBE_KEY\",\"Name\":\"Ensure each container has a configured liveness probe\",\"MessageOnFailure\":\"Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks\",\"OccurrencesDetails\":[{\"MetadataName\":\"rss-site\",\"Kind\":\"Deployment\"}]}]}],\"PolicySummary\":{\"PolicyName\":\"Default\",\"TotalRulesInPolicy\":21,\"TotalRulesFailed\":4,\"TotalPassedCount\":0},\"EvaluationSummary\":{\"ConfigsCount\":1,\"FilesCount\":1,\"PassedYamlValidationCount\":1,\"PassedK8sValidationCount\":1,\"PassedPolicyCheckCount\":0},\"YamlValidationResults\":null,\"K8sValidationResults\":null}\n",
		yaml: "policyvalidationresults:\n- filename: File1\n  rulerresults:\n  - identifier: CONTAINERS_MISSING_IMAGE_VALUE_VERSION\n    name: Ensure each container image has a pinned (tag) version\n    messageonfailure: Incorrect value for key `image` - specify an image version to\n      avoid unpleasant \"version surprises\" in the future\n    occurrencesdetails:\n    - metadataname: rss-site\n      kind: Deployment\n  - identifier: CONTAINERS_MISSING_MEMORY_LIMIT_KEY\n    name: Ensure each container has a configured memory limit\n    messageonfailure: Missing property object `limits.memory` - value should be within\n      the accepted boundaries recommended by the organization\n    occurrencesdetails:\n    - metadataname: rss-site\n      kind: Deployment\n  - identifier: WORKLOAD_INVALID_LABELS_VALUE\n    name: Ensure workload has valid label values\n    messageonfailure: Incorrect value for key(s) under `labels` - the vales syntax\n      is not valid so the Kubernetes engine will not accept it\n    occurrencesdetails:\n    - metadataname: rss-site\n      kind: Deployment\n  - identifier: CONTAINERS_MISSING_LIVENESSPROBE_KEY\n    name: Ensure each container has a configured liveness probe\n    messageonfailure: Missing property object `livenessProbe` - add a properly configured\n      livenessProbe to catch possible deadlocks\n    occurrencesdetails:\n    - metadataname: rss-site\n      kind: Deployment\npolicysummary:\n  policyname: Default\n  totalrulesinpolicy: 21\n  totalrulesfailed: 4\n  totalpassedcount: 0\nevaluationsummary:\n  configscount: 1\n  filescount: 1\n  passedyamlvalidationcount: 1\n  passedk8svalidationcount: 1\n  passedpolicycheckcount: 0\nyamlvalidationresults: []\nk8svalidationresults: []\n\n",
		xml:  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<FormattedOutput>\n\t<PolicyValidationResults>\n\t\t<FileName>File1</FileName>\n\t\t<RuleRresults>\n\t\t\t<Identifier>CONTAINERS_MISSING_IMAGE_VALUE_VERSION</Identifier>\n\t\t\t<Name>Ensure each container image has a pinned (tag) version</Name>\n\t\t\t<MessageOnFailure>Incorrect value for key `image` - specify an image version to avoid unpleasant &#34;version surprises&#34; in the future</MessageOnFailure>\n\t\t\t<OccurrencesDetails>\n\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t</OccurrencesDetails>\n\t\t</RuleRresults>\n\t\t<RuleRresults>\n\t\t\t<Identifier>CONTAINERS_MISSING_MEMORY_LIMIT_KEY</Identifier>\n\t\t\t<Name>Ensure each container has a configured memory limit</Name>\n\t\t\t<MessageOnFailure>Missing property object `limits.memory` - value should be within the accepted boundaries recommended by the organization</MessageOnFailure>\n\t\t\t<OccurrencesDetails>\n\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t</OccurrencesDetails>\n\t\t</RuleRresults>\n\t\t<RuleRresults>\n\t\t\t<Identifier>WORKLOAD_INVALID_LABELS_VALUE</Identifier>\n\t\t\t<Name>Ensure workload has valid label values</Name>\n\t\t\t<MessageOnFailure>Incorrect value for key(s) under `labels` - the vales syntax is not valid so the Kubernetes engine will not accept it</MessageOnFailure>\n\t\t\t<OccurrencesDetails>\n\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t</OccurrencesDetails>\n\t\t</RuleRresults>\n\t\t<RuleRresults>\n\t\t\t<Identifier>CONTAINERS_MISSING_LIVENESSPROBE_KEY</Identifier>\n\t\t\t<Name>Ensure each container has a configured liveness probe</Name>\n\t\t\t<MessageOnFailure>Missing property object `livenessProbe` - add a properly configured livenessProbe to catch possible deadlocks</MessageOnFailure>\n\t\t\t<OccurrencesDetails>\n\t\t\t\t<MetadataName>rss-site</MetadataName>\n\t\t\t\t<Kind>Deployment</Kind>\n\t\t\t</OccurrencesDetails>\n\t\t</RuleRresults>\n\t</PolicyValidationResults>\n\t<PolicySummary>\n\t\t<PolicyName>Default</PolicyName>\n\t\t<TotalRulesInPolicy>21</TotalRulesInPolicy>\n\t\t<TotalRulesFailed>4</TotalRulesFailed>\n\t\t<TotalPassedCount>0</TotalPassedCount>\n\t</PolicySummary>\n\t<EvaluationSummary>\n\t\t<ConfigsCount>1</ConfigsCount>\n\t\t<FilesCount>1</FilesCount>\n\t\t<PassedYamlValidationCount>1</PassedYamlValidationCount>\n\t\t<PassedK8sValidationCount>1</PassedK8sValidationCount>\n\t\t<PassedPolicyCheckCount>0</PassedPolicyCheckCount>\n\t</EvaluationSummary>\n</FormattedOutput>\n",
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
					PolicySummary: PolicySummary{
						PolicyName:         "Default",
						TotalRulesInPolicy: 0,
						TotalRulesFailed:   0,
						TotalPassedCount:   0,
					},
					FormattedEvaluationResults: []*FormattedEvaluationResults{},
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
