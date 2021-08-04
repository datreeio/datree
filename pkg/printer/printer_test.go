package printer

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// add to test and run to get expected string
// got := out.(*bytes.Buffer).String()
// fmt.Sprintf(got)

func TestWarnings(t *testing.T) {
	t.Run("Test Warning", func(t *testing.T) {

		out = new(bytes.Buffer)

		printer := CreateNewPrinter()

		warnings := []Warning{
			Warning{
				Title:   "Failed with Occurrences",
				Details: []WarningInfo{WarningInfo{Caption: "Caption", Occurrences: 5, Suggestion: "Suggestion"}},
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

		expected := []byte("Failed with Occurrences[V] YAML validation\n[V] Kubernetes schema validation\n[X] Policy check\nFailed with yaml validation[X] YAML validation\n[?] Kubernetes schema validation didn’t run for this file\n[?] Policy check didn’t run for this file\nFailed with k8s validation[V] YAML validation\n[X] Kubernetes schema validation\n[?] Policy check didn’t run for this file\n\n\n\n❌  Caption  [5 occurrences]\n💡  Suggestion\n\n\n\n❌  yaml validation error\n\n\n\n\n❌  K8S validation error\n\n\n")

		fmt.Print(got)
		//res := bytes.Compare(got, expected)
		assert.Equal(t, expected, got)
	})
}
