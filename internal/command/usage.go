package command

import (
	"flag"
	"fmt"
	"strings"
	"text/tabwriter"
)

func DefaultUsage(c *Command) string {
	var b strings.Builder

	if c.ShortHelp != "" {
		fmt.Fprintf(&b, "%s\n", c.ShortHelp)
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "USAGE\n")
	if c.ShortUsage != "" {
		fmt.Fprintf(&b, "  %s\n", c.ShortUsage)
	} else {
		fmt.Fprintf(&b, "  %s\n", DefaultShortUsage(c))
	}
	fmt.Fprintf(&b, "\n")

	if c.LongHelp != "" {
		fmt.Fprintf(&b, "%s\n\n", c.LongHelp)
	}

	if len(c.Subcommands) > 0 {
		fmt.Fprintf(&b, "COMMANDS\n")
		tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
		for _, subcommand := range c.Subcommands {
			fmt.Fprintf(tw, "  %s\t%s\n", subcommand.Name, subcommand.ShortHelp)
		}
		tw.Flush()
		fmt.Fprintf(&b, "\n")
	}

	if countFlags(c.Flags) > 0 {
		fmt.Fprintf(&b, "OPTIONS\n")
		tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
		c.Flags.VisitAll(func(f *flag.Flag) {
			_, usage := flag.UnquoteUsage(f)

			if len(f.Name) == 1 {
				fmt.Fprintf(tw, "  -%s\t%s\n", f.Name, usage)
			} else {
				fmt.Fprintf(tw, "  --%s\t%s\n", f.Name, usage)
			}
		})
		tw.Flush()
		fmt.Fprintf(&b, "\n")
	}

	return strings.TrimSpace(b.String()) + "\n"
}

func DefaultShortUsage(c *Command) string {
	builder := strings.Builder{}

	u := c.Name
	p := c.parent
	for p != nil {
		u = fmt.Sprintf("%s %s", p.Name, u)
		p = p.parent
	}
	builder.WriteString(u)

	if len(c.Subcommands) > 0 {
		builder.WriteString(" [command]")
	}

	if countFlags(c.Flags) > 0 {
		builder.WriteString(" [option]...")
	}

	builder.WriteString(" [arg]...")

	return builder.String()
}

func countFlags(fs *flag.FlagSet) (n int) {
	fs.VisitAll(func(*flag.Flag) { n++ })
	return n
}
