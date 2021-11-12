package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type getHistoricalFeatureOption struct {
	types.ExportFeatureValuesOpt
	groupName string
}

var getHistoricalFeatureOpt getHistoricalFeatureOption

var getHistoricalFeatureCmd = &cobra.Command{
	Use:   "historical-feature",
	Short: "get historical features in a group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("limit") {
			getHistoricalFeatureOpt.Limit = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if group, err := oomStore.GetFeatureGroupByName(ctx, getHistoricalFeatureOpt.groupName); err != nil {
			log.Fatalf("failed to get feature group name=%s: %v", getHistoricalFeatureOpt.groupName, err)
		} else {
			getHistoricalFeatureOpt.GroupID = group.ID
		}

		if err := getHistoricalFeature(ctx, oomStore, getHistoricalFeatureOpt.ExportFeatureValuesOpt, *getOutput); err != nil {
			log.Fatalf("failed exporting features: %v\n", err)
		}
	},
}

func init() {
	getCmd.AddCommand(getHistoricalFeatureCmd)

	flags := getHistoricalFeatureCmd.Flags()

	flags.StringSliceVar(&getHistoricalFeatureOpt.FeatureNames, "feature", nil, "select feature names")

	flags.StringVarP(&getHistoricalFeatureOpt.groupName, "group", "g", "", "feature group name")
	_ = getHistoricalFeatureCmd.MarkFlagRequired("group")

	flags.Int32VarP(&getHistoricalFeatureOpt.RevisionID, "revision-id", "r", 0, "group revision id")
	_ = getHistoricalFeatureCmd.MarkFlagRequired("revision-id")

	getHistoricalFeatureOpt.Limit = flags.Uint64P("limit", "l", 0, "max records to export")
}
