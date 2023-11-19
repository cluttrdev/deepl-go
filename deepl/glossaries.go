package deepl

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GlossaryEntry struct {
	Source string
	Target string
}

type GlossaryInfo struct {
	GlossaryId   string `json:"glossary_id"`
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	SourceLang   string `json:"source_lang"`
	TargetLang   string `json:"target_lang"`
	CreationTime string `json:"creation_time"`
	EntryCount   int    `json:"entry_count"`
}

func (t *Translator) CreateGlossary(name string, sourceLang string, targetLang string, entries []GlossaryEntry) (*GlossaryInfo, error) {
	const (
		endpoint string = "v2/glossaries"
		method   string = http.MethodPost
	)

	data := struct {
		Name          string `json:"name"`
		SourceLang    string `json:"source_lang"`
		TargetLang    string `json:"target_lang"`
		Entries       string `json:"entries"`
		EntriesFormat string `json:"entries_format"`
	}{
		Name:          name,
		SourceLang:    sourceLang,
		TargetLang:    targetLang,
		EntriesFormat: "tsv",
	}
	data.Entries = encodeGlossaryEntries(entries...)

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding request data: %w", err)
	}

	res, err := t.callAPI(method, endpoint, headers, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return nil, httpError(res.StatusCode)
	}

	var glossary GlossaryInfo
	if err := json.NewDecoder(res.Body).Decode(&glossary); err != nil {
		return nil, err
	}

	return &glossary, nil
}

func (t *Translator) ListGlossaries() ([]GlossaryInfo, error) {
	const (
		endpoint string = "v2/glossaries"
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
		Glossaries []GlossaryInfo `json:"glossaries"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Glossaries, nil
}

func (t *Translator) GetGlossary(glossaryId string) (*GlossaryInfo, error) {
	var endpoint string = fmt.Sprintf("v2/glossaries/%s", glossaryId)
	const method string = http.MethodGet

	res, err := t.callAPI(method, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	var glossary GlossaryInfo
	if err := json.NewDecoder(res.Body).Decode(&glossary); err != nil {
		return nil, err
	}

	return &glossary, nil
}

func (t *Translator) DeleteGlossary(glossaryId string) error {
	var endpoint string = fmt.Sprintf("v2/glossaries/%s", glossaryId)
	const method string = http.MethodDelete

	res, err := t.callAPI(method, endpoint, nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return httpError(res.StatusCode)
	}

	return nil
}

func (t *Translator) GetGlossaryEntries(glossaryId string) ([]GlossaryEntry, error) {
	var endpoint string = fmt.Sprintf("v2/glossaries/%s/entries", glossaryId)
	const method string = http.MethodGet

	headers := make(http.Header)
	headers.Set("Accept", "text/tab-separated-values")

	res, err := t.callAPI(method, endpoint, headers, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	r := csv.NewReader(res.Body)
	r.Comma = '\t'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make([]GlossaryEntry, 0, len(records))
	for _, rec := range records {
		entries = append(entries, GlossaryEntry{Source: rec[0], Target: rec[1]})
	}

	return entries, nil
}

func encodeGlossaryEntries(entries ...GlossaryEntry) string {
	var encoded = make([]string, 0, len(entries))
	for _, entry := range entries {
		encoded = append(encoded, fmt.Sprintf("%s\t%s", entry.Source, entry.Target))
	}
	return strings.Join(encoded, "\n")
}
