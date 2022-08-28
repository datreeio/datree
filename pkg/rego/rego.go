package rego

import (
	"context"
	_ "embed"
	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/util/test"
	"log"
)

//go:embed k8s-demo.yaml
var k8sDemoFileContent string

//go:embed test.rego
var rule1RegoFileContent string

type DenyItem struct {
	message string `json:"message"`
	ruleID  string `json:"ruleID"`
}

type DenyArray []DenyItem

func GetRegoDenyArray(regoRulesFiles *FilesAsStruct, configurationJson string) (denyArray DenyArray) {
	var paths []string
	for k := range *regoRulesFiles {
		paths = append(paths, k)
	}

	test.WithTempFS(*regoRulesFiles, func(rootDir string) {
		denyArray = runRegoRule(paths, configurationJson)
	}
	return denyArray
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
	
	s, ok := rs.(any)
	if s.deny != nil {
		return s.deny
	} else {
		return DenyArray{}
	}
}

func YAMLToStruct(content string) (res any, err error) {
	var result any
	err = yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
