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

type (
	Client struct {
		httpClient *http.Client
		baseURL    string
		apiKey     string
	}
)

func NewClient(baseURL string, apiKey string, timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		httpClient: client,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

type Error struct {
	Code int
}

func (err Error) Error() string {
	return http.StatusText(err.Code)
}

func (c *Client) do(method string, endpoint string, params url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, endpoint), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.apiKey))

	res, err := c.httpClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, Error{Code: res.StatusCode}
	}

	return res, err
}
