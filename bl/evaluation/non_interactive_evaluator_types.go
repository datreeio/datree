package evaluation

import (
	"github.com/datreeio/datree/bl/validation"
)

type FormattedOutput struct {
	PolicyValidationResults []*FormattedEvaluationResults   `yaml:"policyValidationResults" json:"policyValidationResults"`
	PolicySummary           PolicySummary                   `yaml:"policySummary" json:"policySummary"`
	EvaluationSummary       NonInteractiveEvaluationSummary `yaml:"evaluationSummary" json:"evaluationSummary"`
	YamlValidationResults   []*validation.InvalidYamlFile   `yaml:"yamlValidationResults" json:"yamlValidationResults"`
	K8sValidationResults    []*validation.InvalidK8sFile    `yaml:"k8sValidationResults" json:"k8sValidationResults"`
}

type NonInteractiveEvaluationResults struct {
	FormattedEvaluationResults []*FormattedEvaluationResults
	PolicySummary              PolicySummary
}

type FormattedEvaluationResults struct {
	FileName    string        `yaml:"fileName" json:"fileName"`
	RuleResults []*RuleResult `yaml:"ruleResults" json:"ruleResults"`
}

type RuleResult struct {
	Identifier         string              `json:"identifier"`
	Name               string              `json:"name"`
	MessageOnFailure   string              `yaml:"messageOnFailure" json:"messageOnFailure"`
	OccurrencesDetails []OccurrenceDetails `yaml:"occurrencesDetails" json:"occurrencesDetails"`
}

type NonInteractiveEvaluationSummary struct {
	ConfigsCount                int `yaml:"configsCount" json:"configsCount"`
	FilesCount                  int `yaml:"filesCount" json:"filesCount"`
	PassedYamlValidationCount   int `yaml:"passedYamlValidationCount" json:"passedYamlValidationCount"`
	PassedK8sValidationCount    int `yaml:"passedK8sValidationCount" json:"passedK8sValidationCount"`
	PassedPolicyValidationCount int `yaml:"passedPolicyValidationCount" json:"passedPolicyValidationCount"`
}

type PolicySummary struct {
	PolicyName         string `yaml:"policyName" json:"policyName"`
	TotalRulesInPolicy int    `yaml:"totalRulesInPolicy" json:"totalRulesInPolicy"`
	TotalRulesFailed   int    `yaml:"totalRulesFailed"  json:"totalRulesFailed"`
	TotalPassedCount   int    `yaml:"totalPassedCount"  json:"totalPassedCount"`
}
