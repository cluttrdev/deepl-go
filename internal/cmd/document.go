package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	deepl "github.com/cluttrdev/deepl-go/api"
	table "github.com/cluttrdev/deepl-go/internal"
)

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Translate documents",
}

var documentUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a document for translation",
	Args:  cobra.ExactArgs(1),
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
				case "formality":
					options = append(options, deepl.WithFormality(flag.Value.String()))
				case "glossary_id":
					options = append(options, deepl.WithGlossaryId(flag.Value.String()))
				}
			}
		}

		cmd.LocalFlags().VisitAll(visitor)

		targetLang, err := cmd.LocalFlags().GetString("target-lang")
		if err != nil {
			log.Fatal(err)
		}

		document, err := translator.TranslateDocumentUpload(args[0], targetLang, options...)
		if err != nil {
			log.Fatal(err)
		}

		out, err := json.Marshal(document)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(out))
	},
}

var documentStatusCmd = &cobra.Command{
	Use:   "status [id]",
	Short: "Retrieve current status of a document translation process",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		key, err := cmd.LocalFlags().GetString("document-key")
		if err != nil {
			log.Fatal(err)
		}

		status, err := translator.TranslateDocumentStatus(id, key)
		if err != nil {
			log.Fatal(err)
		}

		verbosity, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			verbosity = 0
		}

		printDocumentStatus(*status, verbosity)
	},
}

var documentDownloadCmd = &cobra.Command{
	Use:   "download [id]",
	Short: "Download the document after translation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		key, err := cmd.LocalFlags().GetString("document-key")
		if err != nil {
			log.Fatal(err)
		}

		pr, err := translator.TranslateDocumentDownload(id, key)
		if err != nil {
			log.Fatal(err)
		}
		defer pr.Close()

		out, err := cmd.LocalFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		if out == "-" {
			_, err := io.Copy(os.Stdout, pr)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			f, err := os.Create(out)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			if _, err := f.ReadFrom(pr); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func printDocumentStatus(status deepl.DocumentStatus, verbosity int) {
	if verbosity == 0 {
		fmt.Println(status.Status)
	} else {
		hs := []string{"Document Id", "Status"}
		row := []string{status.DocumentId, status.Status}

		switch status.Status {
		case "translating":
			hs = append(hs, "Seconds Remaining")
			row = append(row, status.SecondsRemaining)
		case "done":
			hs = append(hs, "Billed Characters")
			row = append(row, fmt.Sprint(status.BilledCharacters))
		case "error":
			hs = append(hs, "Message")
			row = append(row, status.Message)
		}

		tbl := table.NewTable(hs...)
		tbl.AddRow(row...)
		tbl.Print()
	}
}

func init() {
	documentUploadCmd.Flags().String("target-lang", "DE", "the language into which the text should be translated")
	documentUploadCmd.Flags().String("source-lang", "", "language of the text to be translated")
	documentUploadCmd.Flags().String("formality", "", "whether the translated text should lean towards formal or informal language")
	documentUploadCmd.Flags().String("glossary-id", "", "the glossary to use for translation")

	documentStatusCmd.Flags().String("document-key", "", "the document encryption key")

	documentDownloadCmd.Flags().StringP("output", "o", "-", "the output to write the downloaded document to")
	documentDownloadCmd.Flags().String("document-key", "", "the document encryption key")

	documentCmd.AddCommand(documentUploadCmd)
	documentCmd.AddCommand(documentStatusCmd)
	documentCmd.AddCommand(documentDownloadCmd)

	rootCmd.AddCommand(documentCmd)
}
