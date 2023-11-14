package deepl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ServerURLPro  = "https://api.deepl.com/v2"
	ServerURLFree = "https://api-free.deepl.com/v2"
)

type Translator struct {
	client    HTTPClient
	serverURL string
	authKey   string
}

// TranslatorOption is a functional option for configuring the Translator
type TranslatorOption func(*Translator) error

// WithServerURL allows overriding the default server url
func WithServerURL(url string) TranslatorOption {
	return func(t *Translator) error {
		t.serverURL = url
		return nil
	}
}

// WithHTTPClient allows overriding the default http client
func WithHTTPClient(c HTTPClient) TranslatorOption {
	return func(t *Translator) error {
		t.client = c
		return nil
	}
}

// NewTranslator creates a new translator
func NewTranslator(authKey string, opts ...TranslatorOption) (*Translator, error) {
	// Determine default server url based on auth key
	var serverURL string
	if authKeyIsFreeAccount(authKey) {
		serverURL = ServerURLFree
	} else {
		serverURL = ServerURLPro
	}

	// Set up with default http client
	t := &Translator{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		serverURL: serverURL,
		authKey:   authKey,
	}

	if err := t.applyOptions(opts...); err != nil {
		return nil, err
	}

	return t, nil
}

// applyOptions apllies the supplied functional options to the Translator
func (t *Translator) applyOptions(opts ...TranslatorOption) error {
	for _, opt := range opts {
		err := opt(t)
		if err != nil {
			return err
		}
	}

	return nil
}

// callAPI calls the supplied API endpoint with the provided parameters and returns the response
func (t *Translator) callAPI(method string, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", t.serverURL, endpoint)

	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", t.authKey))
	if _, ok := headers["Content-Type"]; !ok {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	return t.client.Do(req)
}

// authKeyIsFreeAccount determines whether the supplied auth key belongs to a Free account
func authKeyIsFreeAccount(authKey string) bool {
	return strings.HasSuffix(authKey, ":fx")
}
