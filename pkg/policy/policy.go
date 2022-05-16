package policy

import (
	"github.com/datreeio/datree/pkg/validatePoliciesYaml"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/ghodss/yaml"
)

func GetPoliciesFileFromPath(path string) (*cliClient.EvaluationPrerunPolicies, error) {
	fileReader := fileReader.CreateFileReader(nil)
	policiesStr, err := fileReader.ReadFileContent(path)
	if err != nil {
		return nil, err
	}

	policiesStrBytes := []byte(policiesStr)

	err = validatePoliciesYaml.ValidatePoliciesYaml(policiesStrBytes, path)
	if err != nil {
		return nil, err
	}

	var policies *cliClient.EvaluationPrerunPolicies
	policiesBytes, err := yaml.YAMLToJSON(policiesStrBytes)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(policiesBytes, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}
