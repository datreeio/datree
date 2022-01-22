package extractor

import (
	"os"

    "sigs.k8s.io/yaml"
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

	jsonByte, err := yaml.YAMLToJSON([]byte(content)) 

		var doc = map[string]interface{}{}
		err = yaml.Unmarshal(jsonByte, &doc)

		if len(doc) > 0 {
			configurations = append(configurations, doc)
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
