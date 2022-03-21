package cliClient

import (
	"strings"

	"github.com/datreeio/datree/pkg/httpClient"
)

type HTTPClient interface {
	Request(method string, resourceURI string, body interface{}, headers map[string]string) (httpClient.Response, error)
}
type CliClient struct {
	baseUrl       string
	httpClient    HTTPClient
	timeoutClient HTTPClient
	httpErrors    []string
}

func NewCliClient(url string) *CliClient {
	httpClient := httpClient.NewClient(url, nil)
	return &CliClient{
		baseUrl:       url,
		httpClient:    httpClient,
		timeoutClient: nil,
		httpErrors:    []string{},
	}
}

func (c *CliClient) AddHttpError(errStr string) {
	c.httpErrors = append(c.httpErrors, errStr)
}

func (c *CliClient) IsBackendAvailable() bool {
	for _, httpError := range c.httpErrors {
		if strings.Contains(httpError, "connection refused") || strings.Contains(httpError, "ECONNREFUSED") {
			return false
		}
	}
	return true
}
