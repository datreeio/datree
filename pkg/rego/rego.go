package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/util/test"
	"log"
)

//go:embed k8s-demo.yaml
var k8sDemoFileContent string

//go:embed test.rego
var rule1RegoFileContent string

func main() {
	result := runRegoRule2(rule1RegoFileContent, "data.main.deny", k8sDemoFileContent)
	fmt.Println(result)
}

func runRegoRule2(regoFileContent string, regoQuery string, k8sFileToTestContent string) (result rego.ResultSet) {
	files := map[string]string{
		"/someRandomFileName": regoFileContent,
	}

	test.WithTempFS(files, func(rootDir string) {
		result = runRegoRule(fmt.Sprint(rootDir, "/someRandomFileName"), regoQuery, k8sFileToTestContent)
	})
	return result
}

func runRegoRule(regoFilePath string, regoQuery string, yamlFileToTest string) rego.ResultSet {

	ctx := context.Background()

	// Construct a Rego object that can be prepared or evaluated.
	r := rego.New(
		rego.Query(regoQuery),
		rego.Load([]string{regoFilePath}, nil))
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

	return rs
}

func YAMLToStruct(content string) (res any, err error) {
	var result any
	err = yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

/**
if deny is empty
and violations is empty
then
pass
else
fail, take message from deny/JSON.stringify(violations) if there are any. (fallback to failSuggestion)

regoRules: ./path-to-rego-rules/**.rego
policies:
 - policy1
   rules:
	ID: roy_you_such
    isEnabled: false



deny[res] {
    input.kind == "Deployment"
    res := {"message":"bad${index}", "ruleID": "MY_SPECIAL_RULE_ID"}
}

compile the rego upon publish

fail for invalid deny struct upon run

"data.main.deny"

join + set RuleId


*/
