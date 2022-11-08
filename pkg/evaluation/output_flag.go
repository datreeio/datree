package evaluation

import "strings"

var FormattedOutputOptions = []string{"yaml", "json", "xml", "JUnit", "sarif"}
var InteractiveOutputOptions = []string{"", "simple"}
var ValidOutputOptions = append(FormattedOutputOptions, InteractiveOutputOptions...)
var ExplicitOutputOptions = []string{"simple", "yaml", "json", "xml", "JUnit", "sarif"}

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
