package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Retrieve usage information and account limits",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		usage, err := translator.GetUsage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if v, _ := cmd.Flags().GetCount("verbose"); v > 0 {
			fmt.Println("Usage this billing period:")
		}
		if usage.CharacterLimit > 0 {
			fmt.Printf("Characters: %d of %d\n", usage.CharacterCount, usage.CharacterLimit)
		}
		if usage.DocumentLimit > 0 {
			fmt.Printf("Documents: %d of /%d\n", usage.DocumentCount, usage.DocumentLimit)
		}
		if usage.TeamDocumentLimit > 0 {
			fmt.Printf("Team Documents: %d of /%d\n", usage.TeamDocumentCount, usage.TeamDocumentLimit)
		}
	},
}

func init() {
	rootCmd.AddCommand(usageCmd)
}
