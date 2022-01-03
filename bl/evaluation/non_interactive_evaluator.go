package evaluation

import (
	"github.com/datreeio/datree/bl/validation"
)

type FormattedOutput struct {
	PolicyValidationResults []*FormattedEvaluationResults
	PolicySummary           PolicySummary
	EvaluationSummary       struct {
		ConfigsCount              int
		FilesCount                int
		PassedYamlValidationCount int
		PassedK8sValidationCount  int
		PassedPolicyCheckCount    int
	}
	YamlValidationResults []*validation.InvalidYamlFile
	K8sValidationResults  []*validation.InvalidK8sFile
}

type NonInteractiveEvaluationResults struct {
	FormattedEvaluationResults []*FormattedEvaluationResults
	PolicySummary              PolicySummary
}

type FormattedEvaluationResults struct {
	FileName     string
	RuleRresults []*RulerObject
}

type RulerObject struct {
	Identifier         string
	Name               string
	MessageOnFailure   string
	OccurrencesDetails []OccurrenceDetails
}

type PolicySummary struct {
	PolicyName         string
	TotalRulesInPolicy int
	TotalRulesFailed   int
	TotalPassedCount   int
}
