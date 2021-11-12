package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/spf13/cobra"
)

var updateFeatureOpt metadatav2.UpdateFeatureOpt

var updateFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "update a specified feature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		featureName := args[0]
		feature, err := oomStore.GetFeatureByName(ctx, featureName)
		if err != nil {
			log.Fatalf("failed to get feature name=%s: %v", featureName, err)
		}
		updateFeatureOpt.FeatureID = feature.ID

		if err := oomStore.UpdateFeature(ctx, updateFeatureOpt); err != nil {
			log.Fatalf("failed to update feature %d, err %v\n", feature.ID, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateFeatureCmd)

	flags := updateFeatureCmd.Flags()

	flags.StringVarP(&updateFeatureOpt.NewDescription, "description", "d", "", "new feature description")
	_ = updateFeatureCmd.MarkFlagRequired("description")
}
