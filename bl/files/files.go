package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/datreeio/datree/bl/validation"
	"github.com/datreeio/datree/pkg/extractor"
)

func ToAbsolutePath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	fileInfo, _ := os.Stat(absolutePath)
	if fileInfo != nil && !fileInfo.IsDir() {
		return filepath.Abs(absolutePath)
	}

	return "", fmt.Errorf("failed parsing absolute path %s", path)
}

func ExtractFilesConfigurations(paths []string, concurrency int) (chan *extractor.FileConfigurations, chan *validation.InvalidYamlFile) {
	filesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidFilesChan := make(chan *validation.InvalidYamlFile, concurrency)

	go func() {
		defer func() {
			close(filesConfigurationsChan)
			close(invalidFilesChan)
		}()

		for _, path := range paths {

			absolutePath, err := ToAbsolutePath(path)
			if err != nil {
				invalidFilesChan <- &validation.InvalidYamlFile{Path: path, ValidationErrors: []error{err}}
				continue
			}

			content, err := extractor.ReadFileContent(absolutePath)
			if err != nil {
				invalidFilesChan <- &validation.InvalidYamlFile{Path: absolutePath, ValidationErrors: []error{err}}
				continue
			}

			configurations, err := extractor.ParseYaml(content)
			if err != nil {
				invalidFilesChan <- &validation.InvalidYamlFile{Path: absolutePath, ValidationErrors: []error{err}}
				continue
			}

			filesConfigurationsChan <- &extractor.FileConfigurations{FileName: absolutePath, Configurations: *configurations}
		}
	}()

	return filesConfigurationsChan, invalidFilesChan
}
