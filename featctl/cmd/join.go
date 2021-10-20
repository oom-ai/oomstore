package cmd

import (
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "point-in-time join feature values",
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
