package deepl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TranslatorOptions struct {
	ServerUrl string `default:""`
}

type Translator struct {
	httpClient *Client
	serverURL  string
	authKey    string
}

func NewTranslator(authKey string, options TranslatorOptions) *Translator {
	var serverURL string
	if options.ServerUrl == "" {
		if authKeyIsFreeAccount(authKey) {
			serverURL = BaseURLFree
		} else {
			serverURL = BaseURLPro
		}
	}

	timeout := 10 * time.Second
	httpClient := NewClient(serverURL, authKey, timeout)

	return &Translator{
		httpClient: httpClient,
		serverURL:  serverURL,
		authKey:    authKey,
	}
}

func (t *Translator) callAPI(method string, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", t.serverURL, endpoint)

	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", t.authKey))
	if _, ok := headers["Content-Type"]; !ok {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := t.httpClient.do(method, url, data, headers)

	return res, err
}

func authKeyIsFreeAccount(authKey string) bool {
	return strings.HasSuffix(authKey, ":fx")
}
