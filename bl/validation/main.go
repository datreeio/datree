package validation

import (
	"fmt"
	"strings"
)

type InvalidK8sSchemaError struct {
	ErrorMessage string
}

func (e *InvalidK8sSchemaError) Error() string {
	return fmt.Sprintf("k8s schema validation error: %s\n%s", e.ErrorMessage, e.usageSuggestion())
}

func (e *InvalidK8sSchemaError) usageSuggestion() string {
	if strings.HasPrefix(e.ErrorMessage, "could not find schema for ") {
		return "You can skip files with missing schemas instead of failing by using the `--ignore-missing-schemas` flag\n"
	}
	return ""
}
