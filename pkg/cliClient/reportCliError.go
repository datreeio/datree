package cliClient

import (
	"net/http"
)

type ReportCliErrorRequest struct {
	ClientId     string `json:"clientId"`
	Token        string `json:"token"`
	CliVersion   string `json:"cliVersion"`
	ErrorMessage string `json:"errorMessage"`
	StackTrace   string `json:"stackTrace"`
}

func (c *CliClient) ReportCliError(reportCliErrorRequest ReportCliErrorRequest, uri string) (StatusCode int, Error error) {
	headers := map[string]string{}
	res, err := c.httpClient.Request(
		http.MethodPost,
		"/cli/public"+uri,
		reportCliErrorRequest,
		headers,
	)
	return res.StatusCode, err
}
