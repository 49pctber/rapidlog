/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark a todo list item as done",
	Long:  `Mark a todo list item as done`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entry, err := rapidlog.GetEntry(args[0])
		if err != nil {
			panic(err)
		}

		err = entry.MarkCompleted()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Printf("\"%s\" marked as done.\n", entry.Entry)

		_, err = rapidlog.RenderSummary()
		if err != nil {
			panic(fmt.Errorf("error rendering summary: %v", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
