package validation

import (
	"fmt"
	"io"
	"os"

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

type InvalidFile struct {
	Path             string
	ValidationErrors []error
}
type InvalidYamlFile InvalidFile

type InvalidK8sFile InvalidFile

func (val *K8sValidator) InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) {
	val.validationClient = newKubconformValidator(k8sVersion, ignoreMissingSchemas, schemaLocations)
}

func (val *K8sValidator) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int, onlyK8sFiles bool) (chan *extractor.FileConfigurations, chan *InvalidK8sFile, chan *string) {
	validK8sFilesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidK8sFilesChan := make(chan *InvalidK8sFile, concurrency)
	ignoredFilesChan := make(chan *string, concurrency)

	go func() {
		defer func() {
			close(invalidK8sFilesChan)
			close(validK8sFilesConfigurationsChan)
			close(ignoredFilesChan)
		}()

		for fileConfigurations := range filesConfigurationsChan {
			if onlyK8sFiles {
				if ok := val.validatePotentialK8sFile(fileConfigurations.Configurations); !ok {
					ignoredFilesChan <- &fileConfigurations.FileName
					continue
				}
			}

			isValid, validationErrors, err := val.validateResource(fileConfigurations.FileName)
			if err != nil {
				invalidK8sFilesChan <- &InvalidK8sFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: []error{err},
				}
				continue
			}
			if isValid {
				validK8sFilesConfigurationsChan <- fileConfigurations
			} else {
				invalidK8sFilesChan <- &InvalidK8sFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: validationErrors,
				}
			}
		}
	}()
	return validK8sFilesConfigurationsChan, invalidK8sFilesChan, ignoredFilesChan
}

func (val *K8sValidator) validatePotentialK8sFile(fileConfigurations []extractor.Configuration) bool {
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

	results := val.validationClient.Validate(filepath, f)
	isValid := true
	var validationErrors []error
	for _, res := range results {

		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid || res.Status == kubeconformValidator.Error {
			isValid = false
			validationErrors = append(validationErrors, &InvalidK8sSchemaError{ErrorMessage: res.Err.Error()})
		}
	}

	return isValid, validationErrors, nil
}

func newKubconformValidator(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) ValidationClient {
	v, _ := kubeconformValidator.New(schemaLocations, kubeconformValidator.Opts{Strict: true, KubernetesVersion: k8sVersion, IgnoreMissingSchemas: ignoreMissingSchemas})
	return v
}
