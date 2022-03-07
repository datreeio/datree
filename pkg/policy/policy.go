package policy

import (
	"github.com/datreeio/datree/pkg/fileReader"
	"gopkg.in/yaml.v3"
)

const defaultRulesYamlPath = "./pkg/policy/defaultRules.yaml"

type DefaultRulesDefinitions struct {
	ApiVersion string                  `yaml:"apiVersion"`
	Rules      []DefaultRuleDefinition `yaml:"rules"`
}

type DefaultRuleDefinition struct {
	ID               int                    `yaml:"id"`
	Name             string                 `yaml:"name"`
	UniqueName       string                 `yaml:"uniqueName"`
	EnabledByDefault bool                   `yaml:"enabledByDefault"`
	DocumentationUrl string                 `yaml:"documentationUrl"`
	MessageOnFailure string                 `yaml:"messageOnFailure"`
	Category         string                 `yaml:"category"`
	Schema           map[string]interface{} `yaml:"schema"`
}

func GetDefaultRules() (*DefaultRulesDefinitions, error) {
	fileReader := fileReader.CreateFileReader(nil)
	defaultRulesYaml, err := fileReader.ReadFileContent(defaultRulesYamlPath)
	if err != nil {
		return nil, err
	}

	defaultRulesDefinitions, err := yamlToStruct(defaultRulesYaml)
	return defaultRulesDefinitions, err
}

func yamlToStruct(content string) (*DefaultRulesDefinitions, error) {
	var defaultRulesDefinitions DefaultRulesDefinitions
	err := yaml.Unmarshal([]byte(content), &defaultRulesDefinitions)
	if err != nil {
		return nil, err
	}
	return &defaultRulesDefinitions, err
}
