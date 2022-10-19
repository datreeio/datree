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
		return "To troubleshoot: refer to our docs [https://hub.datree.io/troubleshooting#schema-validation-failure]\nTo ignore this failure: use the CLI flag `--ignore-missing-schemas`\n"
	}
	return ""
}
