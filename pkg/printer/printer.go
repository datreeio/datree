package printer

import (
	"fmt"
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
	OccurrencesDetails []OccurrenceDetails
}

type OccurrenceDetails struct {
	MetadataName string
	Kind         string
}

type InvalidYamlInfo struct {
	ValidationErrors []error
}
type InvalidK8sInfo struct {
	ValidationErrors []error
	K8sVersion       string
}

type ExtraMessage struct {
	Text  string
	Color string
}

type Warning struct {
	Title           string
	FailedRules     []FailedRule
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

func (p *Printer) printK8sValidationWarning(warning Warning) {
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

func (p *Printer) PrintYamlSchemaResults(result *JSONSchemaValidatorResults, error error) {
	if result == nil {
		p.printInColor("INVALID FILE PATH\n", p.Theme.Colors.RedBold)
		return
	}
	if result.Errors() != nil {
		p.printInColor("Input does NOT pass validation against schema\n", p.Theme.Colors.RedBold)
		var errorsAsString = ""
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
			errorsAsString = errorsAsString + desc.String() + "\n"
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
			p.printK8sValidationWarning(warning)
		} else {
			p.printPassedYamlValidation()
			p.printInColor("[V] Kubernetes schema validation\n", p.Theme.Colors.Green)

			fmt.Fprintln(out)
			p.printInColor("[X] Policy check\n", p.Theme.Colors.White)
			fmt.Fprintln(out)

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
	ErrorRow   SummaryItem
	SuccessRow SummaryItem
}

type EvaluationSummary struct {
	ConfigsCount              int
	RulesCount                int
	FilesCount                int
	PassedYamlValidationCount int
	PassedK8sValidationCount  int
	PassedPolicyCheckCount    int
}

func (p *Printer) PrintEvaluationSummary(summary EvaluationSummary, k8sVersion string) {
	p.printInColor("(Summary)\n", p.Theme.Colors.White)
	fmt.Fprintln(out)

	fmt.Fprintf(out, "- Passing YAML validation: %v/%v\n", summary.PassedYamlValidationCount, summary.FilesCount)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing Kubernetes (%s) schema validation: %v/%v\n", k8sVersion, summary.PassedK8sValidationCount, summary.FilesCount)
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
