package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var registerBatchFeatureOpt types.CreateFeatureOpt
var registerBatchFeatureCmd = &cobra.Command{
	Use:   "batch-feature <feature_name>",
	Short: "register a new batch feature",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		registerBatchFeatureOpt.FeatureName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if _, err := oomStore.CreateBatchFeature(ctx, registerBatchFeatureOpt); err != nil {
			log.Fatalf("failed registering new feature: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerBatchFeatureCmd)

	flags := registerBatchFeatureCmd.Flags()

	flags.StringVarP(&registerBatchFeatureOpt.GroupName, "group", "g", "", "feature group")
	_ = registerBatchFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&registerBatchFeatureOpt.DBValueType, "db-value-type", "", "", "feature value type in database")
	_ = registerBatchFeatureCmd.MarkFlagRequired("db-value-type")

	flags.StringVar(&registerBatchFeatureOpt.Description, "description", "", "feature description")
}
