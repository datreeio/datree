package cliClient

import (
	"net/http"

	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/pkg/httpClient"
)

func (c *CliClient) PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) (httpClient.Response, error) {
	headers := map[string]string{"x-cli-id": cliId}
	return c.httpClient.Request(http.MethodPut, "/cli/policy/publish", policiesConfiguration, headers)
}
