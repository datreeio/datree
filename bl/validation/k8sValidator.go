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

type K8sValidationWarningPerValidFile map[string]*ValidationWarning

func New() *K8sValidator {
	return &K8sValidator{}
}

func (val *K8sValidator) InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) {
	val.validationClient = newKubconformValidator(k8sVersion, ignoreMissingSchemas, append(getDefaultSchemaLocation(), schemaLocations...))
}

func (val *K8sValidator) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile, K8sValidationWarningPerValidFile) {
	validK8sFilesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	k8sValidationWarningPerValidFile := make(K8sValidationWarningPerValidFile)
	invalidK8sFilesChan := make(chan *extractor.InvalidFile, concurrency)

	go func() {
		defer func() {
			close(invalidK8sFilesChan)
			close(validK8sFilesConfigurationsChan)
		}()

		for fileConfigurations := range filesConfigurationsChan {

			isValid, validationErrors, err, validationWarning := val.validateResource(fileConfigurations.FileName)
			if err != nil {
				invalidK8sFilesChan <- &extractor.InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: []error{err},
				}
			}
			if isValid {
				validK8sFilesConfigurationsChan <- fileConfigurations
				k8sValidationWarningPerValidFile[fileConfigurations.FileName] = validationWarning
				_ = validationWarning
			} else {
				invalidK8sFilesChan <- &extractor.InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: validationErrors,
				}
			}
		}
	}()
	return validK8sFilesConfigurationsChan, invalidK8sFilesChan, k8sValidationWarningPerValidFile
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

func (val *K8sValidator) validateResource(filepath string) (bool, []error, error, *ValidationWarning) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, []error{}, fmt.Errorf("failed opening %s: %s", filepath, &InvalidK8sSchemaError{ErrorMessage: err.Error()}), nil
	}

	defer func() {
		f.Close()
	}()
	// defer f.Close()

	results := val.validationClient.Validate(filepath, f)

	// Return an error if no valid configurations found
	// Empty files are throwing errors in k8s
	if isEveryResultStatusEmpty(results) {
		return false, []error{&InvalidK8sSchemaError{ErrorMessage: "empty file"}}, nil, nil
	}

	isValid := true
	var validationErrors []error
	for _, res := range results {
		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid || res.Status == kubeconformValidator.Error {
			if val.isNetworkError(res.Err.Error()) {
				noConnectionWarning := &ValidationWarning{Message: "k8s schema validation skipped: no internet connection"}
				return isValid, []error{}, nil, noConnectionWarning
			}
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

	return isValid, validationErrors, nil, nil
}

func (val *K8sValidator) isNetworkError(errorString string) bool {
	return strings.Contains(errorString, "no such host")
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

func getDefaultSchemaLocation() []string {
	defaultSchemaLocations := [...]string{
		"default",
		getDatreeCRDSchemaByName("argo"),
	}
	return (defaultSchemaLocations[:])
}

func getDatreeCRDSchemaByName(crdCatalogName string) string {
	crdCatalog := "https://raw.githubusercontent.com/datreeio/CRDs-catalog/master/" + crdCatalogName + "/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json"
	return crdCatalog
}
