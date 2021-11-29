package cmd

import (
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a resource",
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
