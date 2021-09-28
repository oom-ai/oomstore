package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/update_feature"
	"github.com/spf13/cobra"
)

var updateFeatureOpt update_feature.Option

var updateFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "update a specified feature",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		updateFeatureOpt.RevisionChanged = flags.Changed("revision")
		updateFeatureOpt.DescriptionChanged = flags.Changed("description")
		updateFeatureOpt.StatusChanged = flags.Changed("status")
		updateFeatureOpt.RevisionsLimitChanged = flags.Changed("revisions-limit")
		updateFeatureOpt.DBOption = dbOption

		if updateFeatureOpt.StatusChanged {
			if err := validateStatus(updateFeatureOpt.Status); err != nil {
				return err
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		update_feature.Run(ctx, &updateFeatureOpt)
	},
}

func init() {
	updateCmd.AddCommand(updateFeatureCmd)

	flags := updateFeatureCmd.Flags()

	flags.StringVarP(&updateFeatureOpt.Group, "group", "g", "", "feature group")
	_ = updateFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&updateFeatureOpt.Name, "name", "n", "", "feature name")
	_ = updateFeatureCmd.MarkFlagRequired("name")

	flags.StringVar(&updateFeatureOpt.Revision, "revision", "", "current revision")
	flags.StringVar(&updateFeatureOpt.Description, "description", "", "feature description")
	flags.StringVarP(&updateFeatureOpt.Status, "status", "s", "", "feature status")
	flags.IntVar(&updateFeatureOpt.RevisionsLimit, "revisions-limit", 0, "feature history revisions limit")
}
