package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Language struct {
	Code              string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"`
}

type LanguagePair struct {
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

func (t *Translator) GetLanguages(langType string) ([]Language, error) {
	const (
		endpoint string = "v2/languages"
		method   string = http.MethodGet
	)

	opts := struct {
		Type string `json:"type,omitempty"`
	}{}

	switch langType {
	case "":
		// if omitted, default is `source`
	case "source", "target":
		opts.Type = langType
	default:
		return nil, fmt.Errorf("Invalid languages `type` value: %v", langType)
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("error encoding request data: %w", err)
	}

	res, err := t.callAPI(method, endpoint, headers, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	var languages []Language
	if err := json.NewDecoder(res.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return languages, nil
}

func (t *Translator) GetGlossaryLanguagePairs() ([]LanguagePair, error) {
	const (
		endpoint string = "v2/glossary-language-pairs"
		method   string = http.MethodGet
	)

	res, err := t.callAPI(method, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	var response struct {
		SupportedLanguages []LanguagePair `json:"supported_languages"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.SupportedLanguages, nil
}
