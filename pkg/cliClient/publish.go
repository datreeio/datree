package cliClient

import (
	"errors"
	"github.com/datreeio/datree/bl/files"
	"net/http"
)

func (c *CliClient) PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) error {
	_, err := c.httpClient.Request(http.MethodPost, "/cli/policy/accounts/"+cliId+"/publish/policies", policiesConfiguration, nil)
	if err != nil && err.Error() == "<nil>" {
		return errors.New("unknown error")
	} else {
		return err
	}
}
