package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/datreeio/datree/bl/validation"

	"github.com/datreeio/datree/pkg/printer"
	"gopkg.in/yaml.v2"
)

type Printer interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintEvaluationSummary(summary printer.EvaluationSummary, k8sVersion string)
}

type FormattedOutput struct {
	Results           *EvaluationResults
	EvaluationSummary printer.EvaluationSummary
}

func PrintResults(results *EvaluationResults, invalidFiles []*validation.InvalidFile, evaluationSummary printer.EvaluationSummary, loginURL string, outputFormat string, printer Printer, k8sVersion string) error {
	switch {
	case outputFormat == "json":
		return jsonOutput(&FormattedOutput{Results: results, EvaluationSummary: evaluationSummary})
	case outputFormat == "yaml":
		return yamlOutput(&FormattedOutput{Results: results, EvaluationSummary: evaluationSummary})
	default:
		return textOutput(results, invalidFiles, evaluationSummary, loginURL, printer, k8sVersion)
	}
}

func jsonOutput(formattedOutput *FormattedOutput) error {
	jsonOutput, err := json.Marshal(formattedOutput)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(jsonOutput))
	return nil
}

func yamlOutput(formattedOutput *FormattedOutput) error {
	yamlOutput, err := yaml.Marshal(formattedOutput)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(yamlOutput))
	return nil
}

func textOutput(results *EvaluationResults, invalidFiles []*validation.InvalidFile, evaluationSummary printer.EvaluationSummary, url string, printer Printer, k8sVersion string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	warnings, err := parseToPrinterWarnings(results, invalidFiles, pwd, k8sVersion)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printer.PrintWarnings(warnings)

	summary := parseEvaluationResultsToSummary(results, evaluationSummary, url)

	printer.PrintEvaluationSummary(evaluationSummary, k8sVersion)

	printer.PrintSummaryTable(summary)

	return nil
}

func parseToPrinterWarnings(results *EvaluationResults, invalidFiles []*validation.InvalidFile, pwd string, k8sVersion string) ([]printer.Warning, error) {
	var warnings = []printer.Warning{}

	for _, invalidFile := range invalidFiles {
		warnings = append(warnings, printer.Warning{
			Title:   fmt.Sprintf(">>  File: %s\n", invalidFile.Path),
			Details: []printer.WarningInfo{},
			ValidationInfo: printer.ValidationInfo{
				IsValid:          false,
				ValidationErrors: invalidFile.ValidationErrors,
				K8sVersion:       k8sVersion,
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
					K8sVersion:       k8sVersion,
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
	configsCount := evaluationSummary.ConfigsCount
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
		{LeftCol: EvaluatedConfigurations.String(), RightCol: fmt.Sprint(configsCount), RowIndex: 1},
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
