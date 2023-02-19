package deepl

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

// TranslateOption is a functional option for configuring text translation parameters
type TranslateOption func(url.Values)

// The language to be translated.
// If this parameter is omitted, the API will attempt to detect the language of the text and translate it
func SourceLang(lang string) (TranslateOption, error) {
	return func(vals url.Values) {
		vals.Set("source_lang", lang)
	}, nil
}

// SplitSentences sets whether the translation engine should first split the input into sentences
func SplitSentences(split string) (TranslateOption, error) {
	switch split {
	case "0", "1", "nonewlines":
		return func(vals url.Values) {
			vals.Set("split_sentences", split)
		}, nil
	}
	return nil, errors.Errorf("Invalid SplitSentence value: %s", split)
}

// PreserveFormatting sets whether the translation engine should respect the original formatting, even if it would usually correct some aspects
func PreserveFormatting(preserve string) (TranslateOption, error) {
	switch preserve {
	case "0", "1":
		return func(vals url.Values) {
			vals.Set("preserve_formatting", preserve)
		}, nil
	}
	return errors.Errorf("Invalid PreserveFormatting value: %s", preserve)
}

// Formality sets whether the translated text should lean towards formal or informal language
func Formality(formality string) (TranslateOption, error) {
	switch formality {
	case "default", "more", "less", "prefer_more", "prefer_less":
		return func(vals url.Values) {
			vals.Set("formality", formality)
		}, nil
	}
	return nil, errors.Errorf("Invalid Formality value: %s", formality)
}

// GlossaryId specifies the glossary to use for the translation
func GlossaryId(glossary string) (TranslateOption, error) {
	return func(vals url.Values) {
		vals.Set("glossary_id", glossary)
	}, nil
}

// TagHandling sets which kind of tags should be handled
func TagHandling(handling string) (TranslateOption, error) {
	switch handling {
	case "html", "xml":
		return func(vals url.Values) {
			vals.Set("tag_handling", handling)
		}, nil
	}
	return nil, errors.Errorf("Invalid TagHandling value: %s", handling)
}

// NonSplittingTags specifies a comma-separated list of XML tags which never split sentences
func NonSplittingTags(tags string) (TranslateOption, error) {
	return func(vals url.Values) {
		vals.Set("non_splitting_tags", tags)
	}, nil
}

// OutlineDetection can be used to disable the automatic detection of the XML structure
func OutlineDetection(detect string) (TranslateOption, error) {
	switch detect {
	case "0":
		return func(vals url.Values) {
			vals.Set("outline_detection", detect)
		}, nil
	}
	return nil, errors.Errorf("Invalid OutlineDetection value: %s", detect)
}

// SplittingTags specifies a comma-separated list of XML tags which always cause splts
func SplittingTags(tags string) (TranslateOption, error) {
	return func(vals url.Values) {
		vals.Set("splitting_tags", tags)
	}, nil
}

// IgnoeTags specifies a comma-separated list of XML tags that indicate text not to be translated
func IgnoreTags(tags string) (TranslateOption, error) {
	return func(vals url.Values) {
		vals.Set("ignore_tags", tags)
	}, nil
}

// TranslateText translates the given text(s) into the specified target language
func (t *Translator) TranslateText(texts []string, targetLang string, options ...TranslateOption) ([]Translation, error) {
	vals := make(url.Values)

	for _, text := range texts {
		vals.Add("text", text)
	}

	vals.Set("target_lang", targetLang)

	// Apply translation parameter options
	for _, opt := range options {
		opt(vals)
	}

	res, err := t.callAPI("POST", "translate", vals, nil)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, HTTPError{StatusCode: res.StatusCode}
	}
	defer res.Body.Close()

	var response struct {
		Translations []Translation `json:"translations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Translations, nil
}
