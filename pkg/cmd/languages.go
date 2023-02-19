package cmd

import (
	"fmt"
	"log"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"
)

var langType string

var languagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "Retrieve supported languages",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			verbosity = 0
		}

		if cmd.Flags().Changed("glossary") {
			printGlossaryLanguages(verbosity)
		} else {
			if cmd.Flags().Changed("source") {
				langType = "source"
			} else if cmd.Flags().Changed("target") {
				langType = "target"
			}

			if langType == "" {
				printLanguages("source", verbosity)
				fmt.Println()
				printLanguages("target", verbosity)
			} else {
				printLanguages(langType, verbosity)
			}
		}
	},
}

func printLanguages(langType string, verbosity int) {
	languages, err := translator.GetLanguages(langType)
	if err != nil {
		log.Fatal(err)
	}

	if verbosity > 0 {
		fmt.Printf("%s languages available:\n", cases.Title(language.English).String(langType))
	}
	for _, lang := range languages {
		if lang.SupportsFormality {
			fmt.Printf("%s: %s (supports formality)\n", lang.Code, lang.Name)
		} else {
			fmt.Printf("%s: %s\n", lang.Code, lang.Name)
		}
	}
}

func printGlossaryLanguages(verbosity int) {
	languagePairs, err := translator.GetGlossaryLanguagePairs()
	if err != nil {
		log.Fatal(err)
	}

	if verbosity > 0 {
		fmt.Println("Language pairs supported for glossaries: (source, target)")
	}
	for _, pair := range languagePairs {
		fmt.Printf("%s, %s\n", pair.SourceLang, pair.TargetLang)
	}
}

func init() {
	languagesCmd.Flags().StringVar(&langType, "type", "source", "whether source or target languages should be listed")
	languagesCmd.Flags().Bool("source", false, "shorthand for --type=source")
	languagesCmd.Flags().Bool("target", false, "shorthand for --type=target")
	languagesCmd.Flags().Bool("glossary", false, "list language pairs supported for glossaries")
	languagesCmd.MarkFlagsMutuallyExclusive("type", "source", "target", "glossary")

	rootCmd.AddCommand(languagesCmd)
}
