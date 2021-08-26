package cliClient

import (
	"github.com/datreeio/datree/bl/files"
	"net/http"
)

func (c *CliClient) PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) error {
	headers := make(map[string]string)
	headers["x-cli-id"] = cliId
	_, err := c.httpClient.Request(http.MethodPut, "/cli/policy/publish", policiesConfiguration, headers)
	return err
}
