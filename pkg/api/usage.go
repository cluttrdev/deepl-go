package deepl

import (
	"encoding/json"
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
	res, err := t.httpClient.do("GET", "usage", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var usage Usage
	if err := json.NewDecoder(res.Body).Decode(&usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
