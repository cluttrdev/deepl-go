package deepl

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func httpError(statusCode int) error {
	var statusText string
	switch statusCode {
	case 456:
		statusText = "Quota exceeded. The character limit has been reached."
	default:
		statusText = http.StatusText(statusCode)
	}

	return fmt.Errorf("%d - %s", statusCode, statusText)
}

var ErrorStatusTooManyRequests = errors.New(httpError(http.StatusTooManyRequests).Error())
var ErrorStatusInternalServerError = errors.New(httpError(http.StatusInternalServerError).Error())

func retriableHTTPError(statusCode int) error {
	switch {
	case statusCode == http.StatusTooManyRequests:
		return ErrorStatusTooManyRequests
	case statusCode >= http.StatusInternalServerError:
		return ErrorStatusInternalServerError
	}
	return nil
}
