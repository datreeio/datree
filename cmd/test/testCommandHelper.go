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

func aggregateInvalidK8sFiles(invalidFilesChan chan *validation.InvalidK8sFile) []*validation.InvalidK8sFile {
	var invalidFiles []*validation.InvalidK8sFile
	for invalidFile := range invalidFilesChan {
		invalidFiles = append(invalidFiles, invalidFile)
	}
	return invalidFiles
}
func aggregateInvalidYamlFiles(invalidFilesChan chan *validation.InvalidYamlFile) []*validation.InvalidYamlFile {
	var invalidFiles []*validation.InvalidYamlFile
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
