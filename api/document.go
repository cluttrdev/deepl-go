package deepl

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (t *Translator) TranslateDocumentUpload(filePath string, targetLang string, opts ...TranslateOptionFunc) (*DocumentInfo, error) {
	var (
		err error
		f   *os.File
		fi  os.FileInfo
	)

	if f, err = os.Open(filePath); err != nil {
		log.Fatal(err)
	}
	if fi, err = f.Stat(); err != nil {
		log.Fatal(err)
	}

	options := TranslateOptions{}
	if err := options.Gather(opts...); err != nil {
		return nil, fmt.Errorf("error gathering options: %w", err)
	}

	r, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	go func() {
		var part io.Writer

		defer w.Close()
		defer f.Close()

		mpw.WriteField("filename", filepath.Base(filePath))
		mpw.WriteField("target_lang", targetLang)

		for name, value := range options {
			switch name {
			case "source_lang", "formality", "glossary_id":
				mpw.WriteField(name, value)
			}
		}

		if part, err = mpw.CreateFormFile("file", filepath.Base(fi.Name())); err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(part, f); err != nil {
			log.Fatal(err)
		}

		if err = mpw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	headers := make(http.Header)
	headers.Set("Content-Type", mpw.FormDataContentType())

	res, err := t.callAPI(http.MethodPost, "document", headers, r)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}
	defer res.Body.Close()

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

	body := strings.NewReader(vals.Encode())

	res, err := t.callAPI(http.MethodPost, endpoint, nil, body)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}
	defer res.Body.Close()

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

	body := strings.NewReader(vals.Encode())

	res, err := t.callAPI(http.MethodPost, endpoint, nil, body)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}

	r, w := io.Pipe()
	go func() {
		defer w.Close()
		defer res.Body.Close()

		if _, err := io.Copy(w, res.Body); err != nil {
			log.Fatal(err)
		}
	}()

	return r, nil
}
