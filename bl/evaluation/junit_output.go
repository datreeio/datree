package evaluation

import (
	"encoding/xml"
	"strconv"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/utils"
)

// JUnit specifications:
// https://llg.cubic.org/docs/junit/
// https://www.ibm.com/docs/en/developer-for-zos/14.2.0?topic=formats-junit-xml-format

type JUnitOutput struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Name       string      `xml:"name,attr"`
	Tests      int         `xml:"tests,attr"`
	Failures   int         `xml:"failures,attr"`
	Skipped    int         `xml:"skipped,attr"`
	TestSuites []testSuite `xml:"testsuite"`
}

type testSuite struct {
	Name       string      `xml:"name,attr"`
	Properties *[]property `xml:"properties>property,omitempty"`
	TestCases  []testCase  `xml:"testcase"`
}

type property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type testCase struct {
	Name      string `xml:"name,attr"`
	ClassName string `xml:"classname,attr"`
	Skipped   *skipped
	Failure   *failure
}

type skipped struct {
	XMLName xml.Name `xml:"skipped,omitempty"`
	Message string   `xml:"message,attr"`
}

type failure struct {
	XMLName xml.Name `xml:"failure,omitempty"`
	Message string   `xml:"message,attr"`
	Content string   `xml:",chardata"`
}

type AdditionalJUnitData struct {
	AllEnabledRules            []cliClient.RuleData
	AllFilesThatRanPolicyCheck []string
}

func FormattedOutputToJUnitOutput(formattedOutput FormattedOutput, additionalJUnitData AdditionalJUnitData) JUnitOutput {
	var jUnitOutput JUnitOutput

	if formattedOutput.PolicySummary != nil {
		jUnitOutput = JUnitOutput{
			Name:       formattedOutput.PolicySummary.PolicyName,
			Tests:      formattedOutput.PolicySummary.TotalRulesInPolicy,
			Failures:   formattedOutput.PolicySummary.TotalRulesFailed,
			Skipped:    formattedOutput.PolicySummary.TotalSkippedRules,
			TestSuites: []testSuite{},
		}
	} else {
		jUnitOutput = JUnitOutput{
			TestSuites: []testSuite{},
		}
	}

	if formattedOutput.YamlValidationResults != nil && len(formattedOutput.YamlValidationResults) > 0 {
		jUnitOutput.TestSuites = append(jUnitOutput.TestSuites, getInvalidYamlFilesTestSuite(formattedOutput)...)
	}

	if formattedOutput.K8sValidationResults != nil && len(formattedOutput.K8sValidationResults) > 0 {
		jUnitOutput.TestSuites = append(jUnitOutput.TestSuites, getInvalidK8sFilesTestSuite(formattedOutput)...)
	}

	for _, fileThatRanPolicyCheck := range additionalJUnitData.AllFilesThatRanPolicyCheck {
		policyValidationResult := findFileInPolicyValidationResults(fileThatRanPolicyCheck, formattedOutput.PolicyValidationResults)

		if policyValidationResult != nil {
			jUnitOutput.TestSuites = append(jUnitOutput.TestSuites, getPolicyValidationResultTestSuite(policyValidationResult, additionalJUnitData.AllEnabledRules))
		} else {
			jUnitOutput.TestSuites = append(jUnitOutput.TestSuites, getPassingFileTestSuite(fileThatRanPolicyCheck, additionalJUnitData.AllEnabledRules))
		}
	}

	if formattedOutput.PolicySummary != nil {
		jUnitOutput.TestSuites = append(jUnitOutput.TestSuites, getPolicySummaryTestSuite(formattedOutput))
	}

	jUnitOutput.TestSuites = append(jUnitOutput.TestSuites, getEvaluationSummaryTestSuite(formattedOutput))

	return jUnitOutput
}

func getPassingFileTestSuite(fileName string, allEnabledRules []cliClient.RuleData) testSuite {
	return testSuite{
		Name: fileName,
		TestCases: utils.MapSlice[cliClient.RuleData, testCase](allEnabledRules, func(ruleData cliClient.RuleData) testCase {
			return testCase{
				Name:      ruleData.Name,
				ClassName: ruleData.Identifier,
				Skipped:   nil,
				Failure:   nil,
			}
		}),
	}
}

func getPolicyValidationResultTestSuite(policyValidationResult *FormattedEvaluationResults, allEnabledRules []cliClient.RuleData) testSuite {
	suite := testSuite{
		Name:      policyValidationResult.FileName,
		TestCases: []testCase{},
	}

	for _, rule := range allEnabledRules {
		testCase := testCase{
			Name:      rule.Name,
			ClassName: rule.Identifier,
		}
		ruleResult := findRuleResult(rule, policyValidationResult.RuleResults)

		if ruleResult != nil {
			testCase.Failure = &failure{
				Message: ruleResult.MessageOnFailure,
				Content: getContentFromOccurrencesDetails(ruleResult.OccurrencesDetails),
			}
			if areAllOccurrencesSkipped(ruleResult.OccurrencesDetails) {
				testCase.Skipped = &skipped{Message: "All failing configs skipped"}
			}
		}
		suite.TestCases = append(suite.TestCases, testCase)
	}

	return suite
}

func findRuleResult(ruleData cliClient.RuleData, ruleResults []*RuleResult) *RuleResult {
	for _, ruleResult := range ruleResults {
		if ruleResult.Identifier == ruleData.Identifier {
			return ruleResult
		}
	}
	return nil
}

func getPolicySummaryTestSuite(formattedOutput FormattedOutput) testSuite {
	return testSuite{
		Name: "policySummary",
		Properties: &[]property{{
			Name:  "policyName",
			Value: formattedOutput.PolicySummary.PolicyName,
		}, {
			Name:  "totalRulesInPolicy",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalRulesInPolicy),
		}, {
			Name:  "totalSkippedRules",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalSkippedRules),
		}, {
			Name:  "totalRulesFailed",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalRulesFailed),
		}, {
			Name:  "totalPassedCount",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalPassedCount),
		}},
	}
}

func getEvaluationSummaryTestSuite(formattedOutput FormattedOutput) testSuite {
	return testSuite{
		Name: "evaluationSummary",
		Properties: &[]property{{
			Name:  "configsCount",
			Value: strconv.Itoa(formattedOutput.EvaluationSummary.ConfigsCount),
		}, {
			Name:  "filesCount",
			Value: strconv.Itoa(formattedOutput.EvaluationSummary.FilesCount),
		}, {
			Name:  "passedYamlValidationCount",
			Value: strconv.Itoa(formattedOutput.EvaluationSummary.PassedYamlValidationCount),
		}, {
			Name:  "k8sValidation",
			Value: formattedOutput.EvaluationSummary.K8sValidation,
		}, {
			Name:  "passedPolicyValidationCount",
			Value: strconv.Itoa(formattedOutput.EvaluationSummary.PassedPolicyValidationCount),
		}},
	}
}

func getContentFromOccurrencesDetails(occurrencesDetails []OccurrenceDetails) string {
	totalOccurrences := 0
	totalSkipped := 0
	var occurrencesLines string
	var skipLines string

	for _, occurrenceDetails := range occurrencesDetails {
		currentLine := "- metadata.name: " + occurrenceDetails.MetadataName + " (kind: " + occurrenceDetails.Kind + ")\n"

		totalOccurrences += occurrenceDetails.Occurrences
		occurrencesLines += currentLine
		if occurrenceDetails.IsSkipped {
			totalSkipped++
			skipLines += currentLine
		}
	}

	return strconv.Itoa(totalOccurrences) + " occurrences\n" + occurrencesLines + strconv.Itoa(totalSkipped) + " skipped\n" + skipLines
}

func areAllOccurrencesSkipped(occurrencesDetails []OccurrenceDetails) bool {
	for _, occurrenceDetails := range occurrencesDetails {
		if !occurrenceDetails.IsSkipped {
			return false
		}
	}
	return true
}

func findFileInPolicyValidationResults(fileName string, policyValidationResults []*FormattedEvaluationResults) *FormattedEvaluationResults {
	for _, policyValidationResult := range policyValidationResults {
		if policyValidationResult.FileName == fileName {
			return policyValidationResult
		}
	}
	return nil
}

func getInvalidYamlFilesTestSuite(formattedOutput FormattedOutput) []testSuite {
	var suites []testSuite

	for _, invalidYamlFile := range formattedOutput.YamlValidationResults {
		suite := testSuite{
			Name: invalidYamlFile.Path,
			TestCases: []testCase{
				{
					Name:      "invalid yaml file",
					ClassName: "yaml validation",
					Skipped:   nil,
					Failure: &failure{
						Message: "Invalid yaml file",
						Content: invalidYamlFile.ValidationErrors[0].Error(),
					},
				},
			},
		}
		suites = append(suites, suite)
	}

	return suites
}

func getInvalidK8sFilesTestSuite(formattedOutput FormattedOutput) []testSuite {
	var suites []testSuite

	for _, invalidK8sFile := range formattedOutput.K8sValidationResults {
		suite := testSuite{
			Name: invalidK8sFile.Path,
		}

		for _, k8sError := range invalidK8sFile.ValidationErrors {
			testCase := testCase{
				Name:      "invalid k8s file",
				ClassName: "k8s validation",
				Skipped:   nil,
				Failure: &failure{
					Message: "Invalid k8s file",
					Content: k8sError.Error(),
				},
			}
			suite.TestCases = append(suite.TestCases, testCase)
		}

		suites = append(suites, suite)
	}

	return suites
}
