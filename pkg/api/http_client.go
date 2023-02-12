package deepl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURLPro  = "https://api.deepl.com/v2"
	BaseURLFree = "https://api-free.deepl.com/v2"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	authKey    string
}

func NewClient(baseURL string, authKey string, timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		httpClient: client,
		baseURL:    baseURL,
		authKey:    authKey,
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

func (c *Client) do(method string, endpoint string, params url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, endpoint), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.authKey))

	res, err := c.httpClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}

	return res, err
}
