package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/list_feature"
	"github.com/spf13/cobra"
)

var listFeatureOpt list_feature.Option

// listFeatureCmd represents the list feature command
var listFeatureCmd = &cobra.Command{
	Use:   "list_feature feature",
	Short: "list_feature all existing features given a specific group",
	Run: func(cmd *cobra.Command, args []string) {
		listFeatureOpt.DBOption = dbOption
		list_feature.ListFeature(context.Background(), &listFeatureOpt)
	},
}

func init() {
	rootCmd.AddCommand(listFeatureCmd)

	flags := listFeatureCmd.Flags()

	flags.StringVarP(&listFeatureOpt.Group, "group", "g", "", "feature group")
	_ = listFeatureCmd.MarkFlagRequired("group")
}
