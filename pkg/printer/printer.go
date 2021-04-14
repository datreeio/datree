package printer

import (
	"fmt"
	"os"

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

type Warning struct {
	Title   string
	Details []struct {
		Caption     string
		Occurrences int
		Suggestion  string
	}
}

func (p *Printer) PrintWarnings(warnings []Warning) {
	for _, warning := range warnings {
		p.printInColor(warning.Title, p.theme.Colors.Warning)

		fmt.Println()

		for _, d := range warning.Details {
			formattedOccurrences := fmt.Sprintf(" [%d occurrences]", d.Occurrences)
			occurrences := p.theme.Colors.Plain.Sprintf(formattedOccurrences)

			caption := p.theme.Colors.Error.Sprint(d.Caption)

			fmt.Printf("%v %v %v\n", p.theme.Emoji.Error, caption, occurrences)
			fmt.Printf("%v %v\n", p.theme.Emoji.Suggestion, d.Suggestion)

			fmt.Println()
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
	summaryTable.Rich(errorRow, []tablewriter.Colors{tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor}})
	rowIndex++

	successRow := []string{summary.SuccessRow.LeftCol, summary.SuccessRow.RightCol}
	summaryTable.Rich(successRow, []tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}})
	rowIndex++

	for plainRowsIndex < len(summary.PlainRows) && summary.PlainRows[plainRowsIndex].RowIndex >= rowIndex {
		summaryTable.Append([]string{summary.PlainRows[plainRowsIndex].LeftCol, summary.PlainRows[plainRowsIndex].RightCol})
		rowIndex++
	}

	summaryTable.Render()
}

func (p *Printer) printInColor(title string, color *color.Color) {
	warningColorPrintFn := color.PrintfFunc()
	warningColorPrintFn(title)
}
