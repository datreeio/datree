package validation

import "fmt"

type InvalidK8sSchemaError struct {
	ErrorMessage string
}

func (e *InvalidK8sSchemaError) Error() string {
	return fmt.Sprintf("k8s schema validation error: %s\n", e.ErrorMessage)
}

type InvalidYamlError struct {
	ErrorMessage string
}

func (e *InvalidYamlError) Error() string {
	return fmt.Sprintf("yaml validation error: %s\n", e.ErrorMessage)
}
