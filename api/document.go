package deepl

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

func (t *Translator) TranslateDocumentUpload(path string, targetLang string, opts ...TranslateOptionFunc) (*DocumentInfo, error) {
	// Gather translate options
	options := TranslateOptions{}
	if err := options.Gather(opts...); err != nil {
		return nil, fmt.Errorf("error gathering options: %w", err)
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

		for name, value := range options {
			switch name {
			case "source_lang", "formality", "glossary_id":
				if err := mpw.WriteField(name, value); err != nil {
					errchan <- fmt.Errorf("error writing form field: %w", err)
					return
				}
			}
		}

		if part, err = mpw.CreateFormFile("file", filename); err != nil {
			errchan <- fmt.Errorf("error creating form file: %w", err)
			return
		}
		if _, err = io.Copy(part, f); err != nil {
			errchan <- fmt.Errorf("error writing file: %w", err)
			return
		}

		if err = mpw.Close(); err != nil {
			errchan <- fmt.Errorf("error closing multipart writer: %w", err)
			return
		}
	}()

	headers := make(http.Header)
	headers.Set("Content-Type", mpw.FormDataContentType())

	res, err := t.callAPI(http.MethodPost, "document", headers, r)
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
	endpoint := fmt.Sprintf("document/%s", id)

	vals := make(url.Values)
	vals.Set("document_key", key)

	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	body := strings.NewReader(vals.Encode())

	res, err := t.callAPI(http.MethodPost, endpoint, headers, body)
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
	endpoint := fmt.Sprintf("document/%s/result", id)

	vals := make(url.Values)
	vals.Set("document_key", key)

	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	body := strings.NewReader(vals.Encode())

	res, err := t.callAPI(http.MethodPost, endpoint, headers, body)
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
