package policy

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
)

//go:embed defaultRules.yaml
var embeddedDefaultRulesYamlContent string

//go:embed policiesSchema.json
var policiesSchemaContent string

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

func GetPoliciesFileFromPath(path string) (*cliClient.EvaluationPrerunPolicies, error) {
	fileReader := fileReader.CreateFileReader(nil)
	policiesStr, err := fileReader.ReadFileContent(path)
	if err != nil {
		return nil, err
	}

	err = validatePoliciesYaml(policiesStr, path)
	if err != nil {
		return nil, err
	}

	var policies *cliClient.EvaluationPrerunPolicies
	policiesBytes, err := yaml.YAMLToJSON([]byte(policiesStr))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(policiesBytes, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func validatePoliciesYaml(content string, policyYamlPath string) error {
	jsonSchemaValidator := jsonSchemaValidator.New()
	errorsResult, err := jsonSchemaValidator.Validate(policiesSchemaContent, content)

	if err != nil {
		return err
	}

	if errorsResult != nil {
		validationErrors := fmt.Errorf("Found errors in policies file %s:\n", policyYamlPath)

		for _, validationError := range errorsResult {
			validationErrors = fmt.Errorf("%s\n%s", validationErrors, validationError.Error)
		}

		return validationErrors
	}

	return nil
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
