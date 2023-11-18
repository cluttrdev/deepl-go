package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ParseOptions struct {
	envVarEnabled bool
	envVarPrefix  string
}

type ParseOption func(*ParseOptions) error

func WithEnvVars() ParseOption {
	return func(po *ParseOptions) error {
		po.envVarEnabled = true
		return nil
	}
}

func WithEnvVarPrefix(prefix string) ParseOption {
	return func(po *ParseOptions) error {
		po.envVarEnabled = true
		po.envVarPrefix = prefix
		return nil
	}
}

func (cmd *Command) Parse(args []string, options ...ParseOption) error {
	if cmd.Name == "" {
		return errors.New("name is required")
	}
	if cmd.Flags == nil {
		cmd.Flags = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	}

	cmd.Flags.Usage = func() {
		fmt.Fprintln(cmd.Flags.Output(), DefaultUsage(cmd))
	}

	if err := parse(cmd.Flags, args, options...); err != nil {
		return fmt.Errorf("%s: %w", cmd.Name, err)
	}

	cmd.args = cmd.Flags.Args()

	// check for subcommands
	if len(cmd.args) > 0 {
		for _, subcmd := range cmd.Subcommands {
			if strings.EqualFold(args[0], subcmd.Name) {
				cmd.selected = subcmd
				subcmd.parent = cmd

				return subcmd.Parse(cmd.args[1:], options...)
			}
		}
	}

	// select self if no subcommand was found
	cmd.selected = cmd

	return nil
}

func parse(fs *flag.FlagSet, args []string, options ...ParseOption) error {
	var opts ParseOptions
	for _, option := range options {
		if err := option(&opts); err != nil {
			return err
		}
	}

	provided := map[string]bool{}

	// command-line flags first
	{
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("parse args: %w", err)
		}

		// mark set flags as provided
		fs.Visit(func(f *flag.Flag) {
			provided[f.Name] = true
		})
	}

	// environment variables next
	if opts.envVarEnabled {
		var visitErr error
		fs.VisitAll(func(f *flag.Flag) {
			// skip flags already provided
			if provided[f.Name] {
				return
			}

			key := getEnvVarKey(f.Name, opts.envVarPrefix)

			val := os.Getenv(key)
			if val == "" {
				return
			}

			if err := fs.Set(f.Name, val); err != nil {
				visitErr = err
			}
		})
		if visitErr != nil {
			return fmt.Errorf("parse env: %w", visitErr)
		}

		// mark set flags as provided
		fs.Visit(func(f *flag.Flag) {
			provided[f.Name] = true
		})
	}

	return nil
}

func getEnvVarKey(name string, prefix string) string {
	replacer := strings.NewReplacer(
		"-", "_",
		".", "_",
		"/", "_",
	)

	key := strings.ToUpper(name)
	key = replacer.Replace(key)
	if prefix != "" {
		key = strings.ToUpper(prefix) + "_" + key
	}

	return key
}
