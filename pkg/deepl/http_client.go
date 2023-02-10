package deepl

import (
	"fmt"
	"net/http"
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

func (c *Client) do(method string, endpoint string, params map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.apiKey))

	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return c.httpClient.Do(req)
}
