package deepl

import (
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
