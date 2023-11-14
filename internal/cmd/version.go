package cmd

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		info, _ := debug.ReadBuildInfo()

		v := getVersion(info)

		verbosity, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			verbosity = 0
		}

		if verbosity == 0 {
			fmt.Printf("%s\n", v)
		} else {
			rev, err := getRevision(info)
			if err != nil {
				fmt.Printf("%s (<unknown>)", v)
			}

			if verbosity == 1 {
				fmt.Printf("%s (%s)\n", v, rev[0:8])
			} else {
				fmt.Printf("%s (%s)\n", v, rev)
			}
		}
	},
}

func getVersion(info *debug.BuildInfo) string {
	if info.Main.Version == "" {
		return "devel"
	} else {
		return info.Main.Version
	}
}

func getRevision(info *debug.BuildInfo) (string, error) {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value, nil
		}
	}
	return "", errors.New("Failed to get vcs.revision")
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
