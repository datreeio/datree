package extractor

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type FileReader interface {
	ReadFileContent(filepath string) (string, error)
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
		configurations = append(configurations, doc)
	}

	if err == io.EOF {
		err = nil
	}

	return &configurations, err
}

func ReadFileContent(filepath string) (string, error) {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
