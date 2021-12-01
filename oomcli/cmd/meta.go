package cmd

import "github.com/spf13/cobra"

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "get metadata resources",
}

func init() {
	getCmd.AddCommand(metaCmd)
}
