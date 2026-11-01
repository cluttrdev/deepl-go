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

	flags *flag.FlagSet

	formatJson    bool
	documentLimit bool
}

func (c *UsageCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)

	fs.BoolVar(&c.formatJson, "json", false, "print usage result in JSON")
	fs.BoolVar(&c.documentLimit, "document-limit", false, "print document limit")
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

	if c.formatJson {
		m, err := json.Marshal(usage)
		if err != nil {
			return err
		}
		fmt.Fprintln(c.stdout, string(m))
	} else {
		var (
			percentOfUsage float32 = 0
			count          int
			countLimit     int
		)

		if c.documentLimit {
			count = usage.DocumentCount
			countLimit = usage.DocumentLimit
		} else {
			count = usage.CharacterCount
			countLimit = usage.CharacterLimit
		}

		if countLimit > 0 {
			percentOfUsage = float32(count) / float32(countLimit) * 100
		}

		fmt.Fprintf(c.stdout, `Used: %.0f%%
  Characters used: %d
  Limit: %d
`,
			percentOfUsage,
			count,
			countLimit,
		)
	}

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
