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
	Use:   "languages [type]",
	Short: "Retreive supported languages",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if langType == "" {
			printLanguages("source")
			fmt.Println()
			printLanguages("target")
		} else {
			printLanguages(langType)
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
	languagesCmd.Flags().StringVarP(&langType, "type", "", "source", "whether source or target languages should be listed")

	rootCmd.AddCommand(languagesCmd)
}
