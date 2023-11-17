package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	deepl "github.com/cluttrdev/deepl-go/api"

	"github.com/cluttrdev/deepl-go/internal/command"
)

type LanguagesCmdConfig struct {
	RootCmdConfig

	langType       string
	sourceLangFlag bool
	targetLangFlag bool
}

func NewLanguagesCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := LanguagesCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("languages", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "languages",
		ShortHelp:  "Retrieve supported languages.",
		ShortUsage: "deepl languages [option]...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *LanguagesCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)

	fs.StringVar(&c.langType, "type", "", "the type of supported languages to retreive (`source` or `target`)")
	fs.BoolVar(&c.sourceLangFlag, "source", false, "shorthand option for `--type=source`")
	fs.BoolVar(&c.targetLangFlag, "target", false, "shorthand option for `--type=target`")
}

func (c *LanguagesCmdConfig) Exec(ctx context.Context, args []string) error {
	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	if btoi(len(c.langType) > 0)+btoi(c.sourceLangFlag)+btoi(c.targetLangFlag) > 1 {
		return errors.New("languages: `--type`,`--source` and `--target` options are mutually exclusive")
	}

	if c.sourceLangFlag {
		c.langType = "source"
	} else if c.targetLangFlag {
		c.langType = "target"
	}

	ls, err := t.GetLanguages(c.langType)
	if err != nil {
		return err
	}

	if int(c.verbosity) > 0 {
		fmt.Fprintf(c.stdout, "%s languages available:\n", cases.Title(language.English).String(c.langType))
	}
	writeLanguages(c.stdout, ls)
	return nil
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func writeLanguages(out io.Writer, langs []deepl.Language) {
	for _, lang := range langs {
		if lang.SupportsFormality {
			fmt.Fprintf(out, "%s: %s (supports formality)\n", lang.Code, lang.Name)
		} else {
			fmt.Fprintf(out, "%s: %s\n", lang.Code, lang.Name)
		}
	}
}
