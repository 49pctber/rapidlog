package cmd

import (
	"fmt"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete the specified entry",
	Long: `Delete the specified entry given by <id>.

Typical IDs look like 2hCcN9LzWH0whbkF8vzSSdKCVfA.
IDs may be obtained by using the list command with the -v (--verbose) flag, or by clicking on the entry after running the summary command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]
		entry, err := rapidlog.GetEntry(id)
		if err != nil {
			fmt.Println("entry not found")
			return
		}

		fmt.Println()
		entry.Print(true)

		fmt.Println("Are you sure you want to delete this entry? (y/n)")

		resp, err := rapidlog.GetUserInput()
		if err != nil {
			panic(err)
		}

		if resp != "y" {
			fmt.Printf("\nEntry NOT deleted.\n\n")
			return
		}

		err = rapidlog.DeleteEntry(id)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nEntry deleted.\n\n")

		_, err = rapidlog.RenderSummary()
		if err != nil {
			panic(fmt.Errorf("error creating summary: %v", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
