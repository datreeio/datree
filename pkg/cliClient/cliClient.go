package cliClient

import (
	"github.com/datreeio/datree/pkg/httpClient"
)

type HTTPClient interface {
	Request(method string, resourceURI string, body interface{}, headers map[string]string) (httpClient.Response, error)
}
type NetworkValidator interface {
	SetIsBackendAvailable(errStr string)
	IsBackendAvailable() bool
}

type CliClient struct {
	baseUrl          string
	httpClient       HTTPClient
	timeoutClient    HTTPClient
	httpErrors       []string
	networkValidator NetworkValidator
}

func NewCliClient(url string, networkValidator NetworkValidator) *CliClient {
	httpClient := httpClient.NewClient(url, nil)
	return &CliClient{
		baseUrl:          url,
		httpClient:       httpClient,
		timeoutClient:    nil,
		httpErrors:       []string{},
		networkValidator: networkValidator,
	}
}

func (c *CliClient) SetIsBackendAvailable(errStr string) {
	c.networkValidator.SetIsBackendAvailable(errStr)
}

func (c *CliClient) IsBackendAvailable() bool {
	return c.networkValidator.IsBackendAvailable()
}
