package evaluation

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/fatih/color"

	"github.com/datreeio/datree/pkg/printer"
	"gopkg.in/yaml.v2"
)

var out = color.Output

type Printer interface {
	GetWarningsText(warnings []printer.Warning) string
	GetSummaryTableText(summary printer.Summary) string
	GetEvaluationSummaryText(summary printer.EvaluationSummary, k8sVersion string) string
}

type PrintResultsData struct {
	Results               FormattedResults
	AdditionalJUnitData   AdditionalJUnitData
	InvalidYamlFiles      []*extractor.InvalidFile
	InvalidK8sFiles       []*extractor.InvalidFile
	EvaluationSummary     printer.EvaluationSummary
	LoginURL              string
	OutputFormat          string
	Printer               Printer
	K8sVersion            string
	Verbose               bool
	PolicyName            string
	K8sValidationWarnings validation.K8sValidationWarningPerValidFile
}

type textOutputData struct {
	results               *EvaluationResults
	invalidYamlFiles      []*extractor.InvalidFile
	invalidK8sFiles       []*extractor.InvalidFile
	evaluationSummary     printer.EvaluationSummary
	url                   string
	printer               Printer
	k8sVersion            string
	Verbose               bool
	policyName            string
	k8sValidationWarnings validation.K8sValidationWarningPerValidFile
}

func PrintResults(resultsData *PrintResultsData) error {
	resultsText, err := GetResultsText(resultsData)
	if err != nil {
		return err
	}

	_, err = out.Write([]byte(resultsText))
	if err != nil {
		return err
	}
	return nil
}

func GetResultsText(resultsData *PrintResultsData) (string, error) {
	if IsFormattedOutputOption(resultsData.OutputFormat) {
		nonInteractiveEvaluationResults := resultsData.Results.NonInteractiveEvaluationResults
		if nonInteractiveEvaluationResults == nil {
			nonInteractiveEvaluationResults = &NonInteractiveEvaluationResults{}
		}
		formattedOutput := FormattedOutput{
			PolicyValidationResults: nonInteractiveEvaluationResults.FormattedEvaluationResults,
			PolicySummary:           nonInteractiveEvaluationResults.PolicySummary,
			EvaluationSummary: NonInteractiveEvaluationSummary{
				ConfigsCount:                resultsData.EvaluationSummary.ConfigsCount,
				FilesCount:                  resultsData.EvaluationSummary.FilesCount,
				PassedYamlValidationCount:   resultsData.EvaluationSummary.PassedYamlValidationCount,
				K8sValidation:               resultsData.EvaluationSummary.K8sValidation,
				PassedPolicyValidationCount: resultsData.EvaluationSummary.PassedPolicyCheckCount,
			},
			YamlValidationResults: resultsData.InvalidYamlFiles,
			K8sValidationResults:  resultsData.InvalidK8sFiles,
		}

		switch resultsData.OutputFormat {
		case "json":
			return getJsonOutput(&formattedOutput)
		case "yaml":
			return getYamlOutput(&formattedOutput)
		case "xml":
			return getXmlOutput(&formattedOutput)
		case "JUnit":
			return getJUnitOutput(&formattedOutput, resultsData.AdditionalJUnitData)
		default:
			panic(errors.New("invalid output format"))
		}
	} else {
		return getTextOutput(textOutputData{
			results:               resultsData.Results.EvaluationResults,
			invalidYamlFiles:      resultsData.InvalidYamlFiles,
			invalidK8sFiles:       resultsData.InvalidK8sFiles,
			evaluationSummary:     resultsData.EvaluationSummary,
			url:                   resultsData.LoginURL,
			printer:               resultsData.Printer,
			k8sVersion:            resultsData.K8sVersion,
			policyName:            resultsData.PolicyName,
			Verbose:               resultsData.Verbose,
			k8sValidationWarnings: resultsData.K8sValidationWarnings,
		})
	}
}

func getJsonOutput(formattedOutput *FormattedOutput) (string, error) {
	jsonOutput, err := json.Marshal(formattedOutput)
	if err != nil {
		return "", err
	}

	return fmt.Sprintln(string(jsonOutput)), nil
}

func getYamlOutput(formattedOutput *FormattedOutput) (string, error) {
	yamlOutput, err := yaml.Marshal(formattedOutput)
	if err != nil {
		return "", err
	}

	return fmt.Sprintln(string(yamlOutput)), nil
}

func getXmlOutput(formattedOutput *FormattedOutput) (string, error) {
	return convertStructToXml(formattedOutput)
}

func getJUnitOutput(formattedOutput *FormattedOutput, additionalJUnitData AdditionalJUnitData) (string, error) {
	return convertStructToXml(FormattedOutputToJUnitOutput(*formattedOutput, additionalJUnitData))
}

func convertStructToXml(output interface{}) (string, error) {
	xmlOutput, err := xml.MarshalIndent(output, "", "\t")
	xmlOutput = []byte(xml.Header + string(xmlOutput))
	if err != nil {
		return "", err
	}

	return fmt.Sprintln(string(xmlOutput)), nil
}

func getTextOutput(outputData textOutputData) (string, error) {
	sb := strings.Builder{}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	warnings, err := parseToPrinterWarnings(outputData.results, outputData.invalidYamlFiles, outputData.invalidK8sFiles, pwd, outputData.k8sVersion, outputData.k8sValidationWarnings, outputData.Verbose)
	if err != nil {
		return "", err
	}

	warningsText := outputData.printer.GetWarningsText(warnings)
	sb.WriteString(warningsText)

	summary := parseEvaluationResultsToSummary(outputData.results, outputData.evaluationSummary, outputData.url, outputData.policyName)

	evaluationSummaryText := outputData.printer.GetEvaluationSummaryText(outputData.evaluationSummary, outputData.k8sVersion)
	sb.WriteString(evaluationSummaryText)

	summaryTableText := outputData.printer.GetSummaryTableText(summary)
	sb.WriteString(summaryTableText)

	return sb.String(), nil
}

func parseInvalidYamlFilesToWarnings(invalidYamlFiles []*extractor.InvalidFile) []printer.Warning {
	var warnings []printer.Warning

	for _, invalidFile := range invalidYamlFiles {
		warnings = append(warnings, printer.Warning{
			Title:       printer.GetFileNameText(invalidFile.Path),
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
			Title:       printer.GetFileNameText(invalidFile.Path),
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

func parseToPrinterWarnings(results *EvaluationResults, invalidYamlFiles []*extractor.InvalidFile, invalidK8sFiles []*extractor.InvalidFile, pwd string, k8sVersion string, k8sValidationWarnings validation.K8sValidationWarningPerValidFile, verbose bool) ([]printer.Warning, error) {
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
			var skippedRules = []printer.FailedRule{}

			rulesUniqueNames := []string{}
			for rulesUniqueName := range rules {
				rulesUniqueNames = append(rulesUniqueNames, rulesUniqueName)
			}

			sort.Strings(rulesUniqueNames)

			for _, ruleUniqueName := range rulesUniqueNames {
				rule := rules[ruleUniqueName]
				var fixLink string
				if verbose {
					fixLink = rule.DocumentationUrl
				}
				failedRule := printer.FailedRule{
					Name:               rule.Name,
					DocumentationUrl:   fixLink,
					Suggestion:         rule.MessageOnFailure,
					Occurrences:        rule.GetFailedOccurrencesCount(),
					OccurrencesDetails: []printer.OccurrenceDetails{},
				}
				skippedRule := failedRule

				hasSkippedOccurrences := false
				hasFailedOccurrences := false

				for _, occurrenceDetails := range rule.OccurrencesDetails {
					if occurrenceDetails.IsSkipped {
						hasSkippedOccurrences = true
						skippedRule.OccurrencesDetails = append(skippedRule.OccurrencesDetails, printer.OccurrenceDetails{
							MetadataName: occurrenceDetails.MetadataName,
							Kind:         occurrenceDetails.Kind,
							SkipMessage:  occurrenceDetails.SkipMessage,
						})
					} else {
						hasFailedOccurrences = true
						failedRule.OccurrencesDetails = append(
							failedRule.OccurrencesDetails,
							printer.OccurrenceDetails{
								MetadataName: occurrenceDetails.MetadataName,
								Kind:         occurrenceDetails.Kind,
							},
						)
					}

				}

				if hasSkippedOccurrences {
					skippedRules = append(skippedRules, skippedRule)
				}
				if hasFailedOccurrences {
					failedRules = append(failedRules, failedRule)
				}
			}

			title := filename
			relativePath, _ := filepath.Rel(pwd, filename)

			if relativePath != "" {
				title = printer.GetFileNameText(relativePath)
			}

			warnings = append(warnings, printer.Warning{
				Title:           title,
				FailedRules:     failedRules,
				SkippedRules:    skippedRules,
				InvalidYamlInfo: printer.InvalidYamlInfo{},
				InvalidK8sInfo: printer.InvalidK8sInfo{
					ValidationWarning: k8sValidationWarnings[filename].Warning,
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
			Text:  "Are you trying to test a raw Helm file? To run Datree with Helm - check out the helm plugin README:\nhttps://github.com/datreeio/helm-datree \n",
			Color: "cyan",
		})
	} else if IsKustomizationFile(invalidFile.Path) {
		extraMessages = append(extraMessages, printer.ExtraMessage{
			Text:  "Are you trying to test Kustomize files? To run Datree with Kustomize, use `datree kustomize test` command, or check out Kustomize support docs:\nhttps://hub.datree.io/kustomize-support \n",
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
	TotalSkippedRules
	TotalRulesFailed
)

func (t OutputTitle) String() string {
	return [...]string{
		"Configs tested against policy",
		"Total rules evaluated",
		"See all rules in policy",
		"Total rules passed",
		"Total rules skipped",
		"Total rules failed"}[t]
}

func buildEnabledRulesTitle(policyName string) string {
	var str strings.Builder
	fmt.Fprintf(&str, "Enabled rules in policy \"%s\"", policyName)
	return str.String()
}

func parseEvaluationResultsToSummary(results *EvaluationResults, evaluationSummary printer.EvaluationSummary, loginURL string, policyName string) printer.Summary {
	configsCount := evaluationSummary.ConfigsCount
	totalRulesEvaluated := 0
	totalFailedRules := 0
	totalSkippedRules := 0
	totalPassedRules := 0

	if results != nil {
		totalRulesEvaluated = evaluationSummary.RulesCount * results.Summary.FilesCount
		totalFailedRules = results.Summary.TotalFailedRules
		totalSkippedRules = results.Summary.TotalSkippedRules
		totalPassedRules = results.Summary.TotalPassedRules
	}

	plainRows := []printer.SummaryItem{
		{LeftCol: buildEnabledRulesTitle(policyName), RightCol: fmt.Sprint(evaluationSummary.RulesCount), RowIndex: 0},
		{LeftCol: EvaluatedConfigurations.String(), RightCol: fmt.Sprint(configsCount), RowIndex: 1},
		{LeftCol: TotalRulesEvaluated.String(), RightCol: fmt.Sprint(totalRulesEvaluated), RowIndex: 2},
		{LeftCol: SeeAll.String(), RightCol: loginURL, RowIndex: 6},
	}

	skipRow := printer.SummaryItem{LeftCol: TotalSkippedRules.String(), RightCol: fmt.Sprint(totalSkippedRules), RowIndex: 3}
	successRow := printer.SummaryItem{LeftCol: TotalRulesPassed.String(), RightCol: fmt.Sprint(totalPassedRules), RowIndex: 5}
	errorRow := printer.SummaryItem{LeftCol: TotalRulesFailed.String(), RightCol: fmt.Sprint(totalFailedRules), RowIndex: 4}

	summary := &printer.Summary{
		SkipRow:    skipRow,
		ErrorRow:   errorRow,
		SuccessRow: successRow,
		PlainRows:  plainRows,
	}
	return *summary
}
