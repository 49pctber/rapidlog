package cmd

import (
	"fmt"
	"os"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a given entry",
	Long: `Edit the specified entry given by <id>.
Typical IDs look like 2hCcN9LzWH0whbkF8vzSSdKCVfA.
IDs may be obtained by using the list command with the -v (--verbose) flag, or by clicking on the entry after running the summary command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]
		entry, err := rapidlog.GetEntry(id)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		tempFile, err := os.CreateTemp("", "tempfile-*.txt")
		if err != nil {
			fmt.Printf("Error creating temporary file: %v\n", err)
			return
		}
		defer os.Remove(tempFile.Name()) // Clean up the temporary file when done

		if _, err := tempFile.Write([]byte(entry.Entry)); err != nil {
			fmt.Printf("Error writing to temporary file: %v\n", err)
			return
		}

		if err := tempFile.Close(); err != nil {
			fmt.Printf("Error closing temporary file: %v\n", err)
			return
		}

		if err := rapidlog.OpenEditor(tempFile.Name()); err != nil {
			fmt.Printf("Error opening file in text editor: %v\n", err)
		}

		buf, err := os.ReadFile(tempFile.Name())
		if err != nil {
			fmt.Printf("Error reading temporary file: %v\n", err)
		}

		entry.Entry = string(buf)
		err = entry.Log()
		if err != nil {
			fmt.Printf("Error updating entry: %v\n", err)
		}

		fmt.Printf("Updated entry.\n%v\n", entry)

		_, err = rapidlog.RenderSummary()
		if err != nil {
			fmt.Printf("Error creating summary: %v\n", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
