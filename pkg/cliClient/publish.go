package cliClient

import (
	"encoding/json"
	"net/http"

	"github.com/datreeio/datree/bl/files"
)

type PublishFailedResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Payload []string `json:"payload"`
}

func (c *CliClient) PublishPolicies(policiesConfiguration files.UnknownStruct, token string) (*PublishFailedResponse, error) {
	headers := map[string]string{"x-datree-token": token}
	res, publishErr := c.httpClient.Request(http.MethodPut, "/cli/policy/publish", policiesConfiguration, headers)
	if publishErr != nil {
		if res.StatusCode != 0 {
			publishFailedResponse := &PublishFailedResponse{}
			err := json.Unmarshal(res.Body, publishFailedResponse)
			if err != nil {
				return nil, publishErr
			}
			return publishFailedResponse, publishErr
		}
		return nil, publishErr
	}
	return nil, nil
}
