package files

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/datreeio/datree/pkg/extractor"
)

type UnknownStruct map[string]interface{}



func ExtractFilesConfigurations(paths []string, concurrency int) (chan *extractor.FileConfigurations, chan *extractor.InvalidFile) {
	filesConfigurationsChan := make(chan *extractor.FileConfigurations, concurrency)
	invalidFilesChan := make(chan *extractor.InvalidFile, concurrency)

	go func() {
		defer func() {
			close(filesConfigurationsChan)
			close(invalidFilesChan)
		}()

		for _, path := range paths {

			configurations, absolutePath, invalidYamlFile := extractor.ExtractConfigurationsFromYamlFile(path)

			if invalidYamlFile != nil {
				invalidFilesChan <- invalidYamlFile
				continue
			}

			filesConfigurationsChan <- &extractor.FileConfigurations{FileName: absolutePath, Configurations: *configurations}
		}
	}()

	return filesConfigurationsChan, invalidFilesChan
}

func ExtractYamlFileToUnknownStruct(path string) (UnknownStruct, error) {
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
