package cliClient

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/datreeio/datree/pkg/httpClient"
)

type PublishReturnType struct {
	CliVersion   string `json:"cliVersion"`
	MessageText  string `json:"messageText"`
	MessageColor string `json:"messageColor"`
}

func (c *CliClient) PublishPolicies(policiesJson string, timeout int) (*PublishReturnType, error) {
	_timeout := time.Duration(timeout) * time.Millisecond

	var timeoutClient HTTPClient
	if c.timeoutClient != nil {
		timeoutClient = c.timeoutClient
	} else {
		timeoutClient = httpClient.NewClientTimeout(c.baseUrl, nil, _timeout)
	}

	res, err := timeoutClient.Request(http.MethodPost, "/cli/policy/publish", policiesJson, nil)
	if err != nil {
		return nil, err
	}

	var response = &PublishReturnType{}
	err = json.Unmarshal(res.Body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
