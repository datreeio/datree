package validation

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/datreeio/datree/pkg/extractor"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)

type ValidationClient interface {
	Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result
}

type K8sValidator struct {
	validationClient ValidationClient
}

func New() *K8sValidator {
	return &K8sValidator{}
}

func (val *K8sValidator) InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) {
	val.validationClient = newKubconformValidator(k8sVersion, ignoreMissingSchemas, schemaLocations)
}

func (val *K8sValidator) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile) {
	validK8sFilesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidK8sFilesChan := make(chan *extractor.InvalidFile, concurrency)

	go func() {
		defer func() {
			close(invalidK8sFilesChan)
			close(validK8sFilesConfigurationsChan)
		}()

		for fileConfigurations := range filesConfigurationsChan {

			isValid, validationErrors, err := val.validateResource(fileConfigurations.FileName)
			if err != nil {
				invalidK8sFilesChan <- &extractor.InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: []error{err},
				}
				continue
			}
			if isValid {
				validK8sFilesConfigurationsChan <- fileConfigurations
			} else {
				invalidK8sFilesChan <- &extractor.InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: validationErrors,
				}
			}
		}
	}()
	return validK8sFilesConfigurationsChan, invalidK8sFilesChan
}

func (val *K8sValidator) GetK8sFiles(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.FileConfigurations) {
	k8sFilesChan := make(chan *extractor.FileConfigurations, concurrency)
	ignoredYamlFilesChan := make(chan *extractor.FileConfigurations, concurrency)

	go func() {
		defer func() {
			close(k8sFilesChan)
			close(ignoredYamlFilesChan)
		}()

		for fileConfigurations := range filesConfigurationsChan {
			if ok := val.isK8sFile(fileConfigurations.Configurations); ok {
				k8sFilesChan <- fileConfigurations
			} else {
				ignoredYamlFilesChan <- fileConfigurations
			}
		}
	}()

	return k8sFilesChan, ignoredYamlFilesChan
}

func (val *K8sValidator) isK8sFile(fileConfigurations []extractor.Configuration) bool {
	for _, configuration := range fileConfigurations {
		_, has_apiVersion := configuration["apiVersion"]
		_, has_kind := configuration["kind"]

		if !has_apiVersion || !has_kind {
			return false
		}
	}

	return true
}

func (val *K8sValidator) validateResource(filepath string) (bool, []error, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, []error{}, fmt.Errorf("failed opening %s: %s", filepath, &InvalidK8sSchemaError{ErrorMessage: err.Error()})
	}
	defer f.Close()

	results := val.validationClient.Validate(filepath, f)

	// Return an error if no valid configurations found
	// Empty files are throwing errors in k8s
	if isEveryResultStatusEmpty(results) {
		return false, []error{&InvalidK8sSchemaError{ErrorMessage: "empty file"}}, nil
	}

	isValid := true
	var validationErrors []error
	for _, res := range results {

		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid || res.Status == kubeconformValidator.Error {
			isValid = false
			errorMessages := strings.Split(res.Err.Error(), "-")

			// errorMessages slice is not empty
			if len(errorMessages) > 0 {
				for _, errorMessage := range errorMessages {
					msg := strings.Trim(errorMessage, " ")
					validationErrors = append(validationErrors, &InvalidK8sSchemaError{ErrorMessage: msg})
				}
			}
		}
	}

	return isValid, validationErrors, nil
}

func newKubconformValidator(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) ValidationClient {
	v, _ := kubeconformValidator.New(schemaLocations, kubeconformValidator.Opts{Strict: true, KubernetesVersion: k8sVersion, IgnoreMissingSchemas: ignoreMissingSchemas})
	return v
}

func isEveryResultStatusEmpty(results []kubeconformValidator.Result) bool {
	isEveryResultStatusEmpty := true
	for _, result := range results {
		if result.Status != kubeconformValidator.Empty {
			isEveryResultStatusEmpty = false
		}
	}
	return isEveryResultStatusEmpty
}
