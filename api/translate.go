package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Translation holds the results of a text translation request.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

// TranslateText translates the given text(s) into the specified target language.
//
// The total request body size must not exceed 128 KiB (128 · 1024 bytes).
func (t *Translator) TranslateText(text []string, targetLang string, opts ...TranslateOption) ([]Translation, error) {
	const (
		endpoint string = "translate"
		method   string = http.MethodPost
	)

	data := struct {
		Text       []string `json:"text"`
		TargetLang string   `json:"target_lang"`

		TranslateOptions
	}{
		Text:       text,
		TargetLang: targetLang,
	}
	if err := data.TranslateOptions.Gather(opts...); err != nil {
		return nil, fmt.Errorf("error setting translate option: %w", err)
	}

	// Setup request
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding request data: %w", err)
	}

	// Send request
	res, err := t.callAPI(method, endpoint, headers, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	// Parse response
	var response struct {
		Translations []Translation `json:"translations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Translations, nil
}
