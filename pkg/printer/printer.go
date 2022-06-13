package printer

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/xeipuuv/gojsonschema"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var out = color.Output

type Printer struct {
	Theme *Theme
}

func CreateNewPrinter() *Printer {
	theme := createDefaultTheme()
	return &Printer{
		Theme: theme,
	}
}

type FailedRule struct {
	Name               string
	Occurrences        int
	Suggestion         string
	DocumentationUrl   string
	OccurrencesDetails []OccurrenceDetails
}

type OccurrenceDetails struct {
	MetadataName string
	Kind         string
	SkipMessage  string
}

type InvalidYamlInfo struct {
	ValidationErrors []error
}
type InvalidK8sInfo struct {
	ValidationErrors  []error
	ValidationWarning string
	K8sVersion        string
}

type ExtraMessage struct {
	Text  string
	Color string
}

type Warning struct {
	Title           string
	FailedRules     []FailedRule
	SkippedRules    []FailedRule
	InvalidYamlInfo InvalidYamlInfo
	InvalidK8sInfo  InvalidK8sInfo
	ExtraMessages   []ExtraMessage
}
type JSONSchemaValidatorResults = gojsonschema.Result

func (p *Printer) SetTheme(theme *Theme) {
	p.Theme = theme
}

func (p *Printer) getYamlValidationWarningText(warning Warning) string {
	sb := strings.Builder{}
	sb.WriteString(p.GetYamlValidationErrorsText(warning.InvalidYamlInfo.ValidationErrors))
	sb.WriteString(p.GetTextInColor("[?] Kubernetes schema validation didn't run for this file\n", p.Theme.Colors.White))
	sb.WriteString(p.getSkippedPolicyCheckText())
	sb.WriteString("\n")
	return sb.String()
}

func (p *Printer) GetYamlValidationErrorsText(yamlValidationErrors []error) string {
	sb := strings.Builder{}
	sb.WriteString(p.GetTextInColor("[X] YAML validation\n\n", p.Theme.Colors.White))
	for _, validationError := range yamlValidationErrors {
		validationError := p.Theme.Colors.RedBold.Sprint(validationError.Error())
		sb.WriteString(fmt.Sprintf("%v %v\n", p.Theme.Emoji.Error, validationError))
	}
	sb.WriteString("\n")
	return sb.String()
}

func (p *Printer) getK8sValidationErrorText(warning Warning) string {
	sb := strings.Builder{}
	sb.WriteString(p.getPassedYamlValidationText())
	sb.WriteString(p.GetTextInColor("[X] Kubernetes schema validation\n\n", p.Theme.Colors.White))

	for _, validationError := range warning.InvalidK8sInfo.ValidationErrors {
		validationError := p.Theme.Colors.RedBold.Sprint(validationError.Error())
		sb.WriteString(fmt.Sprintf("%v %v\n", p.Theme.Emoji.Error, validationError))
	}

	for _, extraMessage := range warning.ExtraMessages {
		sb.WriteString(p.GetTextInColor(extraMessage.Text, p.createNewColor(extraMessage.Color)))
	}

	sb.WriteString("\n")
	sb.WriteString(p.getSkippedPolicyCheckText())
	sb.WriteString("\n")
	return sb.String()
}

func (p *Printer) getK8sValidationWarningText(warning Warning) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	sb.WriteString("[?] Kubernetes schema validation")
	sb.WriteString("\n")
	sb.WriteString(warning.InvalidK8sInfo.ValidationWarning)
	sb.WriteString("\n")
	return sb.String()
}

func (p *Printer) PrintYamlSchemaResults(errorsResult []jsonschema.Detailed, error error) {
	fmt.Print(p.getYamlSchemaResultsText(errorsResult, error))
}

func (p *Printer) getYamlSchemaResultsText(errorsResult []jsonschema.Detailed, error error) string {
	sb := strings.Builder{}
	if errorsResult != nil {
		sb.WriteString(p.GetTextInColor("Input does NOT pass validation against schema\n", p.Theme.Colors.RedBold))
		var errorsAsString = ""
		for _, desc := range errorsResult {
			errorsAsString = errorsAsString + desc.InstanceLocation + " - " + desc.Error + "\n"
		}
		sb.WriteString(p.GetTextInColor(errorsAsString, p.Theme.Colors.RedBold))
		return sb.String()
	}
	if error == nil {
		sb.WriteString(p.GetTextInColor("Input PASSES validation against schema\n", p.Theme.Colors.Green))
	} else {
		sb.WriteString(p.GetTextInColor("The File Is Invalid\n", p.Theme.Colors.RedBold))
	}
	return sb.String()
}

func (p *Printer) GetWarningsText(warnings []Warning) string {
	var sb strings.Builder
	for _, warning := range warnings {
		sb.WriteString(p.GetTitleText(warning.Title))

		if len(warning.InvalidYamlInfo.ValidationErrors) > 0 {
			sb.WriteString(p.getYamlValidationWarningText(warning))
		} else if len(warning.InvalidK8sInfo.ValidationErrors) > 0 {
			sb.WriteString(p.getK8sValidationErrorText(warning))
		} else {
			sb.WriteString(p.getPassedYamlValidationText())

			if warning.InvalidK8sInfo.ValidationWarning != "" {
				sb.WriteString(p.getK8sValidationWarningText(warning))
			} else {
				sb.WriteString(p.GetTextInColor("[V] Kubernetes schema validation\n", p.Theme.Colors.Green))
			}

			sb.WriteString("\n")
			sb.WriteString(p.GetTextInColor("[X] Policy check\n", p.Theme.Colors.White))
			sb.WriteString("\n")

			if len(warning.SkippedRules) > 0 {
				sb.WriteString(fmt.Sprintf("%v", p.Theme.Colors.CyanBold.Sprintf("SKIPPED")+"\n\n"))
			}

			for _, skippedRule := range warning.SkippedRules {
				ruleName := p.Theme.Colors.CyanBold.Sprint(skippedRule.Name)

				sb.WriteString(fmt.Sprintf("%v %v\n", p.Theme.Emoji.Skip, ruleName))

				if skippedRule.DocumentationUrl != "" {
					howToFix := p.Theme.Colors.Cyan.Sprint(skippedRule.DocumentationUrl)
					sb.WriteString(fmt.Sprintf("    How to fix: %v\n", howToFix))
				}

				for _, occurrenceDetails := range skippedRule.OccurrencesDetails {
					sb.WriteString(fmt.Sprintf("    - metadata.name: %v (kind: %v)\n", p.getStringOrNotAvailableText(occurrenceDetails.MetadataName), p.getStringOrNotAvailableText(occurrenceDetails.Kind)))
					m := p.Theme.Colors.White.Sprint(occurrenceDetails.SkipMessage)
					sb.WriteString(fmt.Sprintf("%v %v\n", p.Theme.Emoji.Suggestion, m))
				}

				sb.WriteString("\n")
			}

			for _, failedRule := range warning.FailedRules {
				var occurrencesPostfix string

				if failedRule.Occurrences > 1 {
					occurrencesPostfix = "s"
				} else {
					occurrencesPostfix = ""
				}
				formattedOccurrences := fmt.Sprintf(" [%d occurrence%v]", failedRule.Occurrences, occurrencesPostfix)
				occurrences := p.Theme.Colors.White.Sprintf(formattedOccurrences)

				ruleName := p.Theme.Colors.RedBold.Sprint(failedRule.Name)

				sb.WriteString(fmt.Sprintf("%v %v %v\n", p.Theme.Emoji.Error, ruleName, occurrences))

				if failedRule.DocumentationUrl != "" {
					howToFix := p.Theme.Colors.Cyan.Sprint(failedRule.DocumentationUrl)
					sb.WriteString(fmt.Sprintf("    How to fix: %v\n", howToFix))
				}

				for _, occurrenceDetails := range failedRule.OccurrencesDetails {
					sb.WriteString(fmt.Sprintf("    - metadata.name: %v (kind: %v)\n", p.getStringOrNotAvailableText(occurrenceDetails.MetadataName), p.getStringOrNotAvailableText(occurrenceDetails.Kind)))
				}
				sb.WriteString(fmt.Sprintf("%v %v\n", p.Theme.Emoji.Suggestion, failedRule.Suggestion))

				sb.WriteString("\n")
			}
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

type SummaryItem struct {
	RightCol string
	LeftCol  string
	RowIndex int
}

type Summary struct {
	PlainRows  []SummaryItem
	SkipRow    SummaryItem
	ErrorRow   SummaryItem
	SuccessRow SummaryItem
}

type EvaluationSummary struct {
	ConfigsCount              int
	RulesCount                int
	FilesCount                int
	PassedYamlValidationCount int
	K8sValidation             string
	PassedPolicyCheckCount    int
}

func (p *Printer) GetTitleText(title string) string {
	return p.GetTextInColor(title, p.Theme.Colors.Yellow)
}

func GetFileNameText(title string) string {
	return fmt.Sprintf(">>  File: %s\n\n", title)
}

func (p *Printer) GetEvaluationSummaryText(summary EvaluationSummary, k8sVersion string) string {
	var sb strings.Builder
	sb.WriteString(p.GetTextInColor("(Summary)\n\n", p.Theme.Colors.White))

	sb.WriteString(p.GetYamlValidationSummaryText(summary.PassedYamlValidationCount, summary.FilesCount))

	sb.WriteString(fmt.Sprintf("- Passing Kubernetes (%s) schema validation: %s\n\n", k8sVersion, summary.K8sValidation))
	sb.WriteString(fmt.Sprintf("- Passing policy check: %v/%v\n\n", summary.PassedPolicyCheckCount, summary.FilesCount))
	return sb.String()
}

func (p *Printer) GetYamlValidationSummaryText(passedFiles int, allFiles int) string {
	return fmt.Sprintf("- Passing YAML validation: %v/%v\n\n", passedFiles, allFiles)
}

func (p *Printer) GetSummaryTableText(summary Summary) string {
	var sb strings.Builder
	summaryTable := tablewriter.NewWriter(&sb)
	summaryTable.SetAutoWrapText(false)
	summaryTable.SetAlignment(tablewriter.ALIGN_LEFT)

	rowIndex := 0
	plainRowsIndex := 0
	for i, item := range summary.PlainRows {
		if item.RowIndex == i {
			summaryTable.Append([]string{summary.PlainRows[i].LeftCol, summary.PlainRows[i].RightCol})
			plainRowsIndex++
			rowIndex++
		}
	}

	skipRow := []string{summary.SkipRow.LeftCol, summary.SkipRow.RightCol}
	errorRow := []string{summary.ErrorRow.LeftCol, summary.ErrorRow.RightCol}
	successRow := []string{summary.SuccessRow.LeftCol, summary.SuccessRow.RightCol}

	if p.Theme.Name == "Simple" {
		summaryTable.Append(skipRow)
		summaryTable.Append(errorRow)
		summaryTable.Append(successRow)
	} else {
		summaryTable.Rich(skipRow, []tablewriter.Colors{{int(p.Theme.ColorsAttributes.Cyan)}, {int(p.Theme.ColorsAttributes.Cyan)}})
		summaryTable.Rich(errorRow, []tablewriter.Colors{{int(p.Theme.ColorsAttributes.Red)}, {int(p.Theme.ColorsAttributes.Red)}})
		summaryTable.Rich(successRow, []tablewriter.Colors{{int(p.Theme.ColorsAttributes.Green)}, {int(p.Theme.ColorsAttributes.Green)}})
	}

	rowIndex = rowIndex + 3

	for plainRowsIndex < len(summary.PlainRows) && summary.PlainRows[plainRowsIndex].RowIndex >= rowIndex {
		summaryTable.Append([]string{summary.PlainRows[plainRowsIndex].LeftCol, summary.PlainRows[plainRowsIndex].RightCol})
		rowIndex++
	}

	summaryTable.Render()
	return sb.String()
}

func (p *Printer) printInColor(title string, color *color.Color) {
	colorPrintFn := color.FprintfFunc()
	colorPrintFn(out, title)
}

func (p *Printer) GetTextInColor(text string, color *color.Color) string {
	colorSprintFn := color.SprintfFunc()
	return colorSprintFn(text)
}

func (p *Printer) createNewColor(clr string) *color.Color {
	switch clr {
	case "error":
		return p.Theme.Colors.Error
	case "red":
		return p.Theme.Colors.RedBold
	case "yellow":
		return p.Theme.Colors.Yellow
	case "green":
		return p.Theme.Colors.Green
	case "cyan":
		return p.Theme.Colors.Cyan
	default:
		return p.Theme.Colors.White
	}
}

func (p *Printer) PrintMessage(messageText string, messageColor string) {
	colorPrintFn := p.createNewColor(messageColor)
	p.printInColor(messageText, colorPrintFn)
}

func (p *Printer) PrintPromptMessage(promptMessage string) {
	fmt.Fprint(out, color.HiCyanString("\n\n"+promptMessage+" (Y/n)\n"))
}

func (p *Printer) getPassedYamlValidationText() string {
	return p.GetTextInColor("[V] YAML validation\n", p.Theme.Colors.Green)
}

func (p *Printer) getSkippedPolicyCheckText() string {
	return p.GetTextInColor("[?] Policy check didn't run for this file\n", p.Theme.Colors.White)
}

func (p *Printer) getStringOrNotAvailableText(str string) string {
	if str == "" {
		return "N/A"
	} else {
		return str
	}
}
