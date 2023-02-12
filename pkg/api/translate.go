package deepl

import (
	"encoding/json"
	"net/url"
)

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type translationResponse struct {
	Translations []Translation `json:"translations"`
}

type TranslateOption func(url.Values)

// The language to be translated.
// If this parameter is omitted, the API will attempt to detect the language of the text and translate it.
func SourceLang(lang string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("source_lang", lang)
	}
}

// Sets whether the translation engine should first split the input into sentences.
func SplitSentences(split string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("split_sentences", split)
	}
}

// Sets whether the translation engine should respect the original formatting, even if it would usually correct some aspects.
func PreserveFormatting(preserve string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("preserve_formatting", preserve)
	}
}

// Sets whether the translated text should lean towards formal or informal language.
func Formality(formality string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("formality", formality)
	}
}

// Specify the glossary to use for the translation.
func GlossaryId(glossary string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("glossary_id", glossary)
	}
}

// Sets which kind of tags should be handled.
func TagHandling(handling string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("tag_handling", handling)
	}
}

// Comma-separated list of XML tags which never split sentences.
func NonSplittingTags(tags string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("non_splitting_tags", tags)
	}
}

// Disable the automatic detection of the XML structure.
func OutlineDetection(detect string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("outline_detection", detect)
	}
}

// Comma-separated list of XML tags which always cause splts.
func SplittingTags(tags string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("splitting_tags", tags)
	}
}

// Comma-separated list of XML tags that indicate text not to be translated.
func IgnoreTags(tags string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("ignore_tags", tags)
	}
}

// The translate function.
func (c *Client) TranslateText(texts []string, targetLang string, options ...TranslateOption) ([]Translation, error) {
	vals := make(url.Values)

	for _, text := range texts {
		vals.Add("text", text)
	}

	vals.Set("target_lang", targetLang)

	for _, opt := range options {
		opt(vals)
	}

	res, err := c.do("POST", "translate", vals)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response translationResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Translations, nil
}
