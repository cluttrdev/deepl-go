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
	res, err := t.callAPI("GET", "usage", nil, nil)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}
	defer res.Body.Close()

	var usage Usage
	if err := json.NewDecoder(res.Body).Decode(&usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
