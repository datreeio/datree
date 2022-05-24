package validation

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/utils"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)

type ValidationClient interface {
	Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result
}

type K8sValidator struct {
	validationClient ValidationClient
}

type K8sValidationWarningPerValidFile map[string]FileWithWarning

func New() *K8sValidator {
	return &K8sValidator{}
}

func (val *K8sValidator) InitClient(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) {
	val.validationClient = newKubeconformValidator(k8sVersion, ignoreMissingSchemas, getAllSchemaLocations(schemaLocations))
}

type WarningKind int

const (
	_            WarningKind = iota
	NetworkError             // a network error while validating the resource
	Skipped                  // resource has been skipped, for example if its kind was not found and the user added the --ignore-missing-schemas flag
)

type FileWithWarning struct {
	Filename    string
	Warning     string
	WarningKind WarningKind
}

func (val *K8sValidator) ValidateResources(filesConfigurationsChan chan *extractor.FileConfigurations, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile, chan *FileWithWarning) {
	validK8sFilesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidK8sFilesChan := make(chan *extractor.InvalidFile, concurrency)
	k8sValidationWarningPerValidFileChan := make(chan *FileWithWarning, concurrency)

	go func() {

		defer func() {
			close(invalidK8sFilesChan)
			close(validK8sFilesConfigurationsChan)
			close(k8sValidationWarningPerValidFileChan)
		}()

		for fileConfigurations := range filesConfigurationsChan {

			isValid, validationErrors, validationWarning, err := val.validateResource(fileConfigurations.FileName)
			if err != nil {
				invalidK8sFilesChan <- &extractor.InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: []error{err},
				}
			}
			if isValid {
				validK8sFilesConfigurationsChan <- fileConfigurations
				if validationWarning != nil {
					k8sValidationWarningPerValidFileChan <- &FileWithWarning{
						Filename:    fileConfigurations.FileName,
						Warning:     validationWarning.WarningMessage,
						WarningKind: validationWarning.WarningKind,
					}
				}
			} else {
				invalidK8sFilesChan <- &extractor.InvalidFile{
					Path:             fileConfigurations.FileName,
					ValidationErrors: validationErrors,
				}
			}
		}
	}()

	return validK8sFilesConfigurationsChan, invalidK8sFilesChan, k8sValidationWarningPerValidFileChan
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

type validationWarning struct {
	WarningKind    WarningKind
	WarningMessage string
}

func (val *K8sValidator) validateResource(filepath string) (bool, []error, *validationWarning, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, []error{}, nil, fmt.Errorf("failed opening %s: %s", filepath, &InvalidK8sSchemaError{ErrorMessage: err.Error()})
	}

	defer f.Close()

	results := val.validationClient.Validate(filepath, f)

	// Return an error if no valid configurations found
	// Empty files are throwing errors in k8s
	if isEveryResultStatusEmpty(results) {
		return false, []error{&InvalidK8sSchemaError{ErrorMessage: "empty file"}}, nil, nil
	}

	isValid := true
	isAtLeastOneConfigSkipped := false
	var validationErrors []error
	for _, res := range results {
		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Skipped {
			isAtLeastOneConfigSkipped = true
		}
		if res.Status == kubeconformValidator.Invalid || res.Status == kubeconformValidator.Error {
			if utils.IsNetworkError(res.Err.Error()) {
				noConnectionWarning := &validationWarning{
					WarningKind:    NetworkError,
					WarningMessage: "k8s schema validation skipped: no internet connection",
				}
				return true, []error{}, noConnectionWarning, nil
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
	var warning *validationWarning = nil
	if isAtLeastOneConfigSkipped && isValid {
		warning = &validationWarning{
			WarningKind:    Skipped,
			WarningMessage: "k8s schema validation skipped: --ignore-missing-schemas flag was used",
		}
	}
	return isValid, validationErrors, warning, nil
}

func newKubeconformValidator(k8sVersion string, ignoreMissingSchemas bool, schemaLocations []string) ValidationClient {
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

func getAllSchemaLocations(userProvidedSchemaLocations []string) []string {
	// order matters!
	// it's important that provided schema locations (from --schema-locations flag) are *before* the default schema locations
	// this will give them priority and allow using a local schema in offline mode
	return append(userProvidedSchemaLocations, getDefaultSchemaLocations()...)
}

func getDefaultSchemaLocations() []string {
	return []string{
		"default",
		// this is a workaround for https://github.com/yannh/kubeconform/issues/100
		// notice: order here is important because this fallback doesn't have strict mode enabled (in contrast to "default")
		"https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
		getDatreeCRDSchemaByName("argo"),
	}
}

func getDatreeCRDSchemaByName(crdCatalogName string) string {
	crdCatalog := "https://raw.githubusercontent.com/datreeio/CRDs-catalog/master/" + crdCatalogName + "/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json"
	return crdCatalog
}
