package validator

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/datreeio/datree/bl/files"
	kubeconformValidator "github.com/yannh/kubeconform/pkg/validator"
)

type ValidationClient interface {
	Validate(filename string, r io.ReadCloser) []kubeconformValidator.Result
}

type Validator struct {
	validationClient ValidationClient
}

func NewKubconformValidator(k8sVersion string) ValidationClient {
	v, _ := kubeconformValidator.New(nil, kubeconformValidator.Opts{Strict: true, KubernetesVersion: k8sVersion})
	return v
}

func New(k8sVersion string) *Validator {
	kubconformClient := NewKubconformValidator(k8sVersion)
	return &Validator{
		validationClient: kubconformClient,
	}
}

// example filepath: "../fixtures/valid.yaml"
func (val *Validator) validateFile(filepath string) (bool, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("failed opening %s: %s", filepath, err)
	}

	for i, res := range val.validationClient.Validate(filepath, f) {
		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == kubeconformValidator.Invalid {
			return false, fmt.Errorf("resource %d in file %s is not valid: %s", i, filepath, res.Err)
		}
		if res.Status == kubeconformValidator.Error {
			return false, fmt.Errorf("error while processing resource %d in file %s: %s", i, filepath, res.Err)
		}
	}

	return true, nil
}

func (val *Validator) Validate(paths []string) (<-chan string, <-chan string, <-chan error) {
	pathsChan, _ := files.ToAbsolutePaths(paths)

	errorChan := make(chan error, 100)
	invalidFilesPathsChan := make(chan string, 100)
	validFilesPathChan := make(chan string, 100)

	conc := 10
	wg := sync.WaitGroup{}
	wg.Add(conc)

	go func() {
		for i := 0; i < conc; i++ {
			go func() {
				for {
					path, ok := <-pathsChan
					if !ok {
						break
					}

					isValid, err := val.validateFile(path)
					if err != nil {
						errorChan <- err
					}
					if isValid {
						validFilesPathChan <- path
					} else {
						invalidFilesPathsChan <- path
					}

				}
			}()
			wg.Done()
		}

		wg.Wait()
		close(invalidFilesPathsChan)
		close(validFilesPathChan)
		close(errorChan)
	}()

	return validFilesPathChan, invalidFilesPathsChan, errorChan
}
