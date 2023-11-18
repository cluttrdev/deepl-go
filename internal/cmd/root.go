package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cluttrdev/deepl-go/deepl"

	"github.com/cluttrdev/deepl-go/internal/command"
)

type Verbosity int

func (v *Verbosity) String() string {
	return fmt.Sprintf("%v", int(*v))
}

func (v *Verbosity) Set(s string) error {
	if s != "" {
		return fmt.Errorf("verbosity flag does not take a value: %v", s)
	}
	*v++
	return nil
}

type RootCmdConfig struct {
	stdout io.Writer
	stderr io.Writer

	authKey   string
	serverURL string

	verbosity int
}

func NewRootCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := RootCmdConfig{
		stdout: stdout,
		stderr: stderr,
	}

	fs := flag.NewFlagSet("deepl", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "deepl",
		ShortHelp:  "deepl - DeepL language translation cli",
		ShortUsage: "deepl [command] [option]...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *RootCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.authKey, "auth-key", "", "the authentication key as given in your DeepL account.")
	fs.StringVar(&c.serverURL, "server-url", "", "an alternative server URL.")

	fs.BoolFunc("v", "increase output verbosity", func(s string) error {
		c.verbosity++
		return nil
	})
	// fs.Var(&c.output, "o", "Write to file instead of stdout.")
}

func (c *RootCmdConfig) Exec(ctx context.Context, args []string) error {
	return flag.ErrHelp
}

func newTranslator(cfg RootCmdConfig) (*deepl.Translator, error) {
	opts := []deepl.TranslatorOption{}

	if cfg.serverURL != "" {
		opts = append(opts, deepl.WithServerURL(cfg.serverURL))
	}

	return deepl.NewTranslator(cfg.authKey, opts...)
}

func Configure() *command.Command {
	stdout := os.Stdout
	stderr := os.Stderr

	var (
		translateCmd  = NewTranslateTextCmd(stdout, stderr)
		documentCmd   = NewDocumentCmd(stdout, stderr)
		glossariesCmd = NewGlossariesCmd(stdout, stderr)
		usageCmd      = NewUsageCmd(stdout, stderr)
		languagesCmd  = NewLanguagesCmd(stdout, stderr)

		versionCmd = command.DefaultVersionCommand(stdout)
	)

	root := NewRootCmd(stdout, stderr)

	root.Subcommands = []*command.Command{
		translateCmd,
		documentCmd,
		glossariesCmd,
		usageCmd,
		languagesCmd,
		versionCmd,
	}

	return root
}

func Execute() error {
	rootCmd := Configure()
	if rootCmd == nil {
		return errors.New("failed to configure root command")
	}

	args := os.Args[1:]
	opts := []command.ParseOption{
		command.WithEnvVarPrefix("DEEPL"),
	}

	if err := rootCmd.Parse(args, opts...); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		} else {
			return fmt.Errorf("error parsing arguments: %w", err)
		}
	}

	ctx := context.Background()
	return rootCmd.Run(ctx)
}
