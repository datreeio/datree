package rego

import (
	"context"
	_ "embed"
	"github.com/datreeio/datree/pkg/utils"
	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/util/test"
	"log"
)

type DenyItem struct {
	message string `json:"message"`
	ruleID  string `json:"ruleID"`
}

type DenyArray []DenyItem

type RegoRulesResults map[string]string

func GetRegoDenyArray(regoRulesFiles *FilesAsStruct, configurationJson string) (regoRulesResults RegoRulesResults) {
	var paths []string
	for k := range *regoRulesFiles {
		paths = append(paths, k)
	}

	denyArray := make(DenyArray, 0)
	test.WithTempFS(*regoRulesFiles, func(rootDir string) {
		denyArray = runRegoRule(paths, configurationJson)
	})

	regoRulesResults = make(RegoRulesResults)
	for _, denyItem := range denyArray {
		if regoRulesResults[denyItem.ruleID] == "" {
			regoRulesResults[denyItem.ruleID] = denyItem.message
		} else {
			regoRulesResults[denyItem.ruleID] = regoRulesResults[denyItem.ruleID] + ", " + denyItem.message
		}
	}

	return regoRulesResults
}

var pathToQuery = "data.main.deny"

func runRegoRule(regoFilePaths []string, yamlFileToTest string) DenyArray {
	ctx := context.Background()

	// Construct a Rego object that can be prepared or evaluated.
	r := rego.New(
		rego.Query(pathToQuery),
		rego.Load(regoFilePaths, nil))
	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Load the input document from k8sDemoFileContent string
	input, err := YAMLToStruct(yamlFileToTest)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the prepared query.
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Fatal(err)
	}

	actualResult := rs[0].Expressions[0].Value
	rawDenyArray, ok := actualResult.([]any)
	if !ok {
		log.Fatal("Error: could not convert result to DenyArray")
	}

	denyArray := utils.MapSlice(rawDenyArray, func(denyItem any) DenyItem {
		denyItemConverted, ok := denyItem.(map[string]any)
		if !ok {
			log.Fatal("Error: could not convert result to DenyItem")
		}

		itemMessage, ok1 := denyItemConverted["message"].(string)
		itemRuleID, ok2 := denyItemConverted["ruleID"].(string)

		if !ok1 || !ok2 {
			log.Fatal("Error: could not convert result to DenyItem")
		}

		return DenyItem{
			message: itemMessage,
			ruleID:  itemRuleID,
		}
	})

	return denyArray
}

func YAMLToStruct(content string) (res any, err error) {
	var result any
	err = yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
