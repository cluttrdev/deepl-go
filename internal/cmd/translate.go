package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	deepl "github.com/cluttrdev/deepl-go/api"
)

var (
	// translation options
	targetLang         string
	sourceLang         string
	splitSentences     string
	preserveFormatting string
	formality          string
	glossaryId         string
	tagHandling        string
	nonSplittingTags   string
	outlineDetection   string
	splittingTags      string
	ignoreTags         string
)

var translateCmd = &cobra.Command{
	Use:   "translate [text]...",
	Short: "Translate text(s) into the target language",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		options := []deepl.TranslateOptionFunc{}

		visitor := func(flag *pflag.Flag) {
			if flag.Name == "target-lang" {
				return
			}

			if flag.Changed {
				switch flag.Name {
				case "source-lang":
					options = append(options, deepl.WithSourceLang(flag.Value.String()))
				case "split-sentences":
					options = append(options, deepl.WithSplitSentences(flag.Value.String()))
				case "preserve-formatting":
					options = append(options, deepl.WithPreserveFormatting(flag.Value.String()))
				case "formality":
					options = append(options, deepl.WithFormality(flag.Value.String()))
				case "glossary-id":
					options = append(options, deepl.WithGlossaryId(flag.Value.String()))
				case "tag-handling":
					options = append(options, deepl.WithTagHandling(flag.Value.String()))
				case "non-splitting-tags":
					options = append(options, deepl.WithNonSplittingTags(flag.Value.String()))
				case "outline-detection":
					options = append(options, deepl.WithOutlineDetection(flag.Value.String()))
				case "splitting-tags":
					options = append(options, deepl.WithSplittingTags(flag.Value.String()))
				case "ignore-tags":
					options = append(options, deepl.WithIgnoreTags(flag.Value.String()))
				}
			}
		}

		cmd.LocalFlags().VisitAll(visitor)

		translations, err := translator.TranslateText(args, targetLang, options...)
		if err != nil {
			log.Fatal(err)
		}

		for _, translation := range translations {
			if v, _ := cmd.Flags().GetCount("verbose"); v > 0 {
				fmt.Printf("Detected Source Language: %s\n", translation.DetectedSourceLanguage)
			}
			fmt.Println(translation.Text)
		}
	},
}

func init() {
	translateCmd.Flags().StringVarP(&targetLang, "target-lang", "", "DE", "the language into which the text should be translated")
	translateCmd.Flags().StringVarP(&sourceLang, "source-lang", "", "", "the language to be translated")
	translateCmd.Flags().StringVarP(&splitSentences, "split-sentences", "", "", "whether to split input into sentences")
	translateCmd.Flags().StringVarP(&preserveFormatting, "preserve-formatting", "", "", "whether the engine should respect original formatting")
	translateCmd.Flags().StringVarP(&formality, "formality", "", "", "whether the engine should lean towards formal or informal language")
	translateCmd.Flags().StringVarP(&glossaryId, "glossary-id", "", "", "the glossary to use for translation")
	translateCmd.Flags().StringVarP(&tagHandling, "tag-handling", "", "", "which kind of tags to handle")
	translateCmd.Flags().StringVarP(&nonSplittingTags, "non-splitting-tags", "", "", "comma-separated list of tags which never split sentences")
	translateCmd.Flags().StringVarP(&outlineDetection, "outline-detection", "", "", "whether to automatically detect XML structure")
	translateCmd.Flags().StringVarP(&splittingTags, "splitting-tags", "", "", "comma-separated list of tags which always cause splits")
	translateCmd.Flags().StringVarP(&ignoreTags, "ignore-tags", "", "", "comma-separated list of tags that indicate text not to be translated")
	translateCmd.Flags().SortFlags = false

	rootCmd.AddCommand(translateCmd)
}
