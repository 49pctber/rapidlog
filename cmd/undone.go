/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// undoneCmd represents the undone command
var undoneCmd = &cobra.Command{
	Use:   "undone",
	Short: "Mark a completed todo list item as not complete",
	Long:  `Mark a completed todo list item as not complete`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entry, err := rapidlog.GetEntry(args[0])
		if err != nil {
			panic(err)
		}

		err = entry.MarkIncomplete()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Printf("\"%s\" marked as incomplete.\n", entry.Entry)
	},
}

func init() {
	rootCmd.AddCommand(undoneCmd)
}
