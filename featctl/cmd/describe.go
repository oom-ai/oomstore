package cmd

import (
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "show details of a specific resource",
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
