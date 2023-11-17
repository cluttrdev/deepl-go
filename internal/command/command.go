package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
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

func (cmd *Command) Parse(args []string) error {
	if cmd.Name == "" {
		return errors.New("name is required")
	}
	if cmd.Flags == nil {
		cmd.Flags = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	}

	cmd.Flags.Usage = func() {
		fmt.Fprintln(cmd.Flags.Output(), DefaultUsage(cmd))
	}

	if err := cmd.Flags.Parse(args); err != nil {
		return fmt.Errorf("%s: %w", cmd.Name, err)
	}

	cmd.args = cmd.Flags.Args()

	// check for subcommands
	if len(cmd.args) > 0 {
		arg0 := cmd.args[0]
		for _, subcmd := range cmd.Subcommands {
			if strings.EqualFold(arg0, subcmd.Name) {
				cmd.selected = subcmd
				subcmd.parent = cmd
				return subcmd.Parse(cmd.args[1:])
			}
		}
	}

	// select self if no subcommand was found
	cmd.selected = cmd

	return nil
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
