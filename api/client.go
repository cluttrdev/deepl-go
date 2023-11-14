package deepl

import (
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPError struct {
	StatusCode int
}

func (err HTTPError) Error() string {
	switch err.StatusCode {
	case 456:
		return fmt.Sprintf("%d - %s", err.StatusCode, "Quota exceeded. The character limit has been reached.")
	default:
		return fmt.Sprintf("%d - %s", err.StatusCode, http.StatusText(err.StatusCode))
	}
}
