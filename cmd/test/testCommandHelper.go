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

func aggregateIgnoredYamlFiles(ignoredFilesChan chan *extractor.FileConfigurations) []extractor.FileConfigurations {
	var ignoredFiles []extractor.FileConfigurations

	for ignoredFile := range ignoredFilesChan {
		ignoredFiles = append(ignoredFiles, *ignoredFile)
	}
	return ignoredFiles
}

func openBrowser(url string) {
	fmt.Printf("Opening %s in your browser.\n", url)
	browser.OpenURL(url)
}
