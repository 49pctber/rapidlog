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
	Short: "mark a todo list item as done",
	Long:  `mark a todo list item as done`,
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
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
