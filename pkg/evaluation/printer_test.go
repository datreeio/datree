package evaluation

import (
	"errors"
	"os"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/defaultRules"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"

	"github.com/datreeio/datree/pkg/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPrinter struct {
	mock.Mock
}

func (m *mockPrinter) GetWarningsText(warnings []printer.Warning) string {
	m.Called(warnings)
	return ""
}

func (c *mockPrinter) GetSummaryTableText(summary printer.Summary) string {
	c.Called(summary)
	return ""
}

func (c *mockPrinter) GetEvaluationSummaryText(summary printer.EvaluationSummary, k8sVersion string) string {
	c.Called(summary, k8sVersion)
	return ""
}

type printResultsTestCaseArgs struct {
	results             FormattedResults
	additionalJUnitData AdditionalJUnitData
	invalidYamlFiles    []*extractor.InvalidFile
	invalidK8sFiles     []*extractor.InvalidFile
	evaluationSummary   printer.EvaluationSummary
	loginURL            string
	outputFormat        string
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
		printResults(""),
		printResults("json"),
		printResults("yaml"),
		printResults("xml"),
		printResults("JUnit"),
	}
	for _, tt := range tests {
		mockedPrinter := &mockPrinter{}
		mockedPrinter.On("GetWarningsText", mock.Anything, mock.Anything, mock.Anything)
		mockedPrinter.On("GetSummaryTableText", mock.Anything)
		mockedPrinter.On("GetEvaluationSummaryText", mock.Anything, mock.Anything)

		t.Run(tt.name, func(t *testing.T) {
			_ = PrintResults(&PrintResultsData{
				Results:               tt.args.results,
				AdditionalJUnitData:   tt.args.additionalJUnitData,
				InvalidYamlFiles:      tt.args.invalidYamlFiles,
				InvalidK8sFiles:       tt.args.invalidK8sFiles,
				EvaluationSummary:     tt.args.evaluationSummary,
				LoginURL:              tt.args.loginURL,
				OutputFormat:          tt.args.outputFormat,
				Printer:               mockedPrinter,
				K8sVersion:            "1.18.0",
				Verbose:               false,
				PolicyName:            "Default",
				K8sValidationWarnings: validation.K8sValidationWarningPerValidFile{},
			})

			if tt.args.outputFormat == "json" {
				mockedPrinter.AssertNotCalled(t, "GetWarningsText")
			} else if tt.args.outputFormat == "yaml" {
				mockedPrinter.AssertNotCalled(t, "GetWarningsText")
			} else if tt.args.outputFormat == "xml" {
				mockedPrinter.AssertNotCalled(t, "GetWarningsText")
			} else if tt.args.outputFormat == "JUnit" {
				mockedPrinter.AssertNotCalled(t, "GetWarningsText")
			} else {
				pwd, _ := os.Getwd()
				warnings, _ := parseToPrinterWarnings(tt.args.results.EvaluationResults, tt.args.invalidYamlFiles, tt.args.invalidK8sFiles, pwd, "1.18.0", validation.K8sValidationWarningPerValidFile{}, false)
				mockedPrinter.AssertCalled(t, "GetWarningsText", warnings)
			}
		})
	}
}

func TestCustomOutputs(t *testing.T) {
	formattedOutput := createFormattedOutput()
	additionalJUnitData := createAdditionalJUnitData()
	expectedOutputs := getExpectedOutputs()

	jsonStdout, _ := getJsonOutput(&formattedOutput)
	assert.Equal(t, expectedOutputs.json, jsonStdout)

	yamlStdout, _ := getYamlOutput(&formattedOutput)
	assert.Equal(t, expectedOutputs.yaml, yamlStdout)

	xmlStdout, _ := getXmlOutput(&formattedOutput)
	assert.Equal(t, expectedOutputs.xml, xmlStdout)

	JUnitStdout, _ := getJUnitOutput(&formattedOutput, additionalJUnitData)
	assert.Equal(t, expectedOutputs.JUnit, JUnitStdout)
}

func TestInvalidK8sCustomOutputs(t *testing.T) {
	formattedOutput := createInvalidK8sFileFormattedOutput()
	additionalJUnitData := createAdditionalJUnitDataInvalidK8sFile()
	expectedOutputs := getInvalidK8sFileExpectedOutputs()

	JUnitStdout, _ := getJUnitOutput(&formattedOutput, additionalJUnitData)
	assert.Equal(t, expectedOutputs.JUnit, JUnitStdout)
}

func createAdditionalJUnitData() AdditionalJUnitData {
	dr, err := defaultRules.GetDefaultRules()
	if err != nil {
		panic(err)
	}
	var result []cliClient.RuleData
	for _, r := range dr.Rules {
		if r.EnabledByDefault {
			result = append(result, cliClient.RuleData{
				Identifier: r.UniqueName,
				Name:       r.Name,
			})
		}
	}
	return AdditionalJUnitData{
		AllEnabledRules:            result,
		AllFilesThatRanPolicyCheck: []string{"File1", "File2"},
	}
}

func createAdditionalJUnitDataInvalidK8sFile() AdditionalJUnitData {
	dr, err := defaultRules.GetDefaultRules()
	if err != nil {
		panic(err)
	}
	var result []cliClient.RuleData
	for _, r := range dr.Rules {
		if r.EnabledByDefault {
			result = append(result, cliClient.RuleData{
				Identifier: r.UniqueName,
				Name:       r.Name,
			})
		}
	}
	return AdditionalJUnitData{
		AllEnabledRules:            result,
		AllFilesThatRanPolicyCheck: []string{},
	}
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

func createInvalidK8sFileFormattedOutput() FormattedOutput {
	err := errors.New("k8s schema validation error: could not find schema for Deploymentt You can skip files with missing schemas instead of failing by using the `--ignore-missing-schemas` flag ")
	err2 := errors.New("k8s schema validation error: For field spec.replicas: Invalid type. Expected: [integer,null], given: string ")
	invalidK8sFile := &extractor.InvalidFile{
		Path:             "File1",
		ValidationErrors: []error{err, err2},
	}
	return FormattedOutput{
		EvaluationSummary: NonInteractiveEvaluationSummary{
			ConfigsCount:                0,
			FilesCount:                  1,
			PassedYamlValidationCount:   1,
			K8sValidation:               "0/1",
			PassedPolicyValidationCount: 0,
		},
		K8sValidationResults: []*extractor.InvalidFile{invalidK8sFile},
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

func getInvalidK8sFileExpectedOutputs() expectedOutputs {
	jUnitOutput, _ := os.ReadFile("./printer_test_expected_outputs/JUnit_invalid_k8s_output.xml")
	return expectedOutputs{
		JUnit: string(jUnitOutput),
	}
}

func printResults(outputFormat string) *printResultsTestCase {
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
			additionalJUnitData: AdditionalJUnitData{
				AllEnabledRules:            []cliClient.RuleData{},
				AllFilesThatRanPolicyCheck: []string{},
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
