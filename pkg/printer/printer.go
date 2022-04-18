package printer

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"io"

	"github.com/xeipuuv/gojsonschema"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var out io.Writer = color.Output

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

func (p *Printer) printYamlValidationWarning(warning Warning) {
	p.printInColor("[X] YAML validation\n", p.Theme.Colors.White)
	fmt.Fprintln(out)
	for _, validationError := range warning.InvalidYamlInfo.ValidationErrors {
		validationError := p.Theme.Colors.RedBold.Sprint(validationError.Error())
		fmt.Fprintf(out, "%v %v\n", p.Theme.Emoji.Error, validationError)
	}
	fmt.Fprintln(out)

	p.printInColor("[?] Kubernetes schema validation didn't run for this file\n", p.Theme.Colors.White)
	p.printSkippedPolicyCheck()
	fmt.Fprintln(out)
}

func (p *Printer) printK8sValidationError(warning Warning) {
	p.printPassedYamlValidation()
	p.printInColor("[X] Kubernetes schema validation\n", p.Theme.Colors.White)
	fmt.Fprintln(out)

	for _, validationError := range warning.InvalidK8sInfo.ValidationErrors {
		validationError := p.Theme.Colors.RedBold.Sprint(validationError.Error())
		fmt.Fprintf(out, "%v %v\n", p.Theme.Emoji.Error, validationError)
	}

	for _, extraMessage := range warning.ExtraMessages {
		p.PrintMessage(extraMessage.Text, extraMessage.Color)
	}

	fmt.Fprintln(out)

	p.printSkippedPolicyCheck()
	fmt.Fprintln(out)
}

func (p *Printer) printK8sValidationWarning(warning Warning) {
	fmt.Println("[?] Kubernetes schema validation")
	fmt.Fprintln(out)

	fmt.Println(warning.InvalidK8sInfo.ValidationWarning)
	fmt.Fprintln(out)
}

func (p *Printer) PrintYamlSchemaResults(result []jsonschema.Detailed, error error) {
	if result == nil {
		p.printInColor("INVALID FILE PATH\n", p.Theme.Colors.RedBold)
		return
	}
	if result != nil {
		p.printInColor("Input does NOT pass validation against schema\n", p.Theme.Colors.RedBold)
		var errorsAsString = ""
		for _, desc := range result {
			errorsAsString = errorsAsString + desc.InstanceLocation + " - " + desc.Error + "\n"
		}
		p.printInColor(errorsAsString, p.Theme.Colors.RedBold)
		return
	}
	if error == nil {
		p.printInColor("Input PASSES validation against schema\n", p.Theme.Colors.Green)
	} else {
		p.printInColor("The File Is Invalid\n", p.Theme.Colors.RedBold)
	}
}

func (p *Printer) PrintWarnings(warnings []Warning) {
	for _, warning := range warnings {
		p.printInColor(warning.Title, p.Theme.Colors.Yellow)
		fmt.Fprintln(out)

		if len(warning.InvalidYamlInfo.ValidationErrors) > 0 {
			p.printYamlValidationWarning(warning)
		} else if len(warning.InvalidK8sInfo.ValidationErrors) > 0 {
			p.printK8sValidationError(warning)
		} else {
			p.printPassedYamlValidation()

			if warning.InvalidK8sInfo.ValidationWarning != "" {
				p.printK8sValidationWarning(warning)
			} else {
				p.printInColor("[V] Kubernetes schema validation\n", p.Theme.Colors.Green)
			}

			fmt.Fprintln(out)
			p.printInColor("[X] Policy check\n", p.Theme.Colors.White)
			fmt.Fprintln(out)

			if len(warning.SkippedRules) > 0 {
				fmt.Fprintf(out, "%v", p.Theme.Colors.CyanBold.Sprintf("SKIPPED")+"\n\n")
			}

			for _, skippedRule := range warning.SkippedRules {
				ruleName := p.Theme.Colors.CyanBold.Sprint(skippedRule.Name)

				fmt.Fprintf(out, "%v %v\n", p.Theme.Emoji.Skip, ruleName)

				if skippedRule.DocumentationUrl != "" {
					howToFix := p.Theme.Colors.Cyan.Sprint(skippedRule.DocumentationUrl)
					fmt.Fprintf(out, "    How to fix: %v\n", howToFix)
				}

				for _, occurrenceDetails := range skippedRule.OccurrencesDetails {
					fmt.Fprintf(out, "    — metadata.name: %v (kind: %v)\n", p.getStringOrNotAvailable(occurrenceDetails.MetadataName), p.getStringOrNotAvailable(occurrenceDetails.Kind))
					m := p.Theme.Colors.White.Sprint(occurrenceDetails.SkipMessage)
					fmt.Fprintf(out, "%v %v\n", p.Theme.Emoji.Suggestion, m)
				}

				fmt.Fprintln(out)
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

				fmt.Fprintf(out, "%v %v %v\n", p.Theme.Emoji.Error, ruleName, occurrences)

				if failedRule.DocumentationUrl != "" {
					howToFix := p.Theme.Colors.Cyan.Sprint(failedRule.DocumentationUrl)
					fmt.Fprintf(out, "    How to fix: %v\n", howToFix)
				}

				for _, occurrenceDetails := range failedRule.OccurrencesDetails {
					fmt.Fprintf(out, "    — metadata.name: %v (kind: %v)\n", p.getStringOrNotAvailable(occurrenceDetails.MetadataName), p.getStringOrNotAvailable(occurrenceDetails.Kind))
				}
				fmt.Fprintf(out, "%v %v\n", p.Theme.Emoji.Suggestion, failedRule.Suggestion)

				fmt.Fprintln(out)
			}
		}
	}

	fmt.Fprintln(out)
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

func (p *Printer) PrintEvaluationSummary(summary EvaluationSummary, k8sVersion string) {
	p.printInColor("(Summary)\n", p.Theme.Colors.White)
	fmt.Fprintln(out)

	fmt.Fprintf(out, "- Passing YAML validation: %v/%v\n", summary.PassedYamlValidationCount, summary.FilesCount)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing Kubernetes (%s) schema validation: %s\n", k8sVersion, summary.K8sValidation)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing policy check: %v/%v\n", summary.PassedPolicyCheckCount, summary.FilesCount)
	fmt.Fprintln(out)
}

func (p *Printer) PrintSummaryTable(summary Summary) {
	summaryTable := tablewriter.NewWriter(out)
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
	summaryTable.Rich(skipRow, []tablewriter.Colors{{int(p.Theme.ColorsAttributes.Cyan)}, {int(p.Theme.ColorsAttributes.Cyan)}})
	rowIndex++

	errorRow := []string{summary.ErrorRow.LeftCol, summary.ErrorRow.RightCol}
	summaryTable.Rich(errorRow, []tablewriter.Colors{{int(p.Theme.ColorsAttributes.Red)}, {int(p.Theme.ColorsAttributes.Red)}})
	rowIndex++

	successRow := []string{summary.SuccessRow.LeftCol, summary.SuccessRow.RightCol}
	summaryTable.Rich(successRow, []tablewriter.Colors{{int(p.Theme.ColorsAttributes.Green)}, {int(p.Theme.ColorsAttributes.Green)}})
	rowIndex++

	for plainRowsIndex < len(summary.PlainRows) && summary.PlainRows[plainRowsIndex].RowIndex >= rowIndex {
		summaryTable.Append([]string{summary.PlainRows[plainRowsIndex].LeftCol, summary.PlainRows[plainRowsIndex].RightCol})
		rowIndex++
	}

	summaryTable.Render()
}

func (p *Printer) printInColor(title string, color *color.Color) {
	colorPrintFn := color.FprintfFunc()
	colorPrintFn(out, title)
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

func (p *Printer) printPassedYamlValidation() {
	p.printInColor("[V] YAML validation\n", p.Theme.Colors.Green)
}

func (p *Printer) printSkippedPolicyCheck() {
	p.printInColor("[?] Policy check didn't run for this file\n", p.Theme.Colors.White)
}

func (p *Printer) getStringOrNotAvailable(str string) string {
	if str == "" {
		return "N/A"
	} else {
		return str
	}
}
