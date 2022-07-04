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
		return &CreateTokenResponse{}, nil
	}

	headers := map[string]string{}
	res, err := c.httpClient.Request(http.MethodPost, "/cli/tokens/", nil, headers)

	if err != nil {
		networkErr := c.networkValidator.IdentifyNetworkError(err)
		if networkErr != nil {
			return nil, networkErr
		}

		if c.networkValidator.IsLocalMode() {
			return &CreateTokenResponse{}, nil
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
