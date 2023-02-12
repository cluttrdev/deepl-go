package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	deepl "github.com/cluttrdev/deepl-go/pkg/api"
)

var client *deepl.Client

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "deepl",
	Short: "DeepL language translation",
	Long:  "",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	timeout := 10 * time.Second
	client = deepl.NewClient(deepl.BaseURLFree, os.Getenv("DEEPL_API_KEY"), timeout)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
