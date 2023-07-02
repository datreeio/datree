package policy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/datreeio/datree/pkg/defaultPolicies"

	"github.com/datreeio/datree/pkg/defaultRules"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

type TestFilesByRuleId = map[int]*FailAndPassTests

type FailAndPassTests struct {
	fails  []*FileWithPath
	passes []*FileWithPath
}

type FileWithPath struct {
	path    string
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

func validatePassing(t *testing.T, validator *jsonSchemaValidator.JSONSchemaValidator, schemaContent map[string]interface{}, ruleId int, files []*FileWithPath, expectPass bool) {
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
			t.Errorf("Expected validation for rule with id %d to pass, but it failed for file %s\n", ruleId, file.path)
		}
		if len(errorsResult) == 0 && !expectPass {
			t.Errorf("Expected validation for rule with id %d to fail, but it passed for file %s\n", ruleId, file.path)
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
	for _, filePath := range files {
		// skip directories
		if !fileExists(filePath) {
			continue
		}

		path, _ := filepath.Split(filePath)
		passOrFail := filepath.Base(path)
		ruleId := filepath.Base(filepath.Dir(filepath.Dir(path)))
		id, err := strconv.Atoi(ruleId)
		if err != nil {
			panic(err)
		}

		isPass := passOrFail == "pass"

		fileContent, err := fileReader.ReadFileContent(filePath)
		if err != nil {
			panic(err)
		}

		if testFilesByRuleId[id] == nil {
			testFilesByRuleId[id] = &FailAndPassTests{}
		}

		fileWithPath := &FileWithPath{path: filePath, content: fileContent}
		if isPass {
			testFilesByRuleId[id].passes = append(testFilesByRuleId[id].passes, fileWithPath)
		} else {
			testFilesByRuleId[id].fails = append(testFilesByRuleId[id].fails, fileWithPath)
		}
	}

	return testFilesByRuleId
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func expectedPoliciesContent(t *testing.T, path string) *defaultPolicies.EvaluationPrerunPolicies {
	fileReader := fileReader.CreateFileReader(nil)
	policiesStr, _ := fileReader.ReadFileContent(path)

	var policiesJson *defaultPolicies.EvaluationPrerunPolicies
	policiesBytes, _ := yaml.YAMLToJSON([]byte(policiesStr))

	err := yaml.Unmarshal(policiesBytes, &policiesJson)
	if err != nil {
		panic(err)
	}
	return policiesJson
}
