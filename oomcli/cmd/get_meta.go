package cmd

import "github.com/spf13/cobra"

var getMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "get metadata",
}

func init() {
	getCmd.AddCommand(getMetaCmd)
}
