package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
