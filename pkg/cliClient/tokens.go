package cliClient

import (
	"encoding/json"
	"net/http"
)

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (c *CliClient) CreateToken() (*CreateTokenResponse, error) {
	headers := map[string]string{}
	res, err := c.httpClient.Request(http.MethodPost, "/cli/tokens/", nil, headers)

	if err != nil {
		c.AddHttpError(err.Error())
		return nil, err
	}

	createTokenResponse := &CreateTokenResponse{}
	err = json.Unmarshal(res.Body, &createTokenResponse)

	if err != nil {
		return nil, err
	}

	return createTokenResponse, nil
}
