package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type DocumentInfo struct {
	DocumentId  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

type DocumentStatus struct {
	DocumentId string `json:"document_id"`
	Status     string `json:"status"`

	// Status dependent additional fields
	SecondsRemaining string `json:"seconds_remaining:"`
	BilledCharacters int    `json:"billed_characters"`
	Message          string `json:"message"`
}

func (t *Translator) TranslateDocumentUpload(path string, targetLang string, opts ...TranslateOption) (*DocumentInfo, error) {
	const (
		endpoint string = "v2/document"
		method   string = http.MethodPost
	)

	// Gather translate options
	options := TranslateOptions{}
	if err := options.Gather(opts...); err != nil {
		return nil, fmt.Errorf("error gathering options: %w", err)
	}

	fields := map[string]*string{
		"source_lang": options.SourceLang,
		"formality":   options.Formality,
		"glossary_id": options.GlossaryID,
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	r, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	errchan := make(chan error)
	go func() {
		defer w.Close()
		defer f.Close()
		defer close(errchan)

		var part io.Writer

		filename := filepath.Base(path)

		if err := mpw.WriteField("filename", filename); err != nil {
			errchan <- fmt.Errorf("error writing form field: %w", err)
			return
		}
		if err := mpw.WriteField("target_lang", targetLang); err != nil {
			errchan <- fmt.Errorf("error writing form field: %w", err)
			return
		}

		for name, value := range fields {
			if value == nil {
				continue
			}

			if err := mpw.WriteField(name, *value); err != nil {
				errchan <- fmt.Errorf("error writing form field: %w", err)
				return
			}
		}

		if part, err = mpw.CreateFormFile("file", filename); err != nil {
			errchan <- fmt.Errorf("error creating form file: %w", err)
			return
		}
		if _, err = io.Copy(part, f); err != nil {
			errchan <- fmt.Errorf("error writing form file: %w", err)
			return
		}

		if err = mpw.Close(); err != nil {
			errchan <- fmt.Errorf("error closing multipart writer: %w", err)
			return
		}
	}()

	headers := make(http.Header)
	headers.Set("Content-Type", mpw.FormDataContentType())

	res, err := t.callAPI(method, endpoint, headers, r)
	merr := <-errchan
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if merr != nil {
		return nil, merr
	}
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	var document DocumentInfo
	if err := json.NewDecoder(res.Body).Decode(&document); err != nil {
		return nil, err
	}

	return &document, err
}

func (t *Translator) TranslateDocumentStatus(id string, key string) (*DocumentStatus, error) {
	var endpoint string = fmt.Sprintf("v2/document/%s", id)
	const method string = http.MethodPost

	data := struct {
		DocumentKey string `json:"document_key"`
	}{
		DocumentKey: key,
	}

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
	if res.StatusCode != http.StatusOK {
		return nil, httpError(res.StatusCode)
	}

	var status DocumentStatus
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (t *Translator) TranslateDocumentDownload(id string, key string) (*io.PipeReader, error) {
	var endpoint string = fmt.Sprintf("v2/document/%s/result", id)
	const method string = http.MethodPost

	data := struct {
		DocumentKey string `json:"document_key"`
	}{
		DocumentKey: key,
	}

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
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		return nil, httpError(res.StatusCode)
	}

	r, w := io.Pipe()
	go func() {
		defer w.Close()
		defer res.Body.Close()

		if _, err := io.Copy(w, res.Body); err != nil {
			_ = w.CloseWithError(err) // always returns nil
		}
	}()

	return r, nil
}
