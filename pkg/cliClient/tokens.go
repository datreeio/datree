package cliClient

import (
	"encoding/json"
	"net/http"
)

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (c *CliClient) CreateToken() (*CreateTokenResponse, error) {
	if c.networkValidator.IsLocalMode() {
		return nil, nil
	}

	headers := map[string]string{}
	res, err := c.httpClient.Request(http.MethodPost, "/cli/tokens/", nil, headers)

	if err != nil {
		validatorErr := c.networkValidator.SetIsBackendAvailable(err.Error())
		if validatorErr != nil {
			return nil, validatorErr
		}

		if c.networkValidator.IsLocalMode() {
			return nil, nil
		}

		return nil, err
	}

	createTokenResponse := &CreateTokenResponse{}
	err = json.Unmarshal(res.Body, &createTokenResponse)

	if err != nil {
		return nil, err
	}

	return createTokenResponse, nil
}
