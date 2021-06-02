package test

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"
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

func countConfigurations(filesConfigurations []*extractor.FileConfigurations) int {
	totalConfigs := 0

	for _, fileConfiguration := range filesConfigurations {
		totalConfigs += len(fileConfiguration.Configurations)
	}

	return totalConfigs
}
