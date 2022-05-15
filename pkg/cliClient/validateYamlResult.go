package cliClient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValidatedFile struct {
	Path   string `json:"path"`
	Status bool   `json:"status"`
}

type ValidatedYamlResult struct {
	Token    string           `json:"token"`
	ClientId string           `json:"clientId"`
	Metadata *Metadata        `json:"metadata"`
	Files    []*ValidatedFile `json:"files"`
	Status   bool             `json:"status"`
}

func (c *CliClient) SendValidateYamlResult(request *ValidatedYamlResult) error {
	if c.networkValidator.IsLocalMode() {
		return nil
	}

	httpRes, err := c.httpClient.Request(http.MethodPost, "/cli/yaml-validation/", request, c.flagsHeaders)
	if err != nil {
		networkErr := c.networkValidator.IdentifyNetworkError(err.Error())
		if networkErr != nil {
			return networkErr
		}

		if c.networkValidator.IsLocalMode() {
			return nil
		}

		return err
	}

	var errorJson map[string]interface{}
	err = json.Unmarshal(httpRes.Body, &errorJson)
	fmt.Println(errorJson)
	if err != nil {
		return err
	}

	return nil
}
