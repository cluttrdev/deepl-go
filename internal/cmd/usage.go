package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/deepl-go/internal/command"
)

type UsageCmdConfig struct {
	RootCmdConfig
}

func (c *UsageCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *UsageCmdConfig) Exec(ctx context.Context, args []string) error {
	t, err := newTranslator(c.RootCmdConfig)
	if err != nil {
		return err
	}

	usage, err := t.GetUsage()
	if err != nil {
		return err
	}

	m, err := json.Marshal(usage)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.stdout, string(m))

	return nil
}

func NewUsageCmd(stdout io.Writer, stderr io.Writer) *command.Command {
	cfg := UsageCmdConfig{
		RootCmdConfig: RootCmdConfig{
			stdout: stdout,
			stderr: stderr,
		},
	}

	fs := flag.NewFlagSet("usage", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "usage",
		ShortHelp:  "Retrieve usage information and account limits",
		ShortUsage: "deepl usage [option]...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}
