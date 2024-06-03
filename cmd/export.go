package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var exportCmd = &cobra.Command{
	Use:   "export <path>",
	Short: "Export your database file to <path> to use as a backup",
	Long: `Export your database file to <path> to use as a backup.
If no path is specified, this will default to the same directory as your active rapid log.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// get backup file name and open it
		var bfname string
		if len(args) == 0 {
			bfname = filepath.Join(rapidlog.GetCacheDir(), "backup.sqlite3")
		} else {
			bfname = args[0]
		}
		backupfile, err := os.Create(bfname)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v", err))
		}
		defer backupfile.Close()

		// open database file
		dbfile, err := os.Open(rapidlog.Dbpath)
		if err != nil {
			panic(fmt.Errorf("error opening database file: %v", err))
		}
		defer dbfile.Close()

		// copy contents of database to backup file
		buf := make([]byte, 2048)
		for {
			n, err := dbfile.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}

			if _, err := backupfile.Write(buf[:n]); err != nil {
				panic(err)
			}
		}
		fmt.Printf("Successfully exported to %s\n", bfname)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
