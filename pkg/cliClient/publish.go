package cliClient

import (
	"github.com/datreeio/datree/bl/files"
	"net/http"
)

func (c *CliClient) PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) error {
	_, err := c.httpClient.Request(http.MethodPut, "/cli/policy/publish/cliIds/"+cliId, policiesConfiguration, nil)
	return err
}
