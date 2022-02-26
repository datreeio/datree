package files

import (
	"bytes"
	"os"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/datreeio/datree/pkg/extractor"
)

type UnknownStruct map[string]interface{}

type FilesExtractor struct{}

type FilesExtractorInterface interface {
	ExtractFilesConfigurations(paths []string, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile)
	ExtractYamlFileToUnknownStruct(path string) (UnknownStruct, error)
}

func New() *FilesExtractor {
	return &FilesExtractor{}
}

func (f *FilesExtractor) ExtractFilesConfigurations(paths []string, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile) {
	filesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidFilesChan := make(chan *extractor.InvalidFile, concurrency)
	pathsChan := make(chan string, concurrency)

	go func() {
		defer close(pathsChan)
		for _, path := range paths {
			pathsChan <- path
		}
	}()

	go func() {
		defer func() {
			close(filesConfigurationsChan)
			close(invalidFilesChan)
		}()

		var wg sync.WaitGroup
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				for path := range pathsChan {
					configurations, absolutePath, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(path)

					if invalidYamlFile != nil {
						invalidFilesChan <- invalidYamlFile
						continue
					}

					filesConfigurationsChan <- &extractor.FileConfigurations{FileName: absolutePath, Configurations: *configurations}
				}
			}()
		}
		wg.Wait()
	}()

	return filesConfigurationsChan, invalidFilesChan
}

func (f *FilesExtractor) ExtractYamlFileToUnknownStruct(path string) (UnknownStruct, error) {
	absolutePath, err := extractor.ToAbsolutePath(path)
	if err != nil {
		return nil, err
	}

	yamlContent, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(yamlContent))
	var policies = UnknownStruct{}
	err = yamlDecoder.Decode(&policies)

	if err != nil {
		return nil, err
	}

	return policies, nil
}
