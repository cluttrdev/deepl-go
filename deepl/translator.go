package deepl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cluttrdev/deepl-go/internal/retry"
)

const (
	ServerURLPro  = "https://api.deepl.com"
	ServerURLFree = "https://api-free.deepl.com"
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
	if isFreeAccountAuthKey(authKey) {
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
func (t *Translator) callAPI(method string, endpoint string, headers http.Header, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", t.serverURL, endpoint)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", t.authKey))
	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}

	res, err := retry.DoWithData(
		func() (*http.Response, error) {
			res, err := t.client.Do(req)
			if err != nil {
				return res, err
			} else if err := retriableHTTPError(res.StatusCode); err != nil {
				return res, err
			}
			return res, nil
		},
		retry.RetryIf(func(err error) bool {
			return isRetriableHTTPError(err)
		}),
		retry.MaxAttempts(5),
		retry.WithBackoff(&retry.Backoff{
			InitialDelay: 1 * time.Second,
			MaxDelay:     120 * time.Second,
			Factor:       1.6,
			Jitter:       0.23,
		}),
	)

	return res, err
}

func isRetriableHTTPError(err error) bool {
	switch {
	case errors.Is(err, ErrorStatusTooManyRequests):
		return true
	case errors.Is(err, ErrorStatusInternalServerError):
		return true
	}
	return false
}

// isFreeAccountAuthKey determines whether the supplied auth key belongs to a Free account
func isFreeAccountAuthKey(authKey string) bool {
	return strings.HasSuffix(authKey, ":fx")
}
