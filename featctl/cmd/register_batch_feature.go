package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/spf13/cobra"
)

type registerBatchFeatureOption struct {
	metadatav2.CreateFeatureOpt
	groupName string
}

var registerBatchFeatureOpt registerBatchFeatureOption

var registerBatchFeatureCmd = &cobra.Command{
	Use:     "batch-feature",
	Short:   "register a new batch feature",
	Example: `featctl register feature model --group device --value-type "varchar(30)" --description 'phone model'`,
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		registerBatchFeatureOpt.Name = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		group, err := oomStore.GetFeatureGroupByName(ctx, registerBatchFeatureOpt.groupName)
		if err != nil {
			log.Fatalf("failed to get feature group name=%s: %v", registerBatchFeatureOpt.groupName, err)
		}
		registerBatchFeatureOpt.GroupID = group.ID

		if _, err := oomStore.CreateBatchFeature(ctx, registerBatchFeatureOpt.CreateFeatureOpt); err != nil {
			log.Fatalf("failed registering new feature: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerBatchFeatureCmd)

	flags := registerBatchFeatureCmd.Flags()

	flags.StringVarP(&registerBatchFeatureOpt.groupName, "group", "g", "", "feature group")
	_ = registerBatchFeatureCmd.MarkFlagRequired("group")

	flags.StringVarP(&registerBatchFeatureOpt.DBValueType, "db-value-type", "", "", "feature value type in database")
	_ = registerBatchFeatureCmd.MarkFlagRequired("db-value-type")

	flags.StringVar(&registerBatchFeatureOpt.Description, "description", "", "feature description")
}
