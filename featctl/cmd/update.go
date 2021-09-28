package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a resource from cli",
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
