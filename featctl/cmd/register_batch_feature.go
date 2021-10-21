package cmd

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var registerBatchFeatureOpt types.CreateFeatureOpt

var registerBatchFeatureCmd = &cobra.Command{
	Use:     "batch-feature",
	Short:   "register a new batch feature",
	Example: `featctl register feature model --group device --value-type "varchar(30)" --description 'phone model'`,
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		registerBatchFeatureOpt.FeatureName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		if _, err := oneStore.CreateBatchFeature(ctx, registerBatchFeatureOpt); err != nil {
			log.Fatalf("failed registering new feature: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerBatchFeatureCmd)

	flags := registerBatchFeatureCmd.Flags()

	flags.StringVarP(&registerBatchFeatureOpt.GroupName, "group", "g", "", "feature group")
	_ = registerBatchFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&registerBatchFeatureOpt.ValueType, "value-type", "", "", "feature value type")
	_ = registerBatchFeatureCmd.MarkFlagRequired("value-type")

	flags.StringVar(&registerBatchFeatureOpt.Description, "description", "", "feature description")
}
