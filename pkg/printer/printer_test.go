package printer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarnings(t *testing.T) {
	t.Run("Test Warning", func(t *testing.T) {

		out = new(bytes.Buffer)

		printer := CreateNewPrinter()

		warnings := []Warning{
			Warning{
				Title: "Failed with Occurrences",
				Details: []WarningInfo{
					WarningInfo{
						Caption:            "Caption",
						Occurrences:        1,
						Suggestion:         "Suggestion",
						OccurrencesDetails: []OccurrenceDetails{OccurrenceDetails{MetadataName: "yishay", Kind: "Pod"}},
					},
				},
			},
			Warning{
				Title:           "Failed with yaml validation",
				Details:         []WarningInfo{},
				InvalidYamlInfo: InvalidYamlInfo{ValidationErrors: []error{fmt.Errorf("yaml validation error")}},
			},
			Warning{
				Title:          "Failed with k8s validation",
				Details:        []WarningInfo{},
				InvalidK8sInfo: InvalidK8sInfo{ValidationErrors: []error{fmt.Errorf("K8S validation error")}, K8sVersion: "1.18.0"},
			},
		}

		printer.PrintWarnings(warnings)

		got := out.(*bytes.Buffer).Bytes()

		expected := []byte(
			`Failed with Occurrences
[V] YAML validation
[V] Kubernetes schema validation

[X] Policy check

❌  Caption  [1 occurrence]
    — metadata.name: yishay (kind: Pod)
💡  Suggestion

Failed with yaml validation
[X] YAML validation

❌  yaml validation error

[?] Kubernetes schema validation didn't run for this file
[?] Policy check didn't run for this file

Failed with k8s validation
[V] YAML validation
[X] Kubernetes schema validation

❌  K8S validation error

[?] Policy check didn't run for this file


`)
		assert.Equal(t, string(expected), string(got))
	})
}
