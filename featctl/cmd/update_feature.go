package cmd

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var updateFeatureOpt types.UpdateFeatureOpt

var updateFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "update a specified feature",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		updateFeatureOpt.FeatureName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		if err := oneStore.UpdateFeature(ctx, updateFeatureOpt); err != nil {
			log.Fatalf("failed updating feature %s, err %v\n", updateFeatureOpt.FeatureName, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateFeatureCmd)

	flags := updateFeatureCmd.Flags()

	flags.StringVarP(&updateFeatureOpt.NewDescription, "description", "d", "", "new feature description")
	_ = updateFeatureCmd.MarkFlagRequired("description")
}
