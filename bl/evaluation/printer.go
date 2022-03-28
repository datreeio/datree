package evaluation

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"

	"github.com/datreeio/datree/pkg/printer"
	"gopkg.in/yaml.v2"
)

type Printer interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
	PrintEvaluationSummary(summary printer.EvaluationSummary, k8sVersion string)
}

func PrintResults(results FormattedResults, invalidYamlFiles []*extractor.InvalidFile, invalidK8sFiles []*extractor.InvalidFile, evaluationSummary printer.EvaluationSummary, loginURL string, outputFormat string, printer Printer, k8sVersion string, policyName string, k8sValidationWarnings validation.K8sValidationWarningPerValidFile) error {
	if outputFormat == "json" || outputFormat == "yaml" || outputFormat == "xml" {
		nonInteractiveEvaluationResults := results.NonInteractiveEvaluationResults
		if nonInteractiveEvaluationResults == nil {
			nonInteractiveEvaluationResults = &NonInteractiveEvaluationResults{}
		}
		formattedOutput := FormattedOutput{
			PolicyValidationResults: nonInteractiveEvaluationResults.FormattedEvaluationResults,
			PolicySummary:           nonInteractiveEvaluationResults.PolicySummary,
			EvaluationSummary: NonInteractiveEvaluationSummary{
				ConfigsCount:                evaluationSummary.ConfigsCount,
				FilesCount:                  evaluationSummary.FilesCount,
				PassedYamlValidationCount:   evaluationSummary.PassedYamlValidationCount,
				K8sValidation:               evaluationSummary.K8sValidation,
				PassedPolicyValidationCount: evaluationSummary.PassedPolicyCheckCount,
			},
			YamlValidationResults: invalidYamlFiles,
			K8sValidationResults:  invalidK8sFiles,
		}

		if outputFormat == "json" {
			return jsonOutput(&formattedOutput)
		} else if outputFormat == "yaml" {
			return yamlOutput(&formattedOutput)
		} else {
			return xmlOutput(&formattedOutput)
		}
	} else {
		return textOutput(results.EvaluationResults, invalidYamlFiles, invalidK8sFiles, evaluationSummary, loginURL, printer, k8sVersion, policyName, k8sValidationWarnings)
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

func xmlOutput(formattedOutput *FormattedOutput) error {
	xmlOutput, err := xml.MarshalIndent(formattedOutput, "", "\t")
	xmlOutput = []byte(xml.Header + string(xmlOutput))
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(xmlOutput))
	return nil
}

func textOutput(results *EvaluationResults, invalidYamlFiles []*extractor.InvalidFile, invalidK8sFiles []*extractor.InvalidFile, evaluationSummary printer.EvaluationSummary, url string, printer Printer, k8sVersion string, policyName string, k8sValidationWarnings validation.K8sValidationWarningPerValidFile) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	warnings, err := parseToPrinterWarnings(results, invalidYamlFiles, invalidK8sFiles, pwd, k8sVersion, k8sValidationWarnings)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printer.PrintWarnings(warnings)

	summary := parseEvaluationResultsToSummary(results, evaluationSummary, url, policyName)

	printer.PrintEvaluationSummary(evaluationSummary, k8sVersion)

	printer.PrintSummaryTable(summary)

	return nil
}

func parseInvalidYamlFilesToWarnings(invalidYamlFiles []*extractor.InvalidFile) []printer.Warning {
	var warnings []printer.Warning

	for _, invalidFile := range invalidYamlFiles {
		warnings = append(warnings, printer.Warning{
			Title:       fmt.Sprintf(">>  File: %s\n", invalidFile.Path),
			FailedRules: []printer.FailedRule{},
			InvalidYamlInfo: printer.InvalidYamlInfo{
				ValidationErrors: invalidFile.ValidationErrors,
			},
		})
	}

	return warnings
}

func parseInvalidK8sFilesToWarnings(invalidK8sFiles []*extractor.InvalidFile, k8sVersion string) []printer.Warning {
	var warnings []printer.Warning

	for _, invalidFile := range invalidK8sFiles {
		warnings = append(warnings, printer.Warning{
			Title:       fmt.Sprintf(">>  File: %s\n", invalidFile.Path),
			FailedRules: []printer.FailedRule{},
			InvalidK8sInfo: printer.InvalidK8sInfo{
				ValidationErrors: invalidFile.ValidationErrors,
				K8sVersion:       k8sVersion,
			},
			ExtraMessages: GetWarningExtraMessages(invalidFile),
		})
	}

	return warnings
}

func parseToPrinterWarnings(results *EvaluationResults, invalidYamlFiles []*extractor.InvalidFile, invalidK8sFiles []*extractor.InvalidFile, pwd string, k8sVersion string, k8sValidationWarnings validation.K8sValidationWarningPerValidFile) ([]printer.Warning, error) {
	var warnings = []printer.Warning{}

	warnings = append(warnings, parseInvalidYamlFilesToWarnings(invalidYamlFiles)...)

	warnings = append(warnings, parseInvalidK8sFilesToWarnings(invalidK8sFiles, k8sVersion)...)

	if results != nil {

		filenames := []string{}
		for key := range results.FileNameRuleMapper {
			filenames = append(filenames, key)
		}
		sort.Strings(filenames)

		for _, filename := range filenames {
			rules := results.FileNameRuleMapper[filename]
			var failedRules = []printer.FailedRule{}

			rulesUniqueNames := []string{}
			for rulesUniqueName := range rules {
				rulesUniqueNames = append(rulesUniqueNames, rulesUniqueName)
			}

			for _, ruleUniqueName := range rulesUniqueNames {
				rule := rules[ruleUniqueName]
				failedRule := printer.FailedRule{
					Name:               rule.Name,
					Occurrences:        rule.GetOccurrencesCount(),
					Suggestion:         rule.MessageOnFailure,
					OccurrencesDetails: []printer.OccurrenceDetails{},
				}
				for _, occurrenceDetails := range rule.OccurrencesDetails {
					failedRule.OccurrencesDetails = append(
						failedRule.OccurrencesDetails,
						printer.OccurrenceDetails{MetadataName: occurrenceDetails.MetadataName, Kind: occurrenceDetails.Kind},
					)
				}

				failedRules = append(failedRules, failedRule)
			}

			relativePath, _ := filepath.Rel(pwd, filename)

			warnings = append(warnings, printer.Warning{
				Title:           fmt.Sprintf(">>  File: %s\n", relativePath),
				FailedRules:     failedRules,
				InvalidYamlInfo: printer.InvalidYamlInfo{},
				InvalidK8sInfo: printer.InvalidK8sInfo{
					ValidationWarning: k8sValidationWarnings[filename],
				},
			})
		}
	}

	return warnings, nil
}

func GetWarningExtraMessages(invalidFile *extractor.InvalidFile) []printer.ExtraMessage {
	var extraMessages []printer.ExtraMessage

	if IsHelmFile(invalidFile.Path) {
		extraMessages = append(extraMessages, printer.ExtraMessage{
			Text:  "Are you trying to test a raw helm file? To run Datree with Helm - check out the helm plugin README:\nhttps://github.com/datreeio/helm-datree",
			Color: "cyan",
		})
	} else if IsKustomizationFile(invalidFile.Path) {
		extraMessages = append(extraMessages, printer.ExtraMessage{
			Text:  "It seems you're trying to test a Kustomize file, Please try using `datree kustomize test` instead. For more information, check out the docs:\nhttps://hub.datree.io/kustomize-support",
			Color: "cyan",
		})
	}

	return extraMessages
}

func IsHelmFile(filePath string) bool {
	cleanFilePath := strings.Replace(filePath, "\n", "", -1)
	fileExtension := filepath.Ext(cleanFilePath)

	if fileExtension != ".yml" && fileExtension != ".yaml" {
		return false
	}

	helmFilesExtensions := [...]string{"Chart", "chart", "Values", "values"}

	for _, extension := range helmFilesExtensions {
		if strings.Contains(cleanFilePath, extension) {
			return true
		}
	}

	return false
}

func IsKustomizationFile(filePath string) bool {
	cleanFilePath := strings.Replace(filePath, "\n", "", -1)

	kustomizeFilesExtensions := [...]string{"kustomization.yml", "kustomization.yaml", "Kustomization"}

	for _, extension := range kustomizeFilesExtensions {
		if strings.Contains(cleanFilePath, extension) {
			return true
		}
	}

	return false
}

type OutputTitle int

const (
	EvaluatedConfigurations OutputTitle = iota
	TotalRulesEvaluated
	SeeAll
	TotalRulesPassed
	TotalRulesFailed
)

func (t OutputTitle) String() string {
	return [...]string{
		"Configs tested against policy",
		"Total rules evaluated",
		"See all rules in policy",
		"Total rules passed",
		"Total rules failed"}[t]
}

func buildEnabledRulesTitle(policyName string) string {
	var str strings.Builder
	fmt.Fprintf(&str, "Enabled rules in policy “%s”", policyName)
	return str.String()
}

func parseEvaluationResultsToSummary(results *EvaluationResults, evaluationSummary printer.EvaluationSummary, loginURL string, policyName string) printer.Summary {
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
		{LeftCol: buildEnabledRulesTitle(policyName), RightCol: fmt.Sprint(evaluationSummary.RulesCount), RowIndex: 0},
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
