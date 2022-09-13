package policy

import (
	"github.com/datreeio/datree/pkg/defaultPolicies"
	"github.com/datreeio/datree/pkg/validatePoliciesYaml"

	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/ghodss/yaml"
)

func GetPoliciesFileFromPath(path string) (*defaultPolicies.EvaluationPrerunPolicies, error) {
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

	var policies *defaultPolicies.EvaluationPrerunPolicies
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
