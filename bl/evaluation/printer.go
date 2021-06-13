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
	EvaluationResults *EvaluationResults
	EvaluationSummary printer.EvaluationSummary
	InvalidYamlFiles  []*validation.InvalidYamlFile
	InvalidK8sFiles   []*validation.InvalidK8sFile
}

func PrintResults(results *EvaluationResults, invalidYamlFiles []*validation.InvalidYamlFile, invalidK8sFiles []*validation.InvalidK8sFile, evaluationSummary printer.EvaluationSummary, loginURL string, outputFormat string, printer Printer, k8sVersion string) error {
	switch {
	case outputFormat == "json":
		return jsonOutput(&FormattedOutput{EvaluationResults: results, EvaluationSummary: evaluationSummary, InvalidYamlFiles: invalidYamlFiles, InvalidK8sFiles: invalidK8sFiles})
	case outputFormat == "yaml":
		return yamlOutput(&FormattedOutput{EvaluationResults: results, EvaluationSummary: evaluationSummary, InvalidYamlFiles: invalidYamlFiles, InvalidK8sFiles: invalidK8sFiles})
	default:
		return textOutput(results, invalidYamlFiles, invalidK8sFiles, evaluationSummary, loginURL, printer, k8sVersion)
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

func textOutput(results *EvaluationResults, invalidYamlFiles []*validation.InvalidYamlFile, invalidK8sFiles []*validation.InvalidK8sFile, evaluationSummary printer.EvaluationSummary, url string, printer Printer, k8sVersion string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	warnings, err := parseToPrinterWarnings(results, invalidYamlFiles, invalidK8sFiles, pwd, k8sVersion)
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

func parseToPrinterWarnings(results *EvaluationResults, invalidYamlFiles []*validation.InvalidYamlFile, invalidK8sFiles []*validation.InvalidK8sFile, pwd string, k8sVersion string) ([]printer.Warning, error) {
	var warnings = []printer.Warning{}

	for _, invalidFile := range invalidYamlFiles {
		warnings = append(warnings, printer.Warning{
			Title:   fmt.Sprintf(">>  File: %s\n", invalidFile.Path),
			Details: []printer.WarningInfo{},
			InvalidYamlInfo: printer.InvalidYamlInfo{
				ValidationErrors: invalidFile.ValidationErrors,
			},
		})
	}

	for _, invalidFile := range invalidK8sFiles {
		warnings = append(warnings, printer.Warning{
			Title:   fmt.Sprintf(">>  File: %s\n", invalidFile.Path),
			Details: []printer.WarningInfo{},
			InvalidK8sInfo: printer.InvalidK8sInfo{
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
				Title:           fmt.Sprintf(">>  File: %s\n", relativePath),
				Details:         warningDetails,
				InvalidYamlInfo: printer.InvalidYamlInfo{},
				InvalidK8sInfo:  printer.InvalidK8sInfo{},
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
	totalRulesEvaluated := 0
	totalFailedRules := 0
	totalPassedRules := 0

	if results != nil {
		totalRulesEvaluated = evaluationSummary.RulesCount * results.Summary.FilesCount
		totalFailedRules = results.Summary.TotalFailedRules
		totalPassedRules = totalRulesEvaluated - totalFailedRules
	}

	plainRows := []printer.SummaryItem{
		{LeftCol: EnabledRules.String(), RightCol: fmt.Sprint(evaluationSummary.RulesCount), RowIndex: 0},
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
