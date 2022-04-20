package evaluation

import (
	"encoding/xml"
	"strconv"
)

type JUnitOutput struct {
	XMLName  xml.Name    `xml:"testsuites"`
	Name     string      `xml:"name,attr"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Skipped  int         `xml:"skipped,attr"`
	Suites   []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Name       string      `xml:"name,attr"`
	Properties *[]Property `xml:"properties>property,omitempty"`
	TestCases  []TestCase  `xml:"testcase"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TestCase struct {
	Name      string `xml:"name,attr"`
	ClassName string `xml:"classname,attr"`
	Skipped   *TestCaseSkipped
	Failure   *TestCaseFailure
}

type TestCaseSkipped struct {
	XMLName xml.Name `xml:"skipped,omitempty"`
	Message string   `xml:"message,attr"`
}

type TestCaseFailure struct {
	XMLName xml.Name `xml:"failure,omitempty"`
	Message string   `xml:"message,attr"`
	Content string   `xml:",chardata"`
}

func FormattedOutputToJUnitOutput(formattedOutput FormattedOutput) JUnitOutput {
	jUnitOutput := JUnitOutput{
		Name:     formattedOutput.PolicySummary.PolicyName,
		Tests:    formattedOutput.PolicySummary.TotalRulesInPolicy,
		Failures: formattedOutput.PolicySummary.TotalRulesFailed,
		Skipped:  formattedOutput.PolicySummary.TotalSkippedRules,
		Suites:   []TestSuite{},
	}

	for _, policyValidationResult := range formattedOutput.PolicyValidationResults {
		suite := TestSuite{
			Name:      policyValidationResult.FileName,
			TestCases: []TestCase{},
		}

		for _, ruleResult := range policyValidationResult.RuleResults {
			testCase := TestCase{
				Name:      ruleResult.Name,
				ClassName: ruleResult.Identifier,
			}
			testCase.Failure = &TestCaseFailure{
				Message: ruleResult.MessageOnFailure,
				Content: getContentFromOccurrencesDetails(ruleResult.OccurrencesDetails),
			}
			if areAllOccurrencesSkipped(ruleResult.OccurrencesDetails) {
				testCase.Skipped = &TestCaseSkipped{Message: "all failing configs"}
			}
			suite.TestCases = append(suite.TestCases, testCase)
		}
		jUnitOutput.Suites = append(jUnitOutput.Suites, suite)
	}

	jUnitOutput.Suites = append(jUnitOutput.Suites, TestSuite{
		Name: "policySummary",
		Properties: &[]Property{{
			Name:  "policyName",
			Value: formattedOutput.PolicySummary.PolicyName,
		}, {
			Name:  "totalRulesInPolicy",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalRulesInPolicy),
		}, {
			Name:  "totalRulesFailed",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalRulesFailed),
		}, {
			Name:  "totalSkippedRules",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalSkippedRules),
		}, {
			Name:  "totalPassedCount",
			Value: strconv.Itoa(formattedOutput.PolicySummary.TotalPassedCount),
		}},
	})
	jUnitOutput.Suites = append(jUnitOutput.Suites, TestSuite{
		Name: "evaluationSummary",
		Properties: &[]Property{{
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
	})

	return jUnitOutput
}

func getContentFromOccurrencesDetails(occurrencesDetails []OccurrenceDetails) string {
	totalOccurrences := 0
	totalSkipped := 0
	var occurrencesLines string
	var skipLines string

	for _, occurrenceDetails := range occurrencesDetails {
		currentLine := "— metadata.name: " + occurrenceDetails.MetadataName + " (kind: " + occurrenceDetails.Kind + ")\n"

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
