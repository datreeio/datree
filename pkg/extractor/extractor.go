package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	yamlConvertor "sigs.k8s.io/yaml"
)

type InvalidFile struct {
	Path             string  `yaml:"path" json:"path" xml:"path"`
	ValidationErrors []error `yaml:"errors" json:"errors" xml:"errors"`
}

type InvalidYamlError struct {
	ErrorMessage string
}

func (e *InvalidYamlError) Error() string {
	return fmt.Sprintf("yaml validation error: %s\n", e.ErrorMessage)
}

type FileReader interface {
	ReadFileContent(filepath string) (string, error)
}

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

func ExtractConfigurationsFromYamlFile(path string) (*[]Configuration, string, *InvalidFile) {
	absolutePath, err := ToAbsolutePath(path)
	if err != nil {
		return nil, "", &InvalidFile{Path: path, ValidationErrors: []error{&InvalidYamlError{ErrorMessage: err.Error()}}}
	}

	content, err := ReadFileContent(absolutePath)
	if err != nil {
		return nil, "", &InvalidFile{Path: absolutePath, ValidationErrors: []error{&InvalidYamlError{ErrorMessage: err.Error()}}}
	}

	configurations, err := ParseYaml(content)
	if err != nil {
		return nil, "", &InvalidFile{Path: absolutePath, ValidationErrors: []error{&InvalidYamlError{ErrorMessage: err.Error()}}}
	}

	return configurations, absolutePath, nil
}

type Configuration struct {
	MetadataName string
	Kind         string
	ApiVersion   string
	Annotations  map[string]interface{}
	Payload      []byte
}

type FileConfigurations struct {
	FileName       string          `json:"fileName"`
	Configurations []Configuration `json:"configurations"`
}

func ParseYaml(content string) (*[]Configuration, error) {
	configurations, err := extractYamlConfigurations(content)
	if err != nil {
		return nil, err
	} else {
		return configurations, nil
	}
}

func extractYamlConfigurations(content string) (*[]Configuration, error) {
	var configurations []Configuration

	yamlDecoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))

	var err error
	for {
		var anyTypeStore interface{}
		err = yamlDecoder.Decode(&anyTypeStore)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		yamlByteArray, err := yaml.Marshal(anyTypeStore)
		if err != nil {
			return nil, err
		}

		jsonByte, err := yamlConvertor.YAMLToJSON(yamlByteArray)
		if err != nil {
			return nil, err
		}

		configurations = append(configurations, extractConfigurationK8sData(jsonByte))
	}

	return &configurations, nil
}

func extractConfigurationK8sData(content []byte) Configuration {
	var configuration Configuration
	var jsonObject map[string]interface{}
	configuration.Payload = content
	err := json.Unmarshal(content, &jsonObject)

	if err != nil {
		return configuration
	}

	if metadata := jsonObject["metadata"]; metadata != nil {
		if metadataName := metadata.(map[string]interface{})["name"]; metadataName != nil {
			configuration.MetadataName = metadataName.(string)
			if annotations := metadata.(map[string]interface{})["annotations"]; annotations != nil {
				configuration.Annotations = annotations.(map[string]interface{})
			}
		}
	}

	if apiVersion := jsonObject["apiVersion"]; apiVersion != nil {
		configuration.ApiVersion = apiVersion.(string)
	}

	if kind := jsonObject["kind"]; kind != nil {
		configuration.Kind = kind.(string)
	}

	return configuration
}

func ReadFileContent(filepath string) (string, error) {
	dat, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
