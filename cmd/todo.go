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
	Short: "list your todo list items",
	Long:  `list your todo list items`,
	Run: func(cmd *cobra.Command, args []string) {
		rows, err := rapidlog.GetEntries(".", "100 years")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var entry rapidlog.Entry
			rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)
			fmt.Printf("\xE2\x98\x90 %s\n\t[%s]\n", entry.Entry, entry.Id)
		}
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
