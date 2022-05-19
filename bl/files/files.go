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

	pathsChan := parsePathsArrayToChan(paths, concurrency)

	go func() {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				for {
					path, ok := <-pathsChan
					if !ok {
						break
					}
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
		close(filesConfigurationsChan)
		close(invalidFilesChan)
	}()

	return filesConfigurationsChan, invalidFilesChan
}

func parsePathsArrayToChan(paths []string, concurrency int) chan string {
	pathsChan := make(chan string, concurrency)

	go func() {
		for _, path := range paths {
			pathsChan <- path
		}
		close(pathsChan)
	}()

	return pathsChan
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

func (f *FilesExtractor) ExtractConfigurationsFromYamlFile(path string) (*[]extractor.Configuration, string, *extractor.InvalidFile) {
	return extractor.ExtractConfigurationsFromYamlFile(path)
}
