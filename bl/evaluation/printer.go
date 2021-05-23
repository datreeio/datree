package evaluation

import (
	"encoding/json"
	"fmt"
	"github.com/datreeio/datree/bl/validation"
	"os"
	"path/filepath"

	"github.com/datreeio/datree/pkg/printer"
	"gopkg.in/yaml.v2"
)

type Printer interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintEvaluationSummary(summary printer.EvaluationSummary)
}

// func PrintAllResults(results *EvaluationResults /*, invalidFiles*/) {
// 	// TODO: foreach invalid file
// }

// url := "https://app.datree.io/login?cliId=" + cliId
func PrintResults(results *EvaluationResults, invalidFiles []*validation.InvalidFile, evaluationSummary printer.EvaluationSummary, loginURL string, outputFormat string, printer Printer) error {
	switch {
	case outputFormat == "json":
		return jsonOutput(results)
	case outputFormat == "yaml":
		return yamlOutput(results)
	default:
		return textOutput(results, invalidFiles, evaluationSummary, loginURL, printer)
	}
}

func jsonOutput(results *EvaluationResults) error {
	jsonOutput, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(jsonOutput))
	return nil
}

func yamlOutput(results *EvaluationResults) error {
	yamlOutput, err := yaml.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(yamlOutput))
	return nil
}

func textOutput(results *EvaluationResults, invalidFiles []*validation.InvalidFile, evaluationSummary printer.EvaluationSummary, url string, printer Printer) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	warnings, err := parseToPrinterWarnings(results, invalidFiles, pwd)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printer.PrintWarnings(warnings)

	summary := parseEvaluationResultsToSummary(results, evaluationSummary, url)

	printer.PrintEvaluationSummary(evaluationSummary)

	printer.PrintSummaryTable(summary)

	return nil
}

func parseToPrinterWarnings(results *EvaluationResults, invalidFiles []*validation.InvalidFile, pwd string) ([]printer.Warning, error) {
	var warnings = []printer.Warning{}

	for _, invalidFile := range invalidFiles {
		warnings = append(warnings, printer.Warning{
			Title:   fmt.Sprintf(">>  File: %s\n", invalidFile.Path),
			Details: []printer.WarningInfo{},
			ValidationInfo: printer.ValidationInfo{
				IsValid:          false,
				ValidationErrors: invalidFile.ValidationErrors,
				K8sVersion:       "1.18.0",
			},
		})
	}

	if results != nil {
		for fileName, rules := range results.FileNameRuleMapper {
			var warningDetails = []printer.WarningInfo{}
			for _, rule := range rules {
				details := printer.WarningInfo{
					Caption:     rule.Name,
					Occurrences: rule.Count,
					Suggestion:  rule.FailSuggestion,
				}

				warningDetails = append(warningDetails, details)
			}

			relativePath, _ := filepath.Rel(pwd, fileName)

			warnings = append(warnings, printer.Warning{
				Title:   fmt.Sprintf(">>  File: %s\n", relativePath),
				Details: warningDetails,
				ValidationInfo: printer.ValidationInfo{
					IsValid:          true,
					ValidationErrors: []error{},
					K8sVersion:       "1.18.0",
				},
			})
		}
	}

	return warnings, nil
}

type OutputTitle int

const (
	EnabledRules OutputTitle = iota
	EvaluatedConfigurations
	TotalRulesEvaluated
	SeeAll
	TotalRulesPassed
	TotalRulesFailed
)

func (t OutputTitle) String() string {
	return [...]string{
		"Enabled rules in policy “default”",
		"Configs tested against policy",
		"Total rules evaluated",
		"See all rules in policy",
		"Total rules passed",
		"Total rules failed"}[t]
}
func parseEvaluationResultsToSummary(results *EvaluationResults, evaluationSummary printer.EvaluationSummary, loginURL string) printer.Summary {
	filesCount := evaluationSummary.FilesCount
	rulesCount := 0
	totalRulesEvaluated := 0
	totalFailedRules := 0
	totalPassedRules := 0

	if results != nil {
		rulesCount = results.Summary.RulesCount
		totalRulesEvaluated = results.Summary.RulesCount * results.Summary.FilesCount
		totalFailedRules = results.Summary.TotalFailedRules
		totalPassedRules = totalRulesEvaluated - totalFailedRules
	}

	var rulesCountRightCol string
	if rulesCount == 0 {
		rulesCountRightCol = fmt.Sprint("N/A")
	} else {
		rulesCountRightCol = fmt.Sprint(rulesCount)
	}

	plainRows := []printer.SummaryItem{
		{LeftCol: EnabledRules.String(), RightCol: rulesCountRightCol, RowIndex: 0},
		{LeftCol: EvaluatedConfigurations.String(), RightCol: fmt.Sprint(filesCount), RowIndex: 1},
		{LeftCol: TotalRulesEvaluated.String(), RightCol: fmt.Sprint(totalRulesEvaluated), RowIndex: 2},
		{LeftCol: SeeAll.String(), RightCol: loginURL, RowIndex: 5},
	}

	successRow := printer.SummaryItem{LeftCol: TotalRulesPassed.String(), RightCol: fmt.Sprint(totalPassedRules), RowIndex: 4}
	errorRow := printer.SummaryItem{LeftCol: TotalRulesFailed.String(), RightCol: fmt.Sprint(totalFailedRules), RowIndex: 3}

	summary := &printer.Summary{
		ErrorRow:   errorRow,
		SuccessRow: successRow,
		PlainRows:  plainRows,
	}
	return *summary
}
