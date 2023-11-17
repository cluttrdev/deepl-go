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

func NewGlossariesCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "glossaries",
		ShortHelp:  "Manage glossaries",
		ShortUsage: "glossaries [command] [option]... [args]...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
		Subcommands: []*command.Command{
			NewGlossariesLanguagePairsCmd(stdout, stderr),
			NewGlossariesCreateCmd(stdout, stderr),
			NewGlossariesListCmd(stdout, stderr),
			NewGlossariesInfoCmd(stdout, stderr),
			NewGlossariesEntriesCmd(stdout, stderr),
			NewGlossariesDeleteCmd(stdout, stderr),
		},
	}
}

type GlossariesCmdConfig struct {
	RootCmdConfig
}

func (c *GlossariesCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *GlossariesCmdConfig) Exec(context.Context, []string) error {
	return flag.ErrHelp
}

/*
 *  PAIRS
 */

func NewGlossariesLanguagePairsCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesLanguagePairsCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries language-pairs", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "language-pairs",
		ShortHelp:  "List language pairs supported by glossaries",
		ShortUsage: "glossaries language-pairs",
		LongHelp:   "Retrieve the list of language pairs supported by the glossary feature.",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type GlossariesLanguagePairsCmdConfig struct {
	RootCmdConfig
}

func (c *GlossariesLanguagePairsCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *GlossariesLanguagePairsCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		fmt.Fprintln(c.stderr, "Error: glossary language-pairs: too many arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	lps, err := t.GetGlossaryLanguagePairs()
	if err != nil {
		return err
	}

	m, err := json.Marshal(lps)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.stdout, string(m))

	return nil
}

/*
 *  CREATE
 */

func NewGlossariesCreateCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesCreateCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries create", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "create",
		ShortHelp:  "Create a glossary",
		ShortUsage: "glossaries create [option]... ENTRY...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type GlossariesCreateCmdConfig struct {
	RootCmdConfig

	name       string
	sourceLang string
	targetLang string
}

func (c *GlossariesCreateCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)

	fs.StringVar(&c.name, "name", "", "the name to be associated with the glossary (required)")
	fs.StringVar(&c.sourceLang, "source-lang", "", "the language in which the source texts in the glossary are specified (required)")
	fs.StringVar(&c.sourceLang, "from", "", "alias option for `--source-lang`")
	fs.StringVar(&c.targetLang, "target-lang", "", "the language in which the target texts in the glossary are specified (required)")
	fs.StringVar(&c.targetLang, "to", "", "alias option for `--target-lang`")
}

func (c *GlossariesCreateCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: glossary create: not enough arguments")
		return flag.ErrHelp
	}

	if c.name == "" || c.sourceLang == "" || c.targetLang == "" {
		fmt.Fprintln(c.stderr, "Error: glossary create: `--name`,`--source-lang` and `--target-lang` are required")
		return flag.ErrHelp
	}

	var entries = make([]deepl.GlossaryEntry, 0, len(args))
	for _, arg := range args {
		pair := strings.Split(arg, "=")
		if len(pair) != 2 {
			fmt.Fprintf(c.stderr, "Error: glossary create: invalid argument: %s", arg)
			return flag.ErrHelp
		}
		entries = append(entries, deepl.GlossaryEntry{Source: pair[0], Target: pair[1]})
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	g, err := t.CreateGlossary(c.name, c.sourceLang, c.targetLang, entries)
	if err != nil {
		return err
	}

	m, err := json.Marshal(g)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.stdout, string(m))

	return nil
}

/*
 *  LIST
 */

func NewGlossariesListCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesListCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries list", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "list",
		ShortHelp:  "List all glossaries",
		ShortUsage: "glossaries list [option]...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type GlossariesListCmdConfig struct {
	RootCmdConfig
}

func (c *GlossariesListCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *GlossariesListCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		fmt.Fprintln(c.stderr, "Error: glossaries list: too many arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	gs, err := t.ListGlossaries()
	if err != nil {
		return err
	}

	m, err := json.Marshal(gs)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.stdout, string(m))

	return nil
}

/*
 *  INFO
 */

func NewGlossariesInfoCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesInfoCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries info", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "info",
		ShortHelp:  "Retrieve glossary details",
		ShortUsage: "glossaries info [option]... ID...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type GlossariesInfoCmdConfig struct {
	RootCmdConfig
}

func (c *GlossariesInfoCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *GlossariesInfoCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: glossaries info: not enough arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	var gs []*deepl.GlossaryInfo

	defer func() {
		m, err := json.Marshal(gs)
		if err != nil {
			fmt.Fprintln(c.stderr, err)
		}
		fmt.Fprintln(c.stdout, string(m))
	}()

	for _, gid := range args {
		g, err := t.GetGlossary(gid)
		if err != nil {
			return err
		}
		gs = append(gs, g)
	}

	return nil
}

/*
 *  ENTRIES
 */

func NewGlossariesEntriesCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesEntriesCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries entries", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "entries",
		ShortHelp:  "Retrieve glossary entries",
		ShortUsage: "glossaries entries [option]... ID",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type GlossariesEntriesCmdConfig struct {
	RootCmdConfig

	entriesFormat string
}

func (c *GlossariesEntriesCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)

	fs.StringVar(&c.entriesFormat, "format", "tsv", "the requested format of the returned glossary entries")
}

func (c *GlossariesEntriesCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: glossaries entries: not enough arguments")
		return flag.ErrHelp
	} else if len(args) > 0 {
		fmt.Fprintln(c.stderr, "Error: glossaries entries: too many arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	gid := args[0]

	ges, err := t.GetGlossaryEntries(gid)
	if err != nil {
		return err
	}

	var sep rune
	switch c.entriesFormat {
	case "tsv":
		sep = '\t'
	case "csv":
		sep = ','
	default:
		return fmt.Errorf("glossaries entries: invalid value for option `format`: %s", c.entriesFormat)
	}

	for _, ge := range ges {
		fmt.Fprintf(c.stdout, "%s%c%s\n", ge.Source, sep, ge.Target)
	}

	return nil
}

/*
 *  DELETE
 */

func NewGlossariesDeleteCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := GlossariesDeleteCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("glossaries delete", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "delete",
		ShortHelp:  "Delete glossaries",
		ShortUsage: "glossaries delete [option]... ID...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

type GlossariesDeleteCmdConfig struct {
	RootCmdConfig
}

func (c *GlossariesDeleteCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *GlossariesDeleteCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "Error: glossaries delete: not enough arguments")
		return flag.ErrHelp
	}

	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	var gids []string

	defer func() {
		fmt.Fprintln(c.stdout, strings.Join(gids, "\n"))
	}()

	for _, gid := range args {
		err := t.DeleteGlossary(gid)
		if err != nil {
			return err
		}
		gids = append(gids, gid)
	}

	return nil
}
