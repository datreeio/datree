package test

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl/validation"
)

func createSpinner(text string, color string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = text
	s.Color(color)
	return s
}

func aggregateInvalidFiles(invalidFilesChan chan *validation.InvalidFile) []*validation.InvalidFile {
	var invalidFiles []*validation.InvalidFile
	for invalidFile := range invalidFilesChan {
		invalidFiles = append(invalidFiles, invalidFile)
	}
	return invalidFiles
}
