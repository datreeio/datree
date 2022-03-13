package policy

import (
	"embed"

	"gopkg.in/yaml.v3"
)

//go:embed defaultRules.yaml
var f embed.FS

type DefaultRulesDefinitions struct {
	ApiVersion string                   `yaml:"apiVersion"`
	Rules      []*DefaultRuleDefinition `yaml:"rules"`
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
	defaultRulesYaml, err := f.ReadFile("defaultRules.yaml")

	if err != nil {
		return nil, err
	}

	defaultRulesDefinitions, err := yamlToStruct(string(defaultRulesYaml))
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
