package deepl

import (
	"fmt"
)

type TranslateOptions struct {
	SourceLang         *string   `json:"source_lang,omitempty"`
	SplitSentences     *string   `json:"split_sentences,omitempty"`
	PreserveFormatting *bool     `json:"preserve_formatting,omitempty"`
	Formality          *string   `json:"formality,omitempty"`
	GlossaryID         *string   `json:"glossary_id,omitempty"`
	TagHandling        *string   `json:"tag_handling,omitempty"`
	OutlineDetection   *bool     `json:"outline_detection,omitempty"`
	NonSplittingTags   []*string `json:"non_splitting_tags,omitempty"`
	SplittingTags      []*string `json:"splitting_tags,omitempty"`
	IgnoreTags         []*string `json:"ignore_tags,omitempty"`
}

func (o *TranslateOptions) Gather(opts ...TranslateOption) error {
	for _, option := range opts {
		if err := option(o); err != nil {
			return err
		}
	}
	return nil
}

// TranslateOption can be used to customize the translation engine.
type TranslateOption func(*TranslateOptions) error

// WithSourceLang specifies the language of the text to be translated.
// If this parameter is omitted, the API will attempt to detect the language of the text and translate it
func WithSourceLang(value string) TranslateOption {
	return func(o *TranslateOptions) error {
		o.SourceLang = &value
		return nil
	}
}

// WithSplitSentences sets whether the translation engine should first split
// the input into sentences.
//
// For text translations where `tag_handling` is not set to `html`, the default
// value is `1`, meaning the engine splits on punctuation and on newlines.
// For text translations where `tag_handling=html`, the default value is
// `nonewlines`, meaning the engine splits on punctuation only, ignoring
// newlines.
//
// The use of `nonewlines` as the default value for text translations where
// `tag_handling=html` is new behavior that was implemented in November 2022,
// when HTML handling was moved out of beta.
//
// Possible values are:
//
// `0` - no splitting at all, whole input is treated as one sentence
// `1` - splits on punctuation and on newlines (default when `tag_handling` is not set to `html`)
// `nonewlines` - splits on punctuation only, ignoring newlines (default when `tag_handling=html`)
//
// For applications that send one sentence per text parameter, we recommend
// setting `split_sentences` to `0`, in order to prevent the engine from splitting
// the sentence unintentionally.
//
// Please note that newlines will split sentences when `split_sentences=1`. We
// recommend cleaning files so they don't contain breaking sentences or setting
// the parameter `split_sentences` to `nonewlines`.
func WithSplitSentences(value string) TranslateOption {
	return func(o *TranslateOptions) error {
		switch value {
		case "0", "1", "nonewlines":
			o.SplitSentences = &value
			return nil
		}
		return translateOptionInvalidValueError("split_sentences", value)
	}
}

// WithPreserveFormatting sets whether the translation engine should respect
// the original formatting, even if it would usually correct some aspects.
//
// The formatting aspects affected by this setting include:
//   - Punctuation at the beginning and end of the sentence
//   - Upper/lower case at the beginning of the sentence
func WithPreserveFormatting(value bool) TranslateOption {
	return func(o *TranslateOptions) error {
		o.PreserveFormatting = &value
		return nil
	}
}

// WithFormality sets whether the translated text should lean towards formal or
// informal language.
//
// This feature currently only works for target languages DE (German), FR
// (French), IT (Italian), ES (Spanish), NL (Dutch), PL (Polish), PT-BR and
// PT-PT (Portuguese), JA (Japanese), and RU (Russian).
// Learn more about the plain/polite feature for Japanese [here][formality-japanese].
//
// Setting this parameter with a target language that does not support
// formality will fail, unless one of the `prefer_...` options are used.
//
// Possible options are:
//   - `default` (default)
//   - `more` - for a more formal language
//   - `less` - for a more informal language
//   - `prefer_more` - for a more formal language if available, otherwise fallback to default formality
//   - `prefer_less` - for a more informal language if available, otherwise fallback to default formality
//
// [formality-japanese]: https://support.deepl.com/hc/en-us/articles/6306700061852-About-the-plain-polite-feature-in-Japanese
func WithFormality(value string) TranslateOption {
	return func(o *TranslateOptions) error {
		switch value {
		case "default", "more", "less", "prefer_more", "prefer_less":
			o.Formality = &value
			return nil
		}
		return translateOptionInvalidValueError("formality", value)
	}
}

// WithGlossaryID specifies the glossary to use for the translation.
//
// *Important*: This requires the `source_lang` parameter to be set and the
// language pair of the glossary has to match the language pair of the request.
func WithGlossaryID(value string) TranslateOption {
	return func(o *TranslateOptions) error {
		o.GlossaryID = &value
		return nil
	}
}

// WithTagHandling sets which kind of tags should be handled.
//
// Options currently available:
//   - `xml`: Enable XML tag handling; see [XML Handling][xml-handling].
//   - `html`: Enable HTML tag handling; see [HTML Handling][html-handling].
//
// [xml-handling]: https://www.deepl.com/docs-api/xml
// [html-handling]: https://www.deepl.com/docs-api/html
func WithTagHandling(value string) TranslateOption {
	return func(o *TranslateOptions) error {
		switch value {
		case "html", "xml":
			o.TagHandling = &value
			return nil
		}
		return translateOptionInvalidValueError("tag_handling", value)
	}
}

// WithOutlineDetection can be used to disable the automatic detection of the
// XML structure.
//
// The automatic detection of the XML structure won't yield best results in all
// XML files. You can disable this automatic mechanism altogether by setting
// the `outline_detection` parameter to `0` and selecting the tags that should
// be considered structure tags. This will split sentences using the
// `splitting_tags` parameter.
//
// In the example below, we achieve the same results as the automatic engine by
// disabling automatic detection with `outline_detection=0` and setting the
// parameters manually to `tag_handling=xml`, `split_sentences=nonewlines`, and
// `splitting_tags=par,title`.
//
// Example request:
// ```
// <document>
//
//	<meta>
//	  <title>A document's title</title>
//	</meta>
//	<content>
//	  <par>This is the first sentence. Followed by a second one.</par>
//	  <par>This is the third sentence.</par>
//	</content>
//
// </document>
// ```
//
// Example response:
// ```
// <document>
//
//	<meta>
//	  <title>Der Titel eines Dokuments</title>
//	</meta>
//	<content>
//	  <par>Das ist der erste Satz. Gefolgt von einem zweiten.</par>
//	  <par>Dies ist der dritte Satz.</par>
//	</content>
//
// </document>
// ```
//
// While this approach is slightly more complicated, it allows for greater
// control over the structure of the translation output.
func WithOutlineDetection(value bool) TranslateOption {
	return func(o *TranslateOptions) error {
		o.OutlineDetection = &value
		return nil
	}
}

// WithNonSplittingTags specifies a list of XML tags which
// never split sentences.
//
// For some XML files, finding tags with textual content and splitting
// sentences using those tags won't yield the best results. The following
// example shows the engine splitting sentences on `par` tags and proceeding to
// translate the parts separately, resulting in an incorrect translation:
//
// Example request:
// ```
// <par>The firm said it had been </par><par> conducting an internal investigation.</par>
// ```
//
// Example response:
// ```
// <par>Die Firma sagte, es sei eine gute Idee gewesen.</par><par> Durchführung einer internen Untersuchung.</par>
// ```
//
// As this can lead to bad translations, this type of structure should either
// be avoided, or the `non_splitting_tags` parameter should be set.
//
// The following example shows the same call, with the parameter set to `par`:
//
// Example request:
// ```
// <par>The firm said it had been </par><par> conducting an internal investigation.</par>
// ```
//
// Example response:
// ```
// <par>Die Firma sagte, dass sie</par><par> eine interne Untersuchung durchgeführt</par><par> habe</par><par>.</par>
// ```
//
// This time, the sentence is translated as a whole. The XML tags are now
// considered markup and copied into the translated sentence. As the
// translation of the words "had been" has moved to another position in the
// German sentence, the two `par` tags are duplicated (which is expected here).
func WithNonSplittingTags(value []string) TranslateOption {
	return func(o *TranslateOptions) error {
		for _, v := range value {
			o.NonSplittingTags = append(o.NonSplittingTags, &v)
		}
		return nil
	}
}

// WithSplittingTags specifies a list of XML tags which always
// cause splits.
//
// See the example in the `outline_detection` parameter's description.
func WithSplittingTags(value []string) TranslateOption {
	return func(o *TranslateOptions) error {
		for _, v := range value {
			o.SplittingTags = append(o.SplittingTags, &v)
		}
		return nil
	}
}

// WithIgnoreTags specifies a list of XML tags that indicate
// text not to be translated.
//
// Use this parameter to ensure that elements in the original text are not
// altered in the translation (e.g., trademarks, product names) and insert tags
// into your original text. In the following example, the `ignore_tags` parameter
// is set to `keep`:
//
// Example request:
// ```
// Please open the page <keep>Settings</keep> to configure your system.
// ```
//
// Example response:
// ```
// Bitte öffnen Sie die Seite <keep>Settings</keep> um Ihr System zu konfigurieren.
// ```
func WithIgnoreTags(value []string) TranslateOption {
	return func(o *TranslateOptions) error {
		for _, v := range value {
			o.IgnoreTags = append(o.IgnoreTags, &v)
		}
		return nil
	}
}

func translateOptionInvalidValueError(name string, value string) error {
	return fmt.Errorf("Invalid value for option `%s`: %s", name, value)
}
