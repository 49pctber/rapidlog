package cmd

import (
	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Open your rapid log in your browser",
	Long:  `Renders your rapid log as an HTML file, then opens that file in your default browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		fname, err := rapidlog.RenderSummary()
		if err != nil {
			panic(err)
		}

		err = rapidlog.OpenFile(*fname)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
