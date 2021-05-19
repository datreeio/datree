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

func New(k8sVersion string) *K8sValidator {
	kubconformClient := newKubconformValidator(k8sVersion)
	return &K8sValidator{
		validationClient: kubconformClient,
	}
}

func (val *K8sValidator) ValidateResources(paths []string) (chan string, []*string, chan error) {
	pathsChan := files.ToAbsolutePaths(paths)

	var invalidFilesPaths = []*string{}
	errorChan := make(chan error)
	validFilesPathChan := make(chan string)

	go func() {
		for path := range pathsChan {
			isValid, err := val.validateResource(path)
			if isValid {
				validFilesPathChan <- path
			} else {
				invalidFilesPaths = append(invalidFilesPaths, &path)
			}
			if err != nil {
				errorChan <- err
			}
		}
		close(validFilesPathChan)
		close(errorChan)
	}()

	return validFilesPathChan, invalidFilesPaths, errorChan
}

func (val *K8sValidator) validateResource(filepath string) (bool, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("failed opening %s: %s", filepath, err)
	}

	results := val.validationClient.Validate(filepath, f)
	isValid := false
	for i, res := range results {
		if res.Status == kubeconformValidator.Valid {
			isValid = true
		}
		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid {
			fmt.Errorf("resource %d in file %s is not valid: %s", i, filepath, res.Err)
		}
		if res.Status == kubeconformValidator.Error {
			fmt.Errorf("error while processing resource %d in file %s: %s", i, filepath, res.Err)
		}
	}

	return isValid, nil
}

func newKubconformValidator(k8sVersion string) ValidationClient {
	v, _ := kubeconformValidator.New(nil, kubeconformValidator.Opts{Strict: true, KubernetesVersion: k8sVersion})
	return v
}
