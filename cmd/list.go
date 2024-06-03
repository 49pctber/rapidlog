package cmd

import (
	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List entries in your rapid log",
	Long: `Prints entries and groups them by date.

Append the -v flag to print more information about each entry, like timestamps and IDs. This is useful for editing or deleting entries.
e.g. rapidlog list -v

Specify an entry type filter using -e to list certain types of items (e.g. to-do items prepened with ., or questions prepended with ?)
e.g. rapidlog list -e ?

Specify a timeframe to search over (e.g. "12 hours", "3 days", "1 year") using the -t flag to only list the most recent entries
e.g. rapidlog list -t "3 days"`,
	Run: func(cmd *cobra.Command, args []string) {
		entry_type, err := cmd.Flags().GetString("entry")
		if err != nil {
			panic(err)
		}

		time_constraint, err := cmd.Flags().GetString("time")
		if err != nil {
			panic(err)
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			panic(err)
		}

		err = rapidlog.ListEntries(entry_type, time_constraint, verbose)
		if err != nil {
			panic(err)
		}

		_, err = rapidlog.RenderSummary()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("entry", "e", "", "Types of entries to print [.-=o?]")
	listCmd.Flags().StringP("time", "t", "", "How far into the past to search for entries")
	listCmd.Flags().BoolP("verbose", "v", false, "Enable verbose printing")
}
