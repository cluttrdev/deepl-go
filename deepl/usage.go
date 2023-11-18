package deepl

import (
	"encoding/json"
	"net/http"
)

type Usage struct {
	CharacterCount    int `json:"character_count"`
	CharacterLimit    int `json:"character_limit"`
	DocumentCount     int `json:"document_count"`
	DocumentLimit     int `json:"document_limit"`
	TeamDocumentCount int `json:"team_document_count"`
	TeamDocumentLimit int `json:"team_document_limit"`
}

func (t *Translator) GetUsage() (*Usage, error) {
	const (
		endpoint string = "v2/usage"
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

	var usage Usage
	if err := json.NewDecoder(res.Body).Decode(&usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
