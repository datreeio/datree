package policy

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed validatePoliciesYamlFixtures/customRulesNull.yaml
var customRulesNull string

func TestCustomRulesNull(t *testing.T) {
	err := validatePoliciesYaml([]byte(customRulesNull), "./validatePoliciesYamlFixtures/customRulesNull.yaml")
	assert.Nil(t, err)
}
