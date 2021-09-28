package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/register_feature"
	"github.com/spf13/cobra"
)

var registerFeatureOpt register_feature.Option

var registerFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "register a new feature",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateCategory(registerFeatureOpt.Category); err != nil {
			return err
		}
		if err := validateStatus(registerFeatureOpt.Status); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		registerFeatureOpt.DBOption = dbOption
		register_feature.Run(ctx, &registerFeatureOpt)
	},
}

func init() {
	registerCmd.AddCommand(registerFeatureCmd)

	flags := registerFeatureCmd.Flags()

	flags.StringVarP(&registerFeatureOpt.Group, "group", "g", "", "feature group")
	_ = registerFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&registerFeatureOpt.Name, "name", "n", "", "feature name")
	_ = registerFeatureCmd.MarkFlagRequired("name")

	flags.StringVarP(&registerFeatureOpt.Category, "category", "c", "", "feature category")
	_ = registerFeatureCmd.MarkFlagRequired("category")

	flags.StringVar(&registerFeatureOpt.Revision, "revision", "", "current revision")
	_ = registerFeatureCmd.MarkFlagRequired("revision")

	flags.StringVar(&registerFeatureOpt.Description, "description", "", "feature description")
	_ = registerFeatureCmd.MarkFlagRequired("description")

	flags.StringVarP(&registerFeatureOpt.Status, "status", "s", "disabled", "feature status")
	flags.IntVar(&registerFeatureOpt.RevisionsLimit, "revisions-limit", 3, "feature history revisions limit")
}
