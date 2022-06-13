package validation

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/datreeio/datree/pkg/extractor"
	"github.com/datreeio/datree/pkg/utils"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)

type ValidationClient interface {
	Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result
}

type K8sValidator struct {
	validationClient              ValidationClient
	isOffline                     bool
	areThereCustomSchemaLocations bool
}

type K8sValidationWarningPerValidFile map[string]FileWithWarning

func New() *K8sValidator {
	return &K8sValidator{}
}

func (val *K8sValidator) InitClient(k8sVersion string, ignoreMissingSchemas bool, userProvidedSchemaLocations []string) {
	val.isOffline = checkIsOffline()
	val.areThereCustomSchemaLocations = len(userProvidedSchemaLocations) > 0
	val.validationClient = newKubeconformValidator(k8sVersion, ignoreMissingSchemas, getAllSchemaLocations(userProvidedSchemaLocations, val.isOffline))
}

func checkIsOffline() bool {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get("https://www.githubstatus.com/api/v2/status.json")
	if err == nil && resp != nil && resp.StatusCode == 200 {
		return false
	} else {
		return true
	}
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
		has_apiVersion := configuration.ApiVersion != ""
		has_kind := configuration.Kind != ""
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

	if val.isOffline && !val.areThereCustomSchemaLocations {
		var noConnectionWarning = &validationWarning{
			WarningKind:    NetworkError,
			WarningMessage: "k8s schema validation skipped: no internet connection",
		}
		return true, []error{}, noConnectionWarning, nil
	}

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
			isValid = false
			errString := res.Err.Error()

			if utils.IsNetworkError(errString) {
				validationErrors = append(validationErrors, &InvalidK8sSchemaError{errString})
			} else {
				errorMessages := strings.Split(errString, "-")
				for _, errorMessage := range errorMessages {
					validationErrors = append(validationErrors, &InvalidK8sSchemaError{ErrorMessage: strings.Trim(errorMessage, " ")})
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

func getAllSchemaLocations(userProvidedSchemaLocations []string, isOffline bool) []string {
	if isOffline {
		return userProvidedSchemaLocations
	} else {
		// order matters! userProvidedSchemaLocations get priority over defaultSchemaLocations
		return append(userProvidedSchemaLocations, getDefaultSchemaLocations()...)
	}
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
	crdCatalog := "https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/" + crdCatalogName + "/{{ .ResourceKind }}_{{ .ResourceAPIVersion }}.json"
	return crdCatalog
}
