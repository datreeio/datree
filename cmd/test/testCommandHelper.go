package test

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/pkg/extractor"
	"github.com/pkg/browser"
)

func createSpinner(text string, color string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = text
	s.Color(color)
	return s
}

func aggregateValidK8sFiles (validK8sFilesConfigurationsChan chan *extractor.FileConfigurations) (validK8sFilesConfigurations []*extractor.FileConfigurations){
	for fileConfigurations := range validK8sFilesConfigurationsChan {
		validK8sFilesConfigurations = append(validK8sFilesConfigurations, fileConfigurations)
	}
	return
}

func aggregateIgnoredYamlFiles(ignoredFilesChan chan *extractor.FileConfigurations) []extractor.FileConfigurations {
	var ignoredFiles []extractor.FileConfigurations

	for ignoredFile := range ignoredFilesChan {
		ignoredFiles = append(ignoredFiles, *ignoredFile)
	}
	return ignoredFiles
}

func countConfigurations(filesConfigurations []*extractor.FileConfigurations) int {
	totalConfigs := 0

	for _, fileConfiguration := range filesConfigurations {
		totalConfigs += len(fileConfiguration.Configurations)
	}

	return totalConfigs
}

func openBrowser(url string) {
	fmt.Printf("Opening %s in your browser.\n", url)
	browser.OpenURL(url)
}
