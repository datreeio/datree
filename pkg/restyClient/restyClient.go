package restyClient

import (
	"github.com/go-resty/resty/v2"
)

func New(baseUrl string) *resty.Client {
	client := resty.New()
	client.SetBaseURL(baseUrl)
	return client
}
