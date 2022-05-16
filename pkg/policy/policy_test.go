package policy

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/defaultRules"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

type TestFilesByRuleId = map[int]*FailAndPassTests

type FailAndPassTests struct {
	fails  []*FileWithName
	passes []*FileWithName
}

type FileWithName struct {
	name    string
	content string
}

func TestGetPoliciesFileFromPath(t *testing.T) {
	policiesYamlPath := "../../internal/fixtures/policyAsCode/policies.yaml"
	policies, err := GetPoliciesFileFromPath(policiesYamlPath)
	if err != nil {
		panic(err)
	}

	expectedPoliciesJson := expectedPoliciesContent(t, policiesYamlPath)
	assert.True(t, reflect.DeepEqual(policies, expectedPoliciesJson))
}

func TestDefaultRulesValidation(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	defaultRules, err := defaultRules.GetDefaultRules()
	if err != nil {
		panic(err)
	}

	testFilesByRuleId := getTestFilesByRuleId(t)
	validator := jsonSchemaValidator.New()

	for _, rule := range defaultRules.Rules {
		validatePassing(t, validator, rule.Schema, rule.ID, testFilesByRuleId[rule.ID].passes, true)
		validatePassing(t, validator, rule.Schema, rule.ID, testFilesByRuleId[rule.ID].fails, false)
	}
}

func validatePassing(t *testing.T, validator *jsonSchemaValidator.JSONSchemaValidator, schemaContent map[string]interface{}, ruleId int, files []*FileWithName, expectPass bool) {
	for _, file := range files {
		schemaBytes, err := yaml.Marshal(schemaContent)
		if err != nil {
			panic(err)
		}

		errorsResult, err := validator.ValidateYamlSchema(string(schemaBytes), file.content)
		if err != nil {
			panic(errors.New(err.Error() + fmt.Sprintf("\nruleId: %d", ruleId)))
		}

		if len(errorsResult) > 0 && expectPass {
			t.Errorf("Expected validation for rule with id %d to pass, but it failed for file %s\n", ruleId, file.name)
		}
		if len(errorsResult) == 0 && !expectPass {
			t.Errorf("Expected validation for rule with id %d to fail, but it passed for file %s\n", ruleId, file.name)
		}
	}
}

func getTestFilesByRuleId(t *testing.T) TestFilesByRuleId {
	dirPath := "./pkg/policy/tests"
	fileReader := fileReader.CreateFileReader(nil)
	files, err := fileReader.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	testFilesByRuleId := make(TestFilesByRuleId)
	for _, file := range files {
		filename, err := fileReader.GetFilename(file)
		if err != nil {
			panic(err)
		}

		fileContent, err := fileReader.ReadFileContent(file)
		if err != nil {
			panic(err)
		}

		id, isPass := getFileData(filename)
		if testFilesByRuleId[id] == nil {
			testFilesByRuleId[id] = &FailAndPassTests{}
		}

		fileWithName := &FileWithName{name: filename, content: fileContent}
		if isPass {
			testFilesByRuleId[id].passes = append(testFilesByRuleId[id].passes, fileWithName)
		} else {
			testFilesByRuleId[id].fails = append(testFilesByRuleId[id].fails, fileWithName)
		}
	}

	return testFilesByRuleId
}

func getFileData(filename string) (int, bool) {
	parts := strings.Split(filename, "-")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

	isPass := strings.Contains(parts[1], "pass")
	return id, isPass
}

func expectedPoliciesContent(t *testing.T, path string) *cliClient.EvaluationPrerunPolicies {
	fileReader := fileReader.CreateFileReader(nil)
	policiesStr, _ := fileReader.ReadFileContent(path)

	var policiesJson *cliClient.EvaluationPrerunPolicies
	policiesBytes, _ := yaml.YAMLToJSON([]byte(policiesStr))

	err := yaml.Unmarshal(policiesBytes, &policiesJson)
	if err != nil {
		panic(err)
	}
	return policiesJson
}
