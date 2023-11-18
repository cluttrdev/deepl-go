package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cluttrdev/deepl-go/deepl"

	"github.com/cluttrdev/deepl-go/internal/command"
)

func NewDocumentCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	return &command.Command{
		Name:       "document",
		ShortHelp:  "Translate documents",
		ShortUsage: "deepl document [command] [option]...",
		LongHelp:   "",
		Subcommands: []*command.Command{
			NewDocumentUploadCmd(stdout, stderr),
			NewDocumentStatusCmd(stdout, stderr),
			NewDocumentDownloadCmd(stdout, stderr),
		},
	}
}

/*
 *  UPLOAD
 */

func NewDocumentUploadCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := DocumentUploadCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},

		flags: flag.NewFlagSet("document upload", flag.ContinueOnError),
	}

	cfg.RegisterFlags(cfg.flags)

	return &command.Command{
		Name:       "upload",
		ShortHelp:  "Upload documents for translation",
		ShortUsage: "deepl document upload [option]... --target-lang=LANG FILE...",
		LongHelp:   "",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

type DocumentUploadCmdConfig struct {
	RootCmdConfig

	flags *flag.FlagSet

	targetLang string
	sourceLang string
	formality  string
	glossaryID string
}

func (c *DocumentUploadCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)

	fs.StringVar(&c.targetLang, "target-lang", "", "the language into which the text should be translated (required)")
	fs.StringVar(&c.targetLang, "to", "", "alias option for `--target-lang`")
	fs.StringVar(&c.sourceLang, "source-lang", "", "the language to be translated")
	fs.StringVar(&c.sourceLang, "from", "", "alias option for `--source-lang`")
	fs.StringVar(&c.formality, "formality", "default", "whether the engine should lean towards formal or informal language")
	fs.StringVar(&c.glossaryID, "glossary_id", "", "the glossary to use for the translation")
}

func (c *DocumentUploadCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: document upload: not enough arguments")
		return flag.ErrHelp
	}

	if c.targetLang == "" {
		fmt.Fprintln(c.stderr, "Error: document upload: `--target-lang` is required")
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
		case "formality":
			opts = append(opts, deepl.WithFormality(c.formality))
		case "glossary-id":
			opts = append(opts, deepl.WithGlossaryID(c.glossaryID))
		}
	})

	type document struct {
		Path string `json:"document_path"`

		deepl.DocumentInfo
	}

	var docs []document

	defer func() {
		marshalled, err := json.Marshal(docs)
		if err != nil {
			fmt.Fprintln(c.stderr, err)
		}
		fmt.Fprintln(c.stdout, string(marshalled))
	}()

	for _, path := range args {
		di, err := t.TranslateDocumentUpload(path, c.targetLang, opts...)
		if err != nil {
			return err
		}
		docs = append(docs, document{
			Path:         path,
			DocumentInfo: *di,
		})
	}

	return nil
}

/*
 *  STATUS
 */

func NewDocumentStatusCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := DocumentStatusCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("document status", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "status",
		ShortHelp:  "Retrieve the current status of a document translation process.",
		ShortUsage: "deepl document status [option]... ID:KEY...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type DocumentStatusCmdConfig struct {
	RootCmdConfig
}

func (c *DocumentStatusCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *DocumentStatusCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: document status: not enough arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	var dis []deepl.DocumentInfo
	for _, arg := range args {
		info := strings.Split(arg, ":")
		if len(info) != 2 {
			return fmt.Errorf("invalid argument: %s", arg)
		}
		dis = append(dis, deepl.DocumentInfo{
			DocumentId:  info[0],
			DocumentKey: info[1],
		})
	}

	var dss []*deepl.DocumentStatus

	defer func() {
		marshalled, err := json.Marshal(dss)
		if err != nil {
			fmt.Fprintln(c.stderr, err)
		}
		fmt.Fprintln(c.stdout, string(marshalled))
	}()

	for _, di := range dis {
		ds, err := t.TranslateDocumentStatus(di.DocumentId, di.DocumentKey)
		if err != nil {
			return err
		}
		dss = append(dss, ds)
	}

	return nil
}

/*
 *  DOWNLOAD
 */

func NewDocumentDownloadCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := DocumentDownloadCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("document download", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "download",
		ShortHelp:  "Download a document after translation",
		ShortUsage: "deepl document download [option]... ID:KEY",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type DocumentDownloadCmdConfig struct {
	RootCmdConfig
}

func (c *DocumentDownloadCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *DocumentDownloadCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: document download: not enough arguments")
		return flag.ErrHelp
	} else if len(args) > 1 {
		fmt.Fprintln(c.stderr, "Error: document download: too many arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	info := strings.Split(args[0], ":")
	if len(info) != 2 {
		return fmt.Errorf("invalid argument: %s", args[0])
	}
	di := deepl.DocumentInfo{
		DocumentId:  info[0],
		DocumentKey: info[1],
	}

	pr, err := t.TranslateDocumentDownload(di.DocumentId, di.DocumentKey)
	if err != nil {
		return nil
	}
	defer pr.Close()

	_, err = io.Copy(c.stdout, pr)
	if err != nil {
		return err
	}

	return nil
}
