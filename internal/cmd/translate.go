package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	deepl "github.com/cluttrdev/deepl-go/api"

	"github.com/cluttrdev/deepl-go/internal/command"
)

func NewTranslateTextCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := TranslateTextCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},

		flags: flag.NewFlagSet("translate", flag.ContinueOnError),
	}

	cfg.RegisterFlags(cfg.flags)

	return &command.Command{
		Name:       "translate",
		ShortHelp:  "Translate text(s) into a target language.",
		ShortUsage: "deepl translate [option]... --target-lang=LANG TEXT...",
		LongHelp:   "",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

type TranslateTextCmdConfig struct {
	RootCmdConfig

	flags *flag.FlagSet

	sourceLang         string
	targetLang         string
	splitSentences     string
	preserveFormatting bool
	formality          string
	glossaryID         string
	tagHandling        string
	outlineDetection   bool
	nonSplittingTags   string
	splittingTags      string
	ignoreTags         string

	formatJSON bool
}

func (c *TranslateTextCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)

	fs.StringVar(&c.targetLang, "target-lang", "", "the language into which the text should be translated (required)")
	fs.StringVar(&c.targetLang, "to", "", "alias option for `--target-lang`")
	fs.StringVar(&c.sourceLang, "source-lang", "", "the language to be translated")
	fs.StringVar(&c.sourceLang, "from", "", "alias option for `--source-lang`")
	fs.StringVar(&c.splitSentences, "split-sentences", "0", "whether to split input into sentences")
	fs.BoolVar(&c.preserveFormatting, "preserve-formatting", false, "whether the engine should respect original formatting")
	fs.StringVar(&c.formality, "formality", "default", "whether the engine should lean towards formal or informal language")
	fs.StringVar(&c.glossaryID, "glossary_id", "", "the glossary to use for the translation")
	fs.StringVar(&c.tagHandling, "tag-handling", "", "the kind of tags to handle")
	fs.BoolVar(&c.outlineDetection, "outline-detection", true, "whether to automatically detect XML structure")
	fs.StringVar(&c.nonSplittingTags, "non-splitting-tags", "", "a comma-separated list of XML tags which never split sentences")
	fs.StringVar(&c.splittingTags, "splitting-tags", "", "a comma-separated list of XML tags which always split sentences")
	fs.StringVar(&c.ignoreTags, "ignore-tags", "", "a comma-separated list of XML tags which indicate text not to be translated")

	fs.BoolVar(&c.formatJSON, "json", false, "print translation result in JSON")
}

func (c *TranslateTextCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Fprintln(c.stderr, "Error: translate: not enough arguments")
		return flag.ErrHelp
	}

	if c.targetLang == "" {
		fmt.Fprintln(c.stderr, "Error: translate: `--target-lang` is required")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	opts := []deepl.TranslateOption{}
	c.flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "source-lang":
			opts = append(opts, deepl.WithSourceLang(c.sourceLang))
		case "split-sentences":
			opts = append(opts, deepl.WithSplitSentences(c.splitSentences))
		case "preserve-formatting":
			opts = append(opts, deepl.WithPreserveFormatting(c.preserveFormatting))
		case "formality":
			opts = append(opts, deepl.WithFormality(c.formality))
		case "glossary-id":
			opts = append(opts, deepl.WithGlossaryID(c.glossaryID))
		case "tag-handling":
			opts = append(opts, deepl.WithTagHandling(c.tagHandling))
		case "outline-detection":
			opts = append(opts, deepl.WithOutlineDetection(c.outlineDetection))
		case "non-splitting-tags":
			opts = append(opts, deepl.WithNonSplittingTags(strings.Split(c.nonSplittingTags, ",")))
		case "splitting-tags":
			opts = append(opts, deepl.WithSplittingTags(strings.Split(c.splittingTags, ",")))
		case "ignore-tags":
			opts = append(opts, deepl.WithIgnoreTags(strings.Split(c.ignoreTags, ",")))
		}
	})

	ts, err := t.TranslateText(args, c.targetLang, opts...)
	if err != nil {
		return err
	}

	if c.formatJSON {
		m, err := json.Marshal(ts)
		if err != nil {
			return err
		}
		fmt.Fprintln(c.stdout, string(m))
	} else {
		verbosity := int(c.verbosity)
		for _, tt := range ts {
			if verbosity > 0 {
				fmt.Fprintf(c.stdout, "# Detected source language: %s\n", tt.DetectedSourceLanguage)
			}
			fmt.Fprintln(c.stdout, tt.Text)
		}
	}

	return nil
}
