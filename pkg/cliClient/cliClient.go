package cliClient

import (
	"strings"

	"github.com/datreeio/datree/pkg/httpClient"
)

type HTTPClient interface {
	Request(method string, resourceURI string, body interface{}, headers map[string]string) (httpClient.Response, error)
}
type CliClient struct {
	baseUrl            string
	httpClient         HTTPClient
	timeoutClient      HTTPClient
	httpErrors         []string
	isBackendAvailable bool
}

func NewCliClient(url string) *CliClient {
	httpClient := httpClient.NewClient(url, nil)
	return &CliClient{
		baseUrl:            url,
		httpClient:         httpClient,
		timeoutClient:      nil,
		httpErrors:         []string{},
		isBackendAvailable: true,
	}
}

func (c *CliClient) SetIsBackendAvailable(errStr string) {
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "ECONNREFUSED") {
		c.isBackendAvailable = false
	}
}

func (c *CliClient) IsBackendAvailable() bool {
	return c.isBackendAvailable
}
