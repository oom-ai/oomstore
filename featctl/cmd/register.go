package cmd

import (
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a resource from cli",
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
