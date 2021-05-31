package test

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl/evaluation"
	"github.com/datreeio/datree/bl/validation"
)

func createSpinner(text string, color string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = text
	s.Color(color)
	return s
}

func printEvaluationErrors(errors []*evaluation.Error) {
	fmt.Println("The following files failed:")
	for _, fileError := range errors {
		fmt.Printf("\n\tFilename: %s\n\tError: %s\n---------------------", fileError.Filename, fileError.Message)
	}
	fmt.Println()
}

func handleInvalidFiles(invalidFilesChan chan *validation.InvalidFile) []*validation.InvalidFile {
	time.Sleep(time.Millisecond * 10000)
	var invalidFiles []*validation.InvalidFile
	for invalidFile := range invalidFilesChan {
		invalidFiles = append(invalidFiles, invalidFile)
	}
	return invalidFiles
}
