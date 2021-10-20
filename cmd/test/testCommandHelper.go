package test

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkg/browser"
)

func createSpinner(text string, color string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = text
	s.Color(color)
	return s
}

func openBrowser(url string) {
	fmt.Printf("Opening %s in your browser.\n", url)
	browser.OpenURL(url)
}
