package printer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintWarnings(t *testing.T) {
	printer := CreateNewPrinter()

	warnings := []Warning{{
		Title: "~/.datree/k8-demo.yaml",
		FailedRules: []FailedRule{
			{
				Name:               "Caption",
				Occurrences:        1,
				Suggestion:         "Suggestion",
				OccurrencesDetails: []OccurrenceDetails{{MetadataName: "yishay", Kind: "Pod"}},
			},
		},
	},
		{
			Title:           "/datree/datree/internal/fixtures/kube/yaml-validation-error.yaml\n",
			FailedRules:     []FailedRule{},
			InvalidYamlInfo: InvalidYamlInfo{ValidationErrors: []error{fmt.Errorf("yaml validation error")}},
		},
		{
			Title:          "/datree/datree/internal/fixtures/kube/k8s-validation-error.yaml\n",
			FailedRules:    []FailedRule{},
			InvalidK8sInfo: InvalidK8sInfo{ValidationErrors: []error{fmt.Errorf("K8S validation error")}, K8sVersion: "1.18.0"},
		},
		{
			Title:          "/datree/datree/internal/fixtures/kube/Chart.yaml\n",
			FailedRules:    []FailedRule{},
			InvalidK8sInfo: InvalidK8sInfo{ValidationErrors: []error{fmt.Errorf("K8S validation error")}, K8sVersion: "1.18.0"},
			ExtraMessages: []ExtraMessage{{Text: "Are you trying to test a raw helm file? To run Datree with Helm - check out the helm plugin README:\nhttps://github.com/datreeio/helm-datree",
				Color: "cyan"}},
		},
	}

	t.Run("Test PrintWarnings", func(t *testing.T) {

		out = new(bytes.Buffer)

		printer.PrintWarnings(warnings)

		got := out.(*bytes.Buffer).Bytes()

		expected := []byte(
			`>>  File: ~/.datree/k8-demo.yaml

[V] YAML validation
[V] Kubernetes schema validation

[X] Policy check

❌  Caption  [1 occurrence]
    — metadata.name: yishay (kind: Pod)
💡  Suggestion

>>  File: /datree/datree/internal/fixtures/kube/yaml-validation-error.yaml


[X] YAML validation

❌  yaml validation error

[?] Kubernetes schema validation didn't run for this file
[?] Policy check didn't run for this file

>>  File: /datree/datree/internal/fixtures/kube/k8s-validation-error.yaml


[V] YAML validation
[X] Kubernetes schema validation

❌  K8S validation error

[?] Policy check didn't run for this file

>>  File: /datree/datree/internal/fixtures/kube/Chart.yaml


[V] YAML validation
[X] Kubernetes schema validation

❌  K8S validation error
Are you trying to test a raw helm file? To run Datree with Helm - check out the helm plugin README:
https://github.com/datreeio/helm-datree
[?] Policy check didn't run for this file


`)
		assert.Equal(t, string(expected), string(got))
	})

	t.Run("Test PrintWarnings simple output", func(t *testing.T) {

		out = new(bytes.Buffer)

		printer.SetTheme(CreateSimpleTheme())

		printer.PrintWarnings(warnings)

		got := out.(*bytes.Buffer).Bytes()

		expected := []byte(
			`>>  File: ~/.datree/k8-demo.yaml

[V] YAML validation
[V] Kubernetes schema validation

[X] Policy check

[X]  Caption  [1 occurrence]
    — metadata.name: yishay (kind: Pod)
[*]  Suggestion

>>  File: /datree/datree/internal/fixtures/kube/yaml-validation-error.yaml


[X] YAML validation

[X]  yaml validation error

[?] Kubernetes schema validation didn't run for this file
[?] Policy check didn't run for this file

>>  File: /datree/datree/internal/fixtures/kube/k8s-validation-error.yaml


[V] YAML validation
[X] Kubernetes schema validation

[X]  K8S validation error

[?] Policy check didn't run for this file

>>  File: /datree/datree/internal/fixtures/kube/Chart.yaml


[V] YAML validation
[X] Kubernetes schema validation

[X]  K8S validation error
Are you trying to test a raw helm file? To run Datree with Helm - check out the helm plugin README:
https://github.com/datreeio/helm-datree
[?] Policy check didn't run for this file


`)
		assert.Equal(t, string(expected), string(got))
	})
}

func TestPrintEvaluationSummary(t *testing.T) {
	t.Run("Test PrintEvaluationSummary", func(t *testing.T) {
		out = new(bytes.Buffer)
		printer := CreateNewPrinter()
		summary := EvaluationSummary{
			ConfigsCount:              6,
			RulesCount:                21,
			FilesCount:                5,
			PassedYamlValidationCount: 4,
			K8sValidation:             "3/5",
			PassedPolicyCheckCount:    2,
		}
		k8sVersion := "1.2.3"

		printer.PrintEvaluationSummary(summary, k8sVersion)
		expected := []byte(`(Summary)

- Passing YAML validation: 4/5

- Passing Kubernetes (1.2.3) schema validation: 3/5

- Passing policy check: 2/5

`)

		got := out.(*bytes.Buffer).Bytes()

		assert.Equal(t, string(expected), string(got))

	})

	t.Run("Test PrintEvaluationSummary with no connection warning", func(t *testing.T) {
		out = new(bytes.Buffer)
		printer := CreateNewPrinter()
		summary := EvaluationSummary{
			ConfigsCount:              6,
			RulesCount:                21,
			FilesCount:                5,
			PassedYamlValidationCount: 4,
			K8sValidation:             "no internet connection",
			PassedPolicyCheckCount:    2,
		}
		k8sVersion := "1.2.3"

		printer.PrintEvaluationSummary(summary, k8sVersion)
		expected := []byte(`(Summary)

- Passing YAML validation: 4/5

- Passing Kubernetes (1.2.3) schema validation: no internet connection

- Passing policy check: 2/5

`)

		got := out.(*bytes.Buffer).Bytes()

		assert.Equal(t, string(expected), string(got))

	})
}
