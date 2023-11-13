package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	deepl "github.com/cluttrdev/deepl-go/api"
)

var translator *deepl.Translator

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
	var err error
	translator, err = deepl.NewTranslator(os.Getenv("DEEPL_AUTH_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().CountP("verbose", "v", "verbose output")
}
