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
		err := rapidlog.CliInterface()
		if err != nil {
			panic(err)
		}

		_, err = rapidlog.RenderSummary()
		if err != nil {
			panic(fmt.Errorf("error rendering summary: %v", err))
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
