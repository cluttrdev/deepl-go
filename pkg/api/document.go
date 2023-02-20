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

func (t *Translator) TranslateDocumentUpload(filePath string, targetLang string, options ...TranslateOption) (*DocumentInfo, error) {
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

	r, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	go func() {
		var part io.Writer

		defer w.Close()
		defer f.Close()

		mpw.WriteField("filename", filepath.Base(filePath))
		mpw.WriteField("target_lang", targetLang)

		for _, opt := range options {
			switch opt.Key {
			case "source_lang", "formality", "glossary_id":
				mpw.WriteField(opt.Key, opt.Value)
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

	url := fmt.Sprintf("%s/%s", t.serverURL, "document")

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", t.authKey))
	req.Header.Add("Content-Type", mpw.FormDataContentType())

	res, err := t.httpClient.httpClient.Do(req)
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

	res, err := t.callAPI("POST", endpoint, vals, nil)
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

	res, err := t.callAPI("POST", endpoint, vals, nil)
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
