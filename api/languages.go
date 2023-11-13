package deepl

import (
	"encoding/json"
	"net/http"
	"net/url"
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

	if langType != "" {
		vals.Set("type", langType)
	}

	res, err := t.callAPI("GET", "languages", vals, nil)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}
	defer res.Body.Close()

	var languages []Language
	if err := json.NewDecoder(res.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return languages, nil
}

func (t *Translator) GetGlossaryLanguagePairs() ([]LanguagePair, error) {
	res, err := t.callAPI("GET", "glossary-language-pairs", nil, nil)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}
	defer res.Body.Close()

	var response struct {
		SupportedLanguages []LanguagePair `json:"supported_languages"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.SupportedLanguages, nil
}
