package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	deepl "github.com/cluttrdev/deepl-go/pkg/api"
)

var translator *deepl.Translator

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
	translator = deepl.NewTranslator(os.Getenv("DEEPL_AUTH_KEY"), deepl.TranslatorOptions{})

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
