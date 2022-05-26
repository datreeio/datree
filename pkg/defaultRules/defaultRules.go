package defaultRules

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

//go:embed defaultRules.yaml
var embeddedDefaultRulesYamlContent string

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
	configDefaultRulesYamlContent, err := getDefaultRulesFromFile()
	if err == nil {
		return yamlToStruct(configDefaultRulesYamlContent)
	}

	return yamlToStruct(embeddedDefaultRulesYamlContent)
}

func yamlToStruct(content string) (*DefaultRulesDefinitions, error) {
	var defaultRulesDefinitions DefaultRulesDefinitions
	err := yaml.Unmarshal([]byte(content), &defaultRulesDefinitions)
	if err != nil {
		return nil, err
	}
	return &defaultRulesDefinitions, err
}

func getDefaultRulesFromFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	defaultRulesFileContent, err := ioutil.ReadFile(filepath.Join(homeDir, ".datree", "defaultRules.yaml"))
	return string(defaultRulesFileContent), err
}
