package cmd

import (
	"github.com/spf13/cobra"
)

var listOutput *string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list resources",
}

func init() {
	rootCmd.AddCommand(listCmd)

	flags := listCmd.PersistentFlags()
	listOutput = flags.StringP("output", "o", ASCIITable, "output format csv [csv,ascii_table]")
}
