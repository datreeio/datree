package utils

import (
	"fmt"
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

func IsNetworkError(errStr string) bool {
	networkErrors := []string{"network error", "connection refused", "no such host", "i/o timeout", "server misbehaving"}
	return stringInSliceContains(errStr, networkErrors)
}

func stringInSliceContains(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}
