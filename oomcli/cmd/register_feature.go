package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var registerBatchFeatureOpt types.CreateFeatureOpt
var registerBatchFeatureValueType string

var registerBatchFeatureCmd = &cobra.Command{
	Use:   "feature <feature_name>",
	Short: "Register a new feature",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		valueType, err := types.ParseValueType(registerBatchFeatureValueType)
		if err != nil {
			log.Fatalf("failed registering new feature: %v\n", err)
		}
		registerBatchFeatureOpt.ValueType = valueType
		registerBatchFeatureOpt.FeatureName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if _, err := oomStore.CreateFeature(ctx, registerBatchFeatureOpt); err != nil {
			log.Fatalf("failed registering new feature: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerBatchFeatureCmd)

	flags := registerBatchFeatureCmd.Flags()

	flags.StringVarP(&registerBatchFeatureOpt.GroupName, "group", "g", "", "feature group")
	_ = registerBatchFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&registerBatchFeatureValueType, "value-type", "", "", "feature value type")
	_ = registerBatchFeatureCmd.MarkFlagRequired("value-type")

	flags.StringVar(&registerBatchFeatureOpt.Description, "description", "", "feature description")
}
