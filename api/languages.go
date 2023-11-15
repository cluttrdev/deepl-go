package deepl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	vals := make(url.Values)

	switch langType {
	case "":
		// default is `source`
	case "source", "target":
		vals.Set("type", langType)
	default:
		return nil, fmt.Errorf("Invalid language `type` value: %v", langType)
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	body := strings.NewReader(vals.Encode())

	res, err := t.callAPI(http.MethodGet, "languages", headers, body)
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
	res, err := t.callAPI(http.MethodGet, "glossary-language-pairs", nil, nil)
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
