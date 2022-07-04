package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseErrorToString(err interface{}) string {
	switch panicErr := err.(type) {
	case string:
		return panicErr
	case error:
		return panicErr.Error()
	default:
		return fmt.Sprintf("%v", panicErr)
	}
}

func IsNetworkError(err error) bool {
	networkErrors := []string{"network error", "connection refused", "no such host", "i/o timeout", "server misbehaving"}
	return stringInSliceContains(err.Error(), networkErrors) || isUrlErrorType(err)
}

func stringInSliceContains(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func isUrlErrorType(err error) bool {
	_, ok := err.(*url.Error)
	return ok
}
