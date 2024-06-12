package cmd

import (
	"fmt"
	"os"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rapidlog",
	Short: "Add entries to your rapid log.",
	Long:  `Add entries to your rapid log. Each entry must be prepended with ., -, =, o, or ? Type quit or exit to quit adding entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("What's up?")

		for {
			// get user input
			input, err := rapidlog.GetUserInput()
			if err != nil {
				panic(err)
			}

			// check if user wants to quit
			if rapidlog.Re_quit.MatchString(input) {
				// render summary on clean exit
				_, err := rapidlog.RenderSummary()
				if err != nil {
					panic(err)
				}
				break
			}

			// log a new entry
			var entry rapidlog.Entry
			err = entry.ParseString(input)
			if err != nil {
				fmt.Println("Entries must begin with -, ., o, =, or ?. Use the up arrow to edit your previous input.")
				continue
			}
			err = entry.Log(true)
			if err != nil {
				panic(err)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
