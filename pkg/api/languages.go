package deepl

import (
	"encoding/json"
	"net/url"
)

type Language struct {
	Code              string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"`
}

func (t *Translator) GetLanguages(langType string) ([]Language, error) {
	vals := make(url.Values)

	if langType != "" {
		vals.Set("type", langType)
	}

	res, err := t.httpClient.do("GET", "languages", vals)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var languages []Language
	if err := json.NewDecoder(res.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return languages, nil
}
