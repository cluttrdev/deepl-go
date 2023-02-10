package deepl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TranslateOption func(url.Values)

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type translationResponse struct {
	Translations []Translation `json:"translations"`
}

type Error struct {
	Code int
}

func (err Error) Error() string {
	return http.StatusText(err.Code)
}

func (c *Client) TranslateText(texts []string, targetLang string, options ...TranslateOption) ([]Translation, error) {
	vals := make(url.Values)

	for _, text := range texts {
		vals.Add("text", text)
	}

	vals.Set("target_lang", targetLang)

	for _, opt := range options {
		opt(vals)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.baseURL, "translate"), strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.apiKey))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, Error{Code: res.StatusCode}
	}

	var response translationResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Translations, nil
}
