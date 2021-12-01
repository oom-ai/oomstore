package cmd

import "github.com/spf13/cobra"

var getFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "get feature resource",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	metaCmd.AddCommand(getFeatureCmd)
}
