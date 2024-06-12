/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "List your todo list items",
	Long:  `List your todo list items`,
	Run: func(cmd *cobra.Command, args []string) {

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			panic(err)
		}

		show_completed, err := cmd.Flags().GetBool("completed")
		if err != nil {
			panic(err)
		}

		var entries []rapidlog.Entry
		var symbol string

		if show_completed {
			entries, err = rapidlog.GetEntries("x", "100 years")
			if err != nil {
				panic(err)
			}
			symbol = "☑"
		} else {
			entries, err = rapidlog.GetEntries(".", "100 years")
			if err != nil {
				panic(err)
			}
			symbol = "☐"
		}

		for _, entry := range entries {
			if verbose {
				fmt.Printf("%s %s\n  [%s]\n\n", symbol, entry.Entry, entry.Id)
			} else {
				fmt.Printf("%s %s\n", symbol, entry.Entry)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	todoCmd.Flags().BoolP("verbose", "v", false, "Enable verbose printing")
	todoCmd.Flags().BoolP("completed", "c", false, "Show completed tasks")
}
