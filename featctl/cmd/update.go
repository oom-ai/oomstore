package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a resource",
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
