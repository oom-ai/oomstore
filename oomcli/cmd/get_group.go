package cmd

import (
	"github.com/spf13/cobra"
)

var getGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "get group resource",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	metaCmd.AddCommand(getGroupCmd)
}
