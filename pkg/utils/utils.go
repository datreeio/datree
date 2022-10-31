package utils

import (
	"fmt"
	"strings"

	"github.com/pkg/browser"
)

func Example(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalize{s}.trim().indent().string
}

const Indentation = `  `

type normalize struct {
	string
}

func (s normalize) trim() normalize {
	s.string = strings.TrimSpace(s.string)
	return s
}

func (s normalize) indent() normalize {
	indentedLines := []string{}
	for _, line := range strings.Split(s.string, "\n") {
		trimmed := strings.TrimSpace(line)
		indented := Indentation + trimmed
		indentedLines = append(indentedLines, indented)
	}
	s.string = strings.Join(indentedLines, "\n")
	return s
}

func ValidateStdinPathArgument(paths []string) error {
	if len(paths) < 1 {
		return fmt.Errorf("requires at least 1 arg")
	}

	if paths[0] == "-" {
		if len(paths) > 1 {
			return fmt.Errorf(fmt.Sprintf("unexpected args: [%s]", strings.Join(paths[1:], ",")))
		}
	}

	return nil
}

type URLOpener interface {
	OpenURL(url string) error
}

type OpenBrowserContext struct {
	UrlOpener URLOpener
}

func (o *OpenBrowserContext) OpenBrowser(url string) error {
	fmt.Printf("Opening %s in your browser.\n", url)
	err := browser.OpenURL(url)
	if err != nil {
		return fmt.Errorf("error opening url: %v", err)
	}
	return nil
}
