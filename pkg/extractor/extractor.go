package extractor

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type FileReader interface {
	ReadFileContent(filepath string) (string, error)
}

type FileConfiguration struct {
	FileName       string
	Configurations []map[string]interface{}
}

type Error struct {
	Message  string
	Filename string
}

func ExtractConfiguration(path string) (*FileConfiguration, *Error) {
	content, err := readFileContent(path)
	if err != nil {
		return nil, &Error{Filename: path, Message: err.Error()}
	}

	configurations, err := yamlToK8sConfigurations(content)
	if err != nil {
		return nil, &Error{Filename: path, Message: err.Error()}
	} else {
		f := &FileConfiguration{
			Configurations: *configurations,
			FileName:       path,
		}
		return f, nil
	}
}

func yamlToK8sConfigurations(content string) (*[]map[string]interface{}, error) {
	var configurations []map[string]interface{}

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

func readFileContent(filepath string) (string, error) {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
