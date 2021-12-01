package cmd

import (
	"github.com/spf13/cobra"
)

var getMetaOutput *string

var getMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "get metadata",
}

func init() {
	getCmd.AddCommand(getMetaCmd)

	flags := getMetaCmd.PersistentFlags()
	getMetaOutput = flags.StringP("output", "o", Column, "output format [csv,ascii_table,column]")
}
