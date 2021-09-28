package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/describe_feature"
	"github.com/spf13/cobra"
)

var describeFeatureOpt describe_feature.Option

var describeFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "show details of a specific feature",
	Run: func(cmd *cobra.Command, args []string) {
		describeFeatureOpt.DBOption = dbOption
		describe_feature.Run(context.Background(), &describeFeatureOpt)
	},
}

func init() {
	describeCmd.AddCommand(describeFeatureCmd)

	flags := describeFeatureCmd.Flags()

	flags.StringVarP(&describeFeatureOpt.Group, "group", "g", "", "feature group")
	_ = describeFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&describeFeatureOpt.Name, "name", "n", "", "feature name")
	_ = describeFeatureCmd.MarkFlagRequired("name")
}
