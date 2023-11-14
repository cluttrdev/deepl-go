package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	deepl "github.com/cluttrdev/deepl-go/api"
	table "github.com/cluttrdev/deepl-go/internal"
)

var glossaryCmd = &cobra.Command{
	Use:   "glossary",
	Short: "Manage glossaries",
}

var glossaryCreateCmd = &cobra.Command{
	Use:   "create [name] [source_lang] [target_lang] [entry]...",
	Short: "Create a new glossary",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		sourceLang := args[1]
		targetLang := args[2]

		count := len(args) - 3
		var entries = make([]deepl.GlossaryEntry, 0, count)
		for i := 0; i < count; i++ {
			pair := strings.Split(args[3+i], "=")
			entries = append(entries, deepl.GlossaryEntry{Source: pair[0], Target: pair[1]})
		}

		glossary, err := translator.CreateGlossary(name, sourceLang, targetLang, entries)
		if err != nil {
			log.Fatal(err)
		}

		v, err := cmd.Flags().GetCount("verbose")
		if v > 0 {
			fmt.Println("Created glossary")
		}
		printGlossaries([]deepl.GlossaryInfo{*glossary}, v)
	},
}

var glossaryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available glossaries",
	Run: func(cmd *cobra.Command, args []string) {
		glossaries, err := translator.ListGlossaries()
		if err != nil {
			log.Fatal(err)
		}

		v, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			v = 0
		}

		if len(glossaries) > 0 {
			printGlossaries(glossaries, v)
		} else {
			fmt.Println("No glossaries available")
		}
	},
}

var glossaryInfoCmd = &cobra.Command{
	Use:   "info [glossary_id]",
	Short: "Retrieve meta information for a single glossary",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var glossaryId string = args[0]
		glossary, err := translator.GetGlossary(glossaryId)
		if err != nil {
			log.Fatal(err)
		}

		v, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			v = 0
		}

		printGlossaries([]deepl.GlossaryInfo{*glossary}, v)
	},
}

var glossaryDeleteCmd = &cobra.Command{
	Use:   "delete [glossary_id]",
	Short: "Delete the specified glossary",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var glossaryId string = args[0]
		err := translator.DeleteGlossary(glossaryId)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var glossaryEntriesCmd = &cobra.Command{
	Use:   "entries [glossary_id]",
	Short: "List the entries of a single glossary",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var glossaryId string = args[0]
		entries, err := translator.GetGlossaryEntries(glossaryId)
		if err != nil {
			log.Fatal(err)
		}

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			format = "tsv"
		}

		var sep rune
		if format == "tsv" {
			sep = '\t'
		} else if format == "csv" {
			sep = ','
		} else {
			log.Fatal(errors.New(fmt.Sprintf("Invalid format: %s", format)))
		}

		printGlossaryEntries(entries, sep)
	},
}

func printGlossaries(glossaries []deepl.GlossaryInfo, verbosity int) {
	if verbosity == 0 {
		for _, g := range glossaries {
			fmt.Println(g.GlossaryId)
		}
	} else {
		tbl := table.NewTable("Glossary ID", "Name", "Ready", "Source", "Target", "Count", "Created")

		for _, g := range glossaries {
			tbl.AddRow(g.GlossaryId, g.Name, fmt.Sprint(g.Ready), g.SourceLang, g.TargetLang, fmt.Sprint(g.EntryCount), g.CreationTime)
		}

		tbl.Print()
	}
}

func printGlossaryEntries(entries []deepl.GlossaryEntry, sep rune) {
	for _, entry := range entries {
		fmt.Printf("%s%c%s\n", entry.Source, sep, entry.Target)
	}
}

func init() {
	glossaryEntriesCmd.Flags().String("format", "tsv", "the requested format of the returned glossary entries")

	glossaryCmd.AddCommand(glossaryCreateCmd)
	glossaryCmd.AddCommand(glossaryListCmd)
	glossaryCmd.AddCommand(glossaryInfoCmd)
	glossaryCmd.AddCommand(glossaryDeleteCmd)
	glossaryCmd.AddCommand(glossaryEntriesCmd)

	rootCmd.AddCommand(glossaryCmd)
}
