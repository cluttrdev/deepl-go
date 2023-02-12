package deepl

import (
	"strings"
	"time"
)

type TranslatorOptions struct {
	ServerUrl string `default:""`
}

type Translator struct {
	httpClient *Client
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
	}
}

func authKeyIsFreeAccount(authKey string) bool {
	return strings.HasSuffix(authKey, ":fx")
}
