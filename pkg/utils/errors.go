package utils

import "fmt"

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
