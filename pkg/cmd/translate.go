package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

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
		flagSet := cmd.Flags()

		options := []deepl.TranslateOption{}

		if flagSet.Changed("source-lang") {
			options = append(options, deepl.SourceLang(sourceLang))
		}
		if flagSet.Changed("split-sentences") {
			options = append(options, deepl.SplitSentences(splitSentences))
		}
		if flagSet.Changed("preserveformatting") {
			options = append(options, deepl.PreserveFormatting(preserveFormatting))
		}
		if flagSet.Changed("formality") {
			options = append(options, deepl.Formality(formality))
		}
		if flagSet.Changed("glossary-id") {
			options = append(options, deepl.GlossaryId(glossaryId))
		}
		if flagSet.Changed("tag-handling") {
			options = append(options, deepl.TagHandling(tagHandling))
		}
		if flagSet.Changed("non-splitting-tags") {
			options = append(options, deepl.NonSplittingTags(nonSplittingTags))
		}
		if flagSet.Changed("outline-detection") {
			options = append(options, deepl.OutlineDetection(outlineDetection))
		}
		if flagSet.Changed("splitting-tags") {
			options = append(options, deepl.SplittingTags(splittingTags))
		}
		if flagSet.Changed("ignore-tags") {
			options = append(options, deepl.IgnoreTags(ignoreTags))
		}

		translations, err := translator.TranslateText(args, targetLang, options...)
		if err != nil {
			log.Fatal(err)
		}

		for _, translation := range translations {
			if verbose {
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
