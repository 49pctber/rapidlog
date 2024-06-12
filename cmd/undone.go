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
	Short: "mark a completed todo list item as not complete",
	Long:  `mark a completed todo list item as not complete`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entry, err := rapidlog.GetEntry(args[0])
		if err != nil {
			panic(err)
		}

		if entry.Type != "x" {
			fmt.Printf("entry not a completed todo item (type %s, not x)\n", entry.Type)
			return
		}

		entry.Type = "."

		err = entry.Log()
		if err != nil {
			panic(err)
		}

		fmt.Printf("\"%s\" marked as incomplete.\n", entry.Entry)
	},
}

func init() {
	rootCmd.AddCommand(undoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// undoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// undoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
