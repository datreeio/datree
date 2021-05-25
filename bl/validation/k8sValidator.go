package validation

import (
	"fmt"
	"io"
	"os"

	"github.com/datreeio/datree/bl/files"
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

func (val *K8sValidator) InitClient(k8sVersion string) {
	val.validationClient = newKubconformValidator(k8sVersion)
}

func (val *K8sValidator) ValidateResources(paths []string) (chan string, chan *InvalidFile, chan error) {
	pathsChan := files.ToAbsolutePaths(paths)

	errorChan := make(chan error)
	validFilesPathChan := make(chan string)
	invalidFilesPathsChan := make(chan *InvalidFile)

	go func() {
		for path := range pathsChan {
			isValid, validationErrors, err := val.validateResource(path)
			if isValid {
				validFilesPathChan <- path
			} else {
				invalidFilesPathsChan <- &InvalidFile{
					Path:             path,
					ValidationErrors: validationErrors,
				}
			}
			if err != nil {
				errorChan <- err
			}
		}

		close(invalidFilesPathsChan)
		close(validFilesPathChan)
		close(errorChan)
	}()

	return validFilesPathChan, invalidFilesPathsChan, errorChan
}

func (val *K8sValidator) validateResource(filepath string) (bool, []error, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, []error{}, fmt.Errorf("failed opening %s: %s", filepath, err)
	}

	results := val.validationClient.Validate(filepath, f)
	isValid := true
	var validationErrors []error
	for _, res := range results {

		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid || res.Status == kubeconformValidator.Error {
			isValid = false
			validationErrors = append(validationErrors, res.Err)
		}
	}

	return isValid, validationErrors, nil
}

func newKubconformValidator(k8sVersion string) ValidationClient {
	v, _ := kubeconformValidator.New(nil, kubeconformValidator.Opts{Strict: true, KubernetesVersion: k8sVersion})
	return v
}
