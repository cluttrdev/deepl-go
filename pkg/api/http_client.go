package deepl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		httpClient: client,
	}
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

func (c *Client) do(method string, url string, data url.Values, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	res, err := c.httpClient.Do(req)

	return res, err
}
