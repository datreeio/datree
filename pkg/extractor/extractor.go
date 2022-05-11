package extractor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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

type Configuration map[string]interface{}

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
		var doc = map[string]interface{}{}
		err = yamlDecoder.Decode(&doc)
		if err != nil {
			break
		}

		if len(doc) > 0 {
			configurations = append(configurations, doc)
		}
	}

	if err == io.EOF {
		err = nil
	}

	return &configurations, err
}

func ReadFileContent(filepath string) (string, error) {
	dat, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
