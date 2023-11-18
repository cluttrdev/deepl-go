package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

type Command struct {
	Name       string
	ShortHelp  string
	ShortUsage string
	LongHelp   string

	Flags *flag.FlagSet
	Exec  func(ctx context.Context, args []string) error

	Subcommands []*Command

	selected *Command
	parent   *Command
	args     []string
}

func (cmd *Command) Run(ctx context.Context) (err error) {
	if !cmd.Flags.Parsed() {
		return errors.New("not parsed")
	}

	if cmd.selected == nil {
		return errors.New("none selected")
	}

	switch {
	case cmd.selected == cmd && cmd.Exec == nil:
		return fmt.Errorf("%s: %w", cmd.Name, errors.New("no exec function"))
	case cmd.selected == cmd && cmd.Exec != nil:
		defer func() {
			if errors.Is(err, flag.ErrHelp) {
				cmd.Flags.Usage()
				err = nil
			}
		}()
		return cmd.Exec(ctx, cmd.args)
	default:
		return cmd.selected.Run(ctx)
	}
}
