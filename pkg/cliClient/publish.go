package cliClient

import (
	"encoding/json"
	"github.com/datreeio/datree/bl/files"
	"net/http"
)

type PublishResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Errors       []string `json:"errors"`
}

func (c *CliClient) PublishPolicies(policiesConfiguration files.UnknownStruct, cliId string) (*PublishResponse, error) {
	res, err := c.httpClient.Request(http.MethodPost, "/cli/policy/accounts/"+cliId+"/publish/policies", policiesConfiguration, nil)
	if err != nil {
		return nil, err
	}

	var response = &PublishResponse{}
	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
