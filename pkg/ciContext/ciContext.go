package ciContext

import (
	"fmt"
	"os"
	"strings"
)

type CIContext struct {
	IsCI       bool        `json:"isCI"`
	CIMetadata *CIMetadata `json:"ciMetadata"`
}
type CIMetadata struct {
	CIEnvValue       string `json:"ciEnvValue"`
	ShouldHideEmojis bool   `json:"shouldHideEmojis"`
}

var ciMap = map[string]*CIMetadata{
	"GITHUB_ACTIONS":                     {CIEnvValue: "github_actions", ShouldHideEmojis: false},
	"GITLAB_CI":                          {CIEnvValue: "gitlab_ci", ShouldHideEmojis: false},
	"CIRCLECI":                           {CIEnvValue: "circle-ci", ShouldHideEmojis: false},
	"JENKINS_HOME":                       {CIEnvValue: "jenkins", ShouldHideEmojis: true},
	"JENKINS_URL":                        {CIEnvValue: "jenkins", ShouldHideEmojis: true},
	"BUILDKITE":                          {CIEnvValue: "buildkite", ShouldHideEmojis: false},
	"SYSTEM_COLLECTIONURI":               {CIEnvValue: fmt.Sprintf("azure_devops_%s", os.Getenv("BUILD_REPOSITORY_PROVIDER")), ShouldHideEmojis: false},
	"SYSTEM_TEAMFOUNDATIONCOLLECTIONURI": {CIEnvValue: "azure-pipelines", ShouldHideEmojis: false},
	"TFC_RUN_ID":                         {CIEnvValue: "tfc", ShouldHideEmojis: false},
	"ENV0_ENVIRONMENT_ID":                {CIEnvValue: "env0", ShouldHideEmojis: false},
	"CF_BUILD_ID":                        {CIEnvValue: "codefresh", ShouldHideEmojis: false},
	"TRAVIS":                             {CIEnvValue: "travis", ShouldHideEmojis: true},
	"CODEBUILD_CI":                       {CIEnvValue: "codebuild", ShouldHideEmojis: false},
	"TEAMCITY_VERSION":                   {CIEnvValue: "teamcity", ShouldHideEmojis: false},
	"BUDDYBUILD_BRANCH":                  {CIEnvValue: "buddybuild", ShouldHideEmojis: false},
	"BUDDY_WORKSPACE_ID":                 {CIEnvValue: "buddy", ShouldHideEmojis: false},
	"APPVEYOR":                           {CIEnvValue: "appveyor", ShouldHideEmojis: false},
	"WERCKER_GIT_BRANCH":                 {CIEnvValue: "wercker", ShouldHideEmojis: false},
	"WERCKER":                            {CIEnvValue: "wercker", ShouldHideEmojis: false},
	"SHIPPABLE":                          {CIEnvValue: "shippable", ShouldHideEmojis: false},
	"BITBUCKET_BUILD_NUMBER":             {CIEnvValue: "bitbucket-pipelines", ShouldHideEmojis: false},
	"CIRRUS_CI":                          {CIEnvValue: "cirrusci", ShouldHideEmojis: false},
	"DRONE":                              {CIEnvValue: "drone", ShouldHideEmojis: false},
	"GO_PIPELINE_NAME":                   {CIEnvValue: "gocd", ShouldHideEmojis: false},
	"SAIL_CI":                            {CIEnvValue: "sail", ShouldHideEmojis: false},
}

var ciMapPrefix = map[string]CIMetadata{
	"ATLANTIS_":  {CIEnvValue: "atlantis", ShouldHideEmojis: false},
	"BITBUCKET_": {CIEnvValue: "bitbucket", ShouldHideEmojis: false},
	"CONCOURSE_": {CIEnvValue: "concourse", ShouldHideEmojis: false},
	"SPACELIFT_": {CIEnvValue: "spacelift", ShouldHideEmojis: false},
	"HARNESS_":   {CIEnvValue: "harness", ShouldHideEmojis: false},
}

func Extract() *CIContext {
	ciName := getCIName()
	isCI := len(ciName) > 0
	if isCI {
		return &CIContext{
			IsCI:       isCI,
			CIMetadata: ciMap[ciName],
		}
	} else {
		return &CIContext{
			IsCI: false,
			CIMetadata: &CIMetadata{
				CIEnvValue:       "",
				ShouldHideEmojis: false,
			},
		}
	}
}

func getCIName() string {
	for env := range ciMap {
		if isKeyInEnv(env) {
			return env
		}
	}

	for _, key := range os.Environ() {
		for prefix := range ciMapPrefix {
			if strings.HasPrefix(key, prefix) {
				return prefix
			}
		}
	}

	return ""
}

func isKeyInEnv(key string) bool {
	_, present := os.LookupEnv(key)
	return present
}
