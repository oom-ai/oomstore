package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list resources from cli",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
