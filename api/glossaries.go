package deepl

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	vals := make(url.Values)

	vals.Set("name", name)
	vals.Set("source_lang", sourceLang)
	vals.Set("target_lang", targetLang)
	vals.Set("entries_format", "tsv")

	var entriesTSV = make([]string, 0, len(entries))
	for _, entry := range entries {
		entriesTSV = append(entriesTSV, fmt.Sprintf("%s\t%s", entry.Source, entry.Target))
	}
	vals.Set("entries", strings.Join(entriesTSV, "\n"))

	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	body := strings.NewReader(vals.Encode())

	res, err := t.callAPI(http.MethodPost, "glossaries", headers, body)
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
	res, err := t.callAPI(http.MethodGet, "glossaries", nil, nil)
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
	endpoint := fmt.Sprintf("glossaries/%s", glossaryId)

	res, err := t.callAPI(http.MethodGet, endpoint, nil, nil)
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
	endpoint := fmt.Sprintf("glossaries/%s", glossaryId)

	res, err := t.callAPI(http.MethodDelete, endpoint, nil, nil)
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
	endpoint := fmt.Sprintf("glossaries/%s/entries", glossaryId)

	headers := make(http.Header)
	headers.Set("Accept", "text/tab-separated-values")

	res, err := t.callAPI(http.MethodGet, endpoint, headers, nil)
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
