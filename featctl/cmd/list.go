package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list resources",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
