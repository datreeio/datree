package ciContext

import (
	"fmt"
	"os"
	"strings"
)

type CIContext struct {
	IsCI  bool   `json:"isCI"`
	CIEnv string `json:"ciEnv"`
}

var ciMap = map[string]string{
	"GITHUB_ACTIONS":                     "github_actions",
	"GITLAB_CI":                          "gitlab_ci",
	"CIRCLECI":                           "circle-ci",
	"JENKINS_HOME":                       "jenkins",
	"JENKINS_URL":                        "jenkins",
	"BUILDKITE":                          "buildkite",
	"SYSTEM_COLLECTIONURI":               fmt.Sprintf("azure_devops_%s", os.Getenv("BUILD_REPOSITORY_PROVIDER")),
	"SYSTEM_TEAMFOUNDATIONCOLLECTIONURI": "azure-pipelines",
	"TFC_RUN_ID":                         "tfc",
	"ENV0_ENVIRONMENT_ID":                "env0",
	"CF_BUILD_ID":                        "codefresh",
	"TRAVIS":                             "travis",
	"CODEBUILD_CI":                       "codebuild",
	"TEAMCITY_VERSION":                   "teamcity",
	"BUDDYBUILD_BRANCH":                  "buddybuild",
	"BUDDY_WORKSPACE_ID":                 "buddy",
	"APPVEYOR":                           "appveyor",
	"WERCKER_GIT_BRANCH":                 "wercker",
	"WERCKER":                            "wercker",
	"SHIPPABLE":                          "shippable",
	"BITBUCKET_BUILD_NUMBER":             "bitbucket-pipelines",
	"CIRRUS_CI":                          "cirrusci",
	"DRONE":                              "drone",
	"GO_PIPELINE_NAME":                   "gocd",
	"SAIL_CI":                            "sail",
}

var ciMapPrefix = map[string]string{
	"ATLANTIS_":  "atlantis",
	"BITBUCKET_": "bitbucket",
	"CONCOURSE_": "concourse",
	"SPACELIFT_": "spacelift",
	"HARNESS_":   "harness",
}

func Extract() *CIContext {
	var checkCIEnv = isCI()
	return &CIContext{
		IsCI:  checkCIEnv != "none",
		CIEnv: checkCIEnv,
	}
}

func isCI() string {
	for env, name := range ciMap {
		if isKeyInEnv(env) {
			return name
		}
	}

	for _, key := range os.Environ() {
		for prefix, name := range ciMapPrefix {
			if strings.HasPrefix(key, prefix) {
				return name
			}
		}
	}

	return "none"
}

func isKeyInEnv(key string) bool {
	_, present := os.LookupEnv(key)
	return present
}
