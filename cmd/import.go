package cmd

import (
	"fmt"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "Import another database file into your rapidlog",
	Long: `Import the database file located at <path> into your rapidlog.
This command adds all entries from <path> into your current database.
Duplicate entries will be ignored.

This command will not be needed unless you have manually made a copy of your database file.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := rapidlog.ImportDatabase(args[0])
		if err != nil {
			panic(fmt.Errorf("error importing database: %v", err))
		}

		_, err = rapidlog.RenderSummary()
		if err != nil {
			panic(fmt.Errorf("error rendering summary: %v", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
