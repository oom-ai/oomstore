package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get resources",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
