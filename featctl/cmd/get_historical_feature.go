package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var getHistoricalFeatureOpt types.ExportFeatureValuesOpt

var getHistoricalFeatureCmd = &cobra.Command{
	Use:   "historical-feature",
	Short: "get historical features in a group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("revision") {
			getHistoricalFeatureOpt.GroupRevision = nil
		}
		if !cmd.Flags().Changed("limit") {
			getHistoricalFeatureOpt.Limit = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreOpt)
		defer oomStore.Close()

		if err := getHistoricalFeature(ctx, oomStore, getHistoricalFeatureOpt); err != nil {
			log.Fatalf("failed exporting features: %v\n", err)
		}
	},
}

func init() {
	getCmd.AddCommand(getHistoricalFeatureCmd)

	flags := getHistoricalFeatureCmd.Flags()

	flags.StringSliceVar(&getHistoricalFeatureOpt.FeatureNames, "feature", nil, "select feature names")

	flags.StringVarP(&getHistoricalFeatureOpt.GroupName, "group", "g", "", "feature group name")
	_ = getHistoricalFeatureCmd.MarkFlagRequired("group")

	getHistoricalFeatureOpt.Limit = flags.Uint64P("limit", "l", 0, "max records to export")
	getHistoricalFeatureOpt.GroupRevision = flags.Int64P("revision", "r", 0, "feature group revision")
}
