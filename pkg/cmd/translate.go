package cmd

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	deepl "github.com/cluttrdev/deepl-go/pkg/api"
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
		options := []deepl.TranslateOption{}

		visitor := func(flag *pflag.Flag) {
			var (
				opt *deepl.TranslateOption
				err error
			)

			if flag.Name == "target-lang" {
				return
			}

			if flag.Changed {
				switch flag.Name {
				case "source-lang":
					opt, err = deepl.SourceLang(flag.Value.String())
				case "split-sentences":
					opt, err = deepl.SplitSentences(flag.Value.String())
				case "preserve-formatting":
					opt, err = deepl.PreserveFormatting(flag.Value.String())
				case "formality":
					opt, err = deepl.Formality(flag.Value.String())
				case "glossary-id":
					opt, err = deepl.GlossaryId(flag.Value.String())
				case "tag-handling":
					opt, err = deepl.TagHandling(flag.Value.String())
				case "non-splitting-tags":
					opt, err = deepl.NonSplittingTags(flag.Value.String())
				case "outline-detection":
					opt, err = deepl.OutlineDetection(flag.Value.String())
				case "splitting-tags":
					opt, err = deepl.SplittingTags(flag.Value.String())
				case "ignore-tags":
					opt, err = deepl.IgnoreTags(flag.Value.String())

				default:
					opt, err = nil, errors.Errorf("Invalid option: %s", flag.Name)
				}

				if err != nil {
					log.Fatal(err)
				} else {
					options = append(options, *opt)
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
