package evaluation

import (
	"github.com/datreeio/datree/pkg/extractor"
)

type FormattedOutput struct {
	PolicyValidationResults []*FormattedEvaluationResults   `yaml:"policyValidationResults" json:"policyValidationResults" xml:"policyValidationResults"`
	PolicySummary           *PolicySummary                  `yaml:"policySummary" json:"policySummary" xml:"policySummary"`
	EvaluationSummary       NonInteractiveEvaluationSummary `yaml:"evaluationSummary" json:"evaluationSummary" xml:"evaluationSummary"`
	YamlValidationResults   []*extractor.InvalidFile        `yaml:"yamlValidationResults" json:"yamlValidationResults" xml:"yamlValidationResults"`
	K8sValidationResults    []*extractor.InvalidFile        `yaml:"k8sValidationResults" json:"k8sValidationResults" xml:"k8sValidationResults"`
}

type NonInteractiveEvaluationResults struct {
	FormattedEvaluationResults []*FormattedEvaluationResults
	PolicySummary              *PolicySummary
}

type FormattedEvaluationResults struct {
	FileName    string        `yaml:"fileName" json:"fileName" xml:"fileName"`
	RuleResults []*RuleResult `yaml:"ruleResults" json:"ruleResults" xml:"ruleResults"`
}

type RuleResult struct {
	Identifier         string              `yaml:"identifier" json:"identifier" xml:"identifier"`
	Name               string              `yaml:"name" json:"name" xml:"name"`
	MessageOnFailure   string              `yaml:"messageOnFailure" json:"messageOnFailure" xml:"messageOnFailure"`
	OccurrencesDetails []OccurrenceDetails `yaml:"occurrencesDetails" json:"occurrencesDetails" xml:"occurrencesDetails"`
}

type NonInteractiveEvaluationSummary struct {
	ConfigsCount                int    `yaml:"configsCount" json:"configsCount" xml:"configsCount"`
	FilesCount                  int    `yaml:"filesCount" json:"filesCount" xml:"filesCount"`
	PassedYamlValidationCount   int    `yaml:"passedYamlValidationCount" json:"passedYamlValidationCount" xml:"passedYamlValidationCount"`
	K8sValidation               string `yaml:"k8sValidation" json:"k8sValidation" xml:"k8sValidation"`
	PassedPolicyValidationCount int    `yaml:"passedPolicyValidationCount" json:"passedPolicyValidationCount" xml:"passedPolicyValidationCount"`
}

type PolicySummary struct {
	PolicyName         string `yaml:"policyName" json:"policyName" xml:"policyName"`
	TotalRulesInPolicy int    `yaml:"totalRulesInPolicy" json:"totalRulesInPolicy" xml:"totalRulesInPolicy"`
	TotalSkippedRules  int    `yaml:"totalSkippedRules" json:"totalSkippedRules" xml:"totalSkippedRules"`
	TotalRulesFailed   int    `yaml:"totalRulesFailed"  json:"totalRulesFailed" xml:"totalRulesFailed"`
	TotalPassedCount   int    `yaml:"totalPassedCount"  json:"totalPassedCount" xml:"totalPassedCount"`
}
