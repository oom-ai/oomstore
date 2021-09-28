package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/create_feature"
	"github.com/spf13/cobra"
)

var createFeatureOpt create_feature.Option

// featureCmd represents the feature command
var createFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "create a feature from cli",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateCategory(createFeatureOpt.Category); err != nil {
			return err
		}
		if err := validateStatus(createFeatureOpt.Status); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		createFeatureOpt.DBOption = sqlxDbOption
		create_feature.Create(ctx, &createFeatureOpt)
	},
}

func init() {
	createCmd.AddCommand(createFeatureCmd)

	flags := createFeatureCmd.Flags()

	flags.StringVarP(&createFeatureOpt.Group, "group", "g", "", "feature group")
	_ = createFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&createFeatureOpt.Name, "name", "n", "", "feature name")
	_ = createFeatureCmd.MarkFlagRequired("name")

	flags.StringVarP(&createFeatureOpt.Category, "category", "c", "", "feature category")
	_ = createFeatureCmd.MarkFlagRequired("category")

	flags.StringVar(&createFeatureOpt.Revision, "revision", "", "current revision")
	_ = createFeatureCmd.MarkFlagRequired("revision")

	flags.StringVar(&createFeatureOpt.Description, "description", "", "feature description")
	_ = createFeatureCmd.MarkFlagRequired("description")

	flags.StringVarP(&createFeatureOpt.Status, "status", "s", "disabled", "feature status")
	flags.IntVar(&createFeatureOpt.RevisionsLimit, "revisions-limit", 3, "feature history revisions limit")
}
