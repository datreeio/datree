package cliClient

import (
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
	Status   string           `json:"status"`
}

func (c *CliClient) SendValidateYamlResult(request *ValidatedYamlResult) {
	if c.networkValidator.IsLocalMode() {
		return
	}

	c.httpClient.Request(http.MethodPost, "/cli/yaml-validation/", request, c.flagsHeaders)
}
