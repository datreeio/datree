package cliClient

import (
	"github.com/datreeio/datree/pkg/httpClient"
)

type HTTPClient interface {
	Request(method string, resourceURI string, body interface{}, headers map[string]string) (httpClient.Response, error)
}
type CliClient struct {
	baseUrl       string
	httpClient    HTTPClient
	timeoutClient HTTPClient
}

func NewCliClient(url string) *CliClient {
	httpClient := httpClient.NewClient(url, nil)
	return &CliClient{
		baseUrl:       url,
		httpClient:    httpClient,
		timeoutClient: nil,
	}
}
