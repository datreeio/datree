package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var out io.Writer = os.Stdout

type Printer struct {
	theme *theme
}

func CreateNewPrinter() *Printer {
	theme := createTheme()
	return &Printer{
		theme: theme,
	}
}

type WarningInfo struct {
	Caption            string
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
type Warning struct {
	Title           string
	Details         []WarningInfo
	InvalidYamlInfo InvalidYamlInfo
	InvalidK8sInfo  InvalidK8sInfo
}

func (p *Printer) SimplePrintWarnings(warnings []Warning) {
	for _, warning := range warnings {
		fmt.Fprintf(out, warning.Title)
		fmt.Fprintln(out)

		if len(warning.InvalidYamlInfo.ValidationErrors) > 0 {
			fmt.Fprintf(out, "[X] YAML validation\n")
			fmt.Fprintln(out)
			for _, validationError := range warning.InvalidYamlInfo.ValidationErrors {
				fmt.Fprintf(out, "%v %v\n", "[X]", validationError.Error())
			}
			fmt.Fprintln(out)

			fmt.Fprintf(out, "[?] Kubernetes schema validation didn't run for this file\n")
			fmt.Fprintf(out, "[?] Policy check didn't run for this file\n")
			fmt.Fprintln(out)
		} else if len(warning.InvalidK8sInfo.ValidationErrors) > 0 {
			fmt.Fprintf(out, "[V] YAML validation\n")
			fmt.Fprintf(out, "[X] Kubernetes schema validation\n")
			fmt.Fprintln(out)

			for _, validationError := range warning.InvalidK8sInfo.ValidationErrors {
				fmt.Fprintf(out, "%v %v\n", "[X]", validationError.Error())
			}
			fmt.Fprintln(out)
			fmt.Fprintf(out, "[?] Policy check didn't run for this file\n")
			fmt.Fprintln(out)
		} else {
			fmt.Fprintf(out, "[V] YAML validation\n")
			fmt.Fprintf(out, "[V] Kubernetes schema validation\n")

			fmt.Fprintln(out)
			fmt.Fprintf(out, "[X] Policy check\n")
			fmt.Fprintln(out)

			for _, details := range warning.Details {
				var occurrencesPostfix string
				if details.Occurrences > 1 {
					occurrencesPostfix = "s"
				} else {
					occurrencesPostfix = ""
				}
				formattedOccurrences := fmt.Sprintf(" [%d occurrence%v]", details.Occurrences, occurrencesPostfix)

				fmt.Fprintf(out, "%v %v %v\n", "[X]", details.Caption, formattedOccurrences)
				for _, occurrenceDetails := range details.OccurrencesDetails {
					fmt.Fprintf(out, "    — metadata.name: %v (kind: %v)\n", p.getStringOrNotAvailable(occurrenceDetails.MetadataName), p.getStringOrNotAvailable(occurrenceDetails.Kind))
				}
				fmt.Fprintf(out, "%v %v\n", "[*]", details.Suggestion)

				fmt.Fprintln(out)
			}
		}
	}

	fmt.Fprintln(out)
}

func (p *Printer) PrintWarnings(warnings []Warning) {
	for _, warning := range warnings {
		p.printInColor(warning.Title, p.theme.Colors.Yellow)
		fmt.Fprintln(out)

		if len(warning.InvalidYamlInfo.ValidationErrors) > 0 {
			p.printInColor("[X] YAML validation\n", p.theme.Colors.White)
			fmt.Fprintln(out)
			for _, validationError := range warning.InvalidYamlInfo.ValidationErrors {
				validationError := p.theme.Colors.Red.Sprint(validationError.Error())
				fmt.Fprintf(out, "%v %v\n", p.theme.Emoji.Error, validationError)
			}
			fmt.Fprintln(out)

			p.printInColor("[?] Kubernetes schema validation didn't run for this file\n", p.theme.Colors.White)
			p.printSkippedPolicyCheck()
			fmt.Fprintln(out)
		} else if len(warning.InvalidK8sInfo.ValidationErrors) > 0 {
			p.printPassedYamlValidation()
			p.printInColor("[X] Kubernetes schema validation\n", p.theme.Colors.White)
			fmt.Fprintln(out)

			for _, validationError := range warning.InvalidK8sInfo.ValidationErrors {
				validationError := p.theme.Colors.Red.Sprint(validationError.Error())
				fmt.Fprintf(out, "%v %v\n", p.theme.Emoji.Error, validationError)
			}
			fmt.Fprintln(out)
			p.printSkippedPolicyCheck()
			fmt.Fprintln(out)
		} else {
			p.printPassedYamlValidation()
			p.printInColor("[V] Kubernetes schema validation\n", p.theme.Colors.Green)

			fmt.Fprintln(out)
			p.printInColor("[X] Policy check\n", p.theme.Colors.White)
			fmt.Fprintln(out)

			for _, details := range warning.Details {
				var occurrencesPostfix string
				if details.Occurrences > 1 {
					occurrencesPostfix = "s"
				} else {
					occurrencesPostfix = ""
				}
				formattedOccurrences := fmt.Sprintf(" [%d occurrence%v]", details.Occurrences, occurrencesPostfix)
				occurrences := p.theme.Colors.White.Sprintf(formattedOccurrences)

				caption := p.theme.Colors.Red.Sprint(details.Caption)

				fmt.Fprintf(out, "%v %v %v\n", p.theme.Emoji.Error, caption, occurrences)
				for _, occurrenceDetails := range details.OccurrencesDetails {
					fmt.Fprintf(out, "    — metadata.name: %v (kind: %v)\n", p.getStringOrNotAvailable(occurrenceDetails.MetadataName), p.getStringOrNotAvailable(occurrenceDetails.Kind))
				}
				fmt.Fprintf(out, "%v %v\n", p.theme.Emoji.Suggestion, details.Suggestion)

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

func (p *Printer) SimplePrintEvaluationSummary(summary EvaluationSummary, k8sVersion string) {
	fmt.Fprintf(out, "(Summary)\n")
	fmt.Fprintln(out)

	fmt.Fprintf(out, "- Passing YAML validation: %v/%v\n", summary.PassedYamlValidationCount, summary.FilesCount)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing Kubernetes (%s) schema validation: %v/%v\n", k8sVersion, summary.PassedK8sValidationCount, summary.FilesCount)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing policy check: %v/%v\n", summary.PassedPolicyCheckCount, summary.FilesCount)
	fmt.Fprintln(out)
}

func (p *Printer) PrintEvaluationSummary(summary EvaluationSummary, k8sVersion string) {
	p.printInColor("(Summary)\n", p.theme.Colors.White)
	fmt.Fprintln(out)

	fmt.Fprintf(out, "- Passing YAML validation: %v/%v\n", summary.PassedYamlValidationCount, summary.FilesCount)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing Kubernetes (%s) schema validation: %v/%v\n", k8sVersion, summary.PassedK8sValidationCount, summary.FilesCount)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- Passing policy check: %v/%v\n", summary.PassedPolicyCheckCount, summary.FilesCount)
	fmt.Fprintln(out)
}

func (p *Printer) SimplePrintSummaryTable(summary Summary) {
	summaryTable := tablewriter.NewWriter(os.Stdout)
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
	summaryTable.Append(errorRow)
	rowIndex++

	successRow := []string{summary.SuccessRow.LeftCol, summary.SuccessRow.RightCol}
	summaryTable.Append(successRow)
	rowIndex++

	for plainRowsIndex < len(summary.PlainRows) && summary.PlainRows[plainRowsIndex].RowIndex >= rowIndex {
		summaryTable.Append([]string{summary.PlainRows[plainRowsIndex].LeftCol, summary.PlainRows[plainRowsIndex].RightCol})
		rowIndex++
	}

	summaryTable.Render()
}

func (p *Printer) PrintSummaryTable(summary Summary) {
	summaryTable := tablewriter.NewWriter(os.Stdout)
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
	summaryTable.Rich(errorRow, []tablewriter.Colors{{tablewriter.FgHiRedColor}, {tablewriter.FgHiRedColor}})
	rowIndex++

	successRow := []string{summary.SuccessRow.LeftCol, summary.SuccessRow.RightCol}
	summaryTable.Rich(successRow, []tablewriter.Colors{{tablewriter.Normal, tablewriter.FgGreenColor}, {tablewriter.Normal, tablewriter.FgGreenColor}})
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
		return p.theme.Colors.Error
	case "red":
		return p.theme.Colors.Red
	case "yellow":
		return p.theme.Colors.Yellow
	case "green":
		return p.theme.Colors.Green
	default:
		return p.theme.Colors.White
	}
}

func (p *Printer) PrintMessage(messageText string, messageColor string) {
	colorPrintFn := p.createNewColor(messageColor)
	p.printInColor(messageText, colorPrintFn)
}

func (p *Printer) printPassedYamlValidation() {
	p.printInColor("[V] YAML validation\n", p.theme.Colors.Green)
}

func (p *Printer) printSkippedPolicyCheck() {
	p.printInColor("[?] Policy check didn't run for this file\n", p.theme.Colors.White)
}

func (p *Printer) getStringOrNotAvailable(str string) string {
	if str == "" {
		return "N/A"
	} else {
		return str
	}
}
