package cmd

import (
	"github.com/spf13/cobra"
)

var getOutput *string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get resources",
}

func init() {
	rootCmd.AddCommand(getCmd)

	flags := getCmd.PersistentFlags()

	getOutput = flags.StringP("output", "o", ASCIITable, "output format [csv,ascii_table]")
}
