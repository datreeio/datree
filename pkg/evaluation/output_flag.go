package evaluation

import "strings"

var FormattedOutputOptions = []string{"yaml", "json", "xml", "JUnit"}
var InteractiveOutputOptions = []string{"", "simple"}
var ValidOutputOptions = append(FormattedOutputOptions, InteractiveOutputOptions...)
var ExplicitOutputOptions = []string{"simple", "yaml", "json", "xml", "JUnit"}

func IsValidOutputOption(option string) bool {
	for _, validOption := range ValidOutputOptions {
		if option == validOption {
			return true
		}
	}
	return false
}

func IsFormattedOutputOption(option string) bool {
	for _, formattedOption := range FormattedOutputOptions {
		if option == formattedOption {
			return true
		}
	}
	return false
}

func OutputFormats() string {
	return strings.Join(ExplicitOutputOptions, ", ")
}
