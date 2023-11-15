package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Translation holds the results of a text translation request.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

// TranslateText translates the given text(s) into the specified target language.
//
// The total request body size must not exceed 128 KiB (128 · 1024 bytes).
func (t *Translator) TranslateText(text []string, targetLang string, opts ...TranslateOptionFunc) ([]Translation, error) {
	// Gather translation parameter options
	options := TranslateOptions{}
	for _, optfunc := range opts {
		if err := optfunc(options); err != nil {
			return nil, fmt.Errorf("error setting translate option: %w", err)
		}
	}

	// Setup request
	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := urlEncode(text, targetLang, options)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %w", err)
	}

	// Send request
	res, err := t.callAPI(http.MethodPost, "translate", headers, bytes.NewReader(body))
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}
	defer res.Body.Close()

	// Parse response
	var response struct {
		Translations []Translation `json:"translations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Translations, nil
}

func urlEncode(text []string, targetLang string, opts map[string]string) ([]byte, error) {
	vals := make(url.Values)

	for _, t := range text {
		vals.Add("text", t)
	}
	vals.Set("target_lang", targetLang)

	for name, value := range opts {
		vals.Set(name, value)
	}

	return []byte(vals.Encode()), nil
}

func jsonEncode(text []string, targetLang string, opts map[string]string) ([]byte, error) {
	data := make(map[string]interface{})

	data["text"] = text
	data["target_lang"] = targetLang

	for name, value := range opts {
		switch name {
		case "non_splitting_tags", "splitting_tags", "ignore_tags":
			data[name] = strings.Split(value, ",")
		default:
			data[name] = value
		}
	}

	return json.Marshal(data)
}
