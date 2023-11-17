package command

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

type VersionInfo interface {
	// The version number
	Version() string

	// The revision identifier of the commit
	Revision() string

	// The modification time of the commit, in RFC3339 format
	Time() string

	// Whether the source tree had uncommitted local changes
	Modified() bool

	// The version of the Go toolchain that built the binary
	GoVersion() string
}

type buildInfo struct {
	debug.BuildInfo
}

func (i *buildInfo) Version() string {
	return i.BuildInfo.Main.Version
}

func (i *buildInfo) Revision() string {
	for _, setting := range i.BuildInfo.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return ""
}

func (i *buildInfo) Time() string {
	for _, setting := range i.BuildInfo.Settings {
		if setting.Key == "vcs.time" {
			return setting.Value
		}
	}
	return ""
}

func (i *buildInfo) Modified() bool {
	for _, setting := range i.BuildInfo.Settings {
		if setting.Key == "vcs.modified" {
			v, err := strconv.ParseBool(setting.Value)
			if err != nil {
				return false
			}
			return v
		}
	}
	return false
}

func (i *buildInfo) GoVersion() string {
	return i.BuildInfo.GoVersion
}

func DefaultVersionInfo() VersionInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		info = &debug.BuildInfo{}
	}

	return &buildInfo{
		BuildInfo: *info,
	}
}

func DefaultVersionCommand(stdout io.Writer) *Command {
	cfg := versionCmdConfig{
		version: DefaultVersionInfo(),
		flags:   flag.NewFlagSet("version", flag.ExitOnError),
		stdout:  stdout,
	}

	if cfg.stdout == nil {
		cfg.stdout = os.Stdout
	}

	cfg.RegisterFlags(cfg.flags)

	return &Command{
		Name:      "version",
		ShortHelp: "Show version information",
		Flags:     cfg.flags,
		Exec:      cfg.Exec,
	}
}

type versionCmdConfig struct {
	version VersionInfo

	flags *flag.FlagSet

	stdout io.Writer
}

func (c *versionCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	a := fs.Bool("all", false, "print all information")
	fs.BoolVar(a, "a", false, "shorthand option for `--all`")

	n := fs.Bool("number", false, "print the version number")
	fs.BoolVar(n, "n", false, "shorthand option for `--number`")
	r := fs.Bool("revision", false, "print the commit revision identifier")
	fs.BoolVar(r, "r", false, "shorthand option for `--revision`")
	t := fs.Bool("time", false, "print the commit revision modification time")
	fs.BoolVar(t, "t", false, "shorthand option for `--time`")
	m := fs.Bool("modified", false, "print the commit revision identifier")
	fs.BoolVar(m, "m", false, "shorthand option for `--modified`")
	g := fs.Bool("go-version", false, "print the Go toolchain version")
	fs.BoolVar(g, "g", false, "shorthand option for `--go-version`")

	fs.Bool("json", false, "print information in JSON")
}

func testFlag(fs *flag.FlagSet, name string) bool {
	f := fs.Lookup(name)
	if f == nil {
		return false
	}

	v, err := strconv.ParseBool(f.Value.String())
	if err != nil {
		return false
	}

	return v
}

func (c *versionCmdConfig) Exec(ctx context.Context, args []string) error {
	if testFlag(c.flags, "json") {
		return c.writeJson()
	}
	return c.writeText()
}

func (c *versionCmdConfig) writeText() error {
	any := false
	c.flags.Visit(func(f *flag.Flag) {
		if f.Name != "number" {
			any = true
		}
	})

	all := testFlag(c.flags, "all")

	builder := strings.Builder{}

	if !any || testFlag(c.flags, "number") || all {
		builder.WriteString(c.version.Version())
	}
	if testFlag(c.flags, "revision") || all {
		builder.WriteString(fmt.Sprintf(" %s", c.version.Revision()))
	}
	if testFlag(c.flags, "time") || all {
		builder.WriteString(fmt.Sprintf(" %s", c.version.Time()))
	}
	if testFlag(c.flags, "go-version") || all {
		builder.WriteString(fmt.Sprintf(" %s", c.version.GoVersion()))
	}
	if testFlag(c.flags, "modified") || all {
		if c.version.Modified() {
			builder.WriteString(" (modified)")
		}
	}

	s := builder.String()

	_, err := fmt.Fprintln(c.stdout, strings.TrimSpace(s))
	if err != nil {
		return fmt.Errorf("error writing version information: %w", err)
	}
	return nil
}

func (c *versionCmdConfig) writeJson() error {
	data := map[string]string{
		"Version":   c.version.Version(),
		"Revision":  c.version.Revision(),
		"Time":      c.version.Time(),
		"GoVersion": c.version.GoVersion(),
		"Modified":  fmt.Sprint(c.version.Modified()),
	}

	m, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding version information: %w", err)
	}

	_, err = fmt.Fprintln(c.stdout, string(m))
	if err != nil {
		return fmt.Errorf("error writing version information: %w", err)
	}
	return nil
}
