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
		if !cmd.Flags().Changed("limit") {
			getHistoricalFeatureOpt.Limit = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := getHistoricalFeature(ctx, oomStore, getHistoricalFeatureOpt, *getOutput); err != nil {
			log.Fatalf("failed exporting features: %v\n", err)
		}
	},
}

func init() {
	getCmd.AddCommand(getHistoricalFeatureCmd)

	flags := getHistoricalFeatureCmd.Flags()

	flags.StringSliceVar(&getHistoricalFeatureOpt.FeatureNames, "feature", nil, "select feature names")

	flags.Int32VarP(&getHistoricalFeatureOpt.RevisionID, "revision-id", "r", 0, "group revision id")
	_ = getHistoricalFeatureCmd.MarkFlagRequired("revision-id")

	getHistoricalFeatureOpt.Limit = flags.Uint64P("limit", "l", 0, "max records to export")
}
