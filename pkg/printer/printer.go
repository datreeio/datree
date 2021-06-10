package printer

import (
	"fmt"
	"os"

	"github.com/datreeio/datree/bl/validation"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

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
	Caption     string
	Occurrences int
	Suggestion  string
}

type ValidationInfo struct {
	IsValid          bool
	ValidationStatus validation.ValidationStatus
	ValidationErrors []error
	K8sVersion       string
}
type Warning struct {
	Title          string
	Details        []WarningInfo
	ValidationInfo ValidationInfo
}

func (p *Printer) PrintWarnings(warnings []Warning) {
	for _, warning := range warnings {
		p.printInColor(warning.Title, p.theme.Colors.Yellow)
		fmt.Println()

		if !warning.ValidationInfo.IsValid {
			if warning.ValidationInfo.ValidationStatus == validation.InvalidYamlFile {
				p.printInColor("[X] YAML validation\n", p.theme.Colors.White)
				fmt.Println()
				for _, validationError := range warning.ValidationInfo.ValidationErrors {
					validationError := p.theme.Colors.Red.Sprint(validationError.Error())
					fmt.Printf("%v %v\n", p.theme.Emoji.Error, validationError)
				}
				fmt.Println()

				p.printInColor("[?] Kubernetes schema validation\n", p.theme.Colors.White)
				p.printInColor("[?] Policy check didn’t run for this file\n", p.theme.Colors.White)

				fmt.Println()

			} else if warning.ValidationInfo.ValidationStatus == validation.InvalidK8sFile {
				p.printInColor("[V] YAML validation\n", p.theme.Colors.Green)
				p.printInColor("[X] Kubernetes schema validation\n", p.theme.Colors.White)
				fmt.Println()

				for _, validationError := range warning.ValidationInfo.ValidationErrors {
					validationError := p.theme.Colors.Red.Sprint(validationError.Error())
					fmt.Printf("%v %v\n", p.theme.Emoji.Error, validationError)
				}
				fmt.Println()
				p.printInColor("[?] Policy check didn’t run for this file\n", p.theme.Colors.White)
				fmt.Println()
			}
		} else {
			p.printInColor("[V] YAML validation\n", p.theme.Colors.Green)
			p.printInColor("[V] Kubernetes schema validation\n", p.theme.Colors.Green)

			fmt.Println()
			p.printInColor("[X] Policy check\n", p.theme.Colors.White)
			fmt.Println()

			for _, details := range warning.Details {
				formattedOccurrences := fmt.Sprintf(" [%d occurrences]", details.Occurrences)
				occurrences := p.theme.Colors.White.Sprintf(formattedOccurrences)

				caption := p.theme.Colors.Red.Sprint(details.Caption)

				fmt.Printf("%v %v %v\n", p.theme.Emoji.Error, caption, occurrences)
				fmt.Printf("%v %v\n", p.theme.Emoji.Suggestion, details.Suggestion)

				fmt.Println()
			}
		}
	}

	fmt.Println()
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
	p.printInColor("(Summary)\n", p.theme.Colors.White)
	fmt.Println()

	fmt.Printf("- Passing YAML validation: %v/%v\n", summary.PassedYamlValidationCount, summary.FilesCount)
	fmt.Println()
	fmt.Printf("- Passing Kubernetes (%s) schema validation: %v/%v\n", k8sVersion, summary.PassedK8sValidationCount, summary.FilesCount)
	fmt.Println()
	fmt.Printf("- Passing policy check: %v/%v\n", summary.PassedPolicyCheckCount, summary.FilesCount)
	fmt.Println()
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
	colorPrintFn := color.PrintfFunc()
	colorPrintFn(title)
}

func (p *Printer) createNewColor(clr string) *color.Color {
	switch clr {
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
