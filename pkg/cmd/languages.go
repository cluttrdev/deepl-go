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
		if cmd.Flags().Changed("glossary") {
			languagePairs, err := translator.GetGlossaryLanguagePairs()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Language pairs supported for glossaries: (source, target)")
			for _, pair := range languagePairs {
				fmt.Printf("%s, %s\n", pair.SourceLang, pair.TargetLang)
			}
		} else {
			if langType == "" {
				printLanguages("source")
				fmt.Println()
				printLanguages("target")
			} else {
				printLanguages(langType)
			}
		}
	},
}

func printLanguages(langType string) {
	languages, err := translator.GetLanguages(langType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s languages available:\n", cases.Title(language.English).String(langType))
	for _, lang := range languages {
		if lang.SupportsFormality {
			fmt.Printf("%s: %s (supports formality)\n", lang.Code, lang.Name)
		} else {
			fmt.Printf("%s: %s\n", lang.Code, lang.Name)
		}
	}
}

func init() {
	languagesCmd.Flags().Bool("glossary", false, "list language pairs supported for glossaries")
	languagesCmd.Flags().StringVar(&langType, "type", "source", "whether source or target languages should be listed")
	languagesCmd.MarkFlagsMutuallyExclusive("type", "glossary")

	rootCmd.AddCommand(languagesCmd)
}
