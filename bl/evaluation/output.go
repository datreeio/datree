package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/pkg/printer"
	"gopkg.in/yaml.v2"
)

type Printer interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
}

func JSONOutput(results *EvaluationResults) error {
	jsonOutput, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(jsonOutput))
	return nil
}

func YAMLOutput(results *EvaluationResults) error {
	yamlOutput, err := yaml.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(yamlOutput))
	return nil
}

func TextOutput(results *EvaluationResults, url string, printer Printer) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	warnings, err := parseToPrinterWarnings(results, pwd)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printer.PrintWarnings(warnings)
	summary := parseEvaluationResultsToSummary(results, url)

	printer.PrintSummaryTable(summary)

	if results.Summary.TotalFailedRules > 0 {
		return fmt.Errorf("failed rules count is %d (>0)", results.Summary.TotalFailedRules)
	}
	return nil
}

func parseToPrinterWarnings(results *EvaluationResults, pwd string) ([]printer.Warning, error) {
	var warnings = []printer.Warning{}

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
		})
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
func parseEvaluationResultsToSummary(results *evaluation.EvaluationResults, loginURL string) printer.Summary {
	totalRulesEvaluated := results.Summary.RulesCount * results.Summary.FilesCount
	plainRows := []printer.SummaryItem{
		{LeftCol: EnabledRules.String(), RightCol: fmt.Sprint(results.Summary.RulesCount), RowIndex: 0},
		{LeftCol: EvaluatedConfigurations.String(), RightCol: fmt.Sprint(results.Summary.FilesCount), RowIndex: 1},
		{LeftCol: TotalRulesEvaluated.String(), RightCol: fmt.Sprint(totalRulesEvaluated), RowIndex: 2},
		{LeftCol: SeeAll.String(), RightCol: loginURL, RowIndex: 5},
	}

	successRow := printer.SummaryItem{LeftCol: TotalRulesPassed.String(), RightCol: fmt.Sprint(totalRulesEvaluated - results.Summary.TotalFailedRules), RowIndex: 4}
	errorRow := printer.SummaryItem{LeftCol: TotalRulesFailed.String(), RightCol: fmt.Sprint(results.Summary.TotalFailedRules), RowIndex: 3}

	summary := &printer.Summary{
		ErrorRow:   errorRow,
		SuccessRow: successRow,
		PlainRows:  plainRows,
	}
	return *summary
}
