package main

import (
	"context"
	_ "embed"
	"github.com/ghodss/yaml"
	"log"

	"github.com/open-policy-agent/opa/rego"
)

//go:embed test.rego
var regoRuleFileContent string

//go:embed k8s-demo.yaml
var k8sDemoFileContent string

func main() {
	runRegoRule(regoRuleFileContent, k8sDemoFileContent)
}

func runRegoRule(regoQuery string, yamlFileToTest string) rego.ResultSet {

	ctx := context.Background()

	// Construct a Rego object that can be prepared or evaluated.
	r := rego.New(rego.Query(regoQuery))

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
