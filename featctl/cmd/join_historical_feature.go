package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

var joinHistoricalFeatureOpt JoinHistoricalFeaturesOpt
var joinHistoricalFeatureOutput *string

var joinHistoricalFeatureCmd = &cobra.Command{
	Use:   "historical-feature",
	Short: "join training label data set with historical feature values",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("output") {
			joinHistoricalFeatureOutput = stringPtr(ASCIITable)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := joinHistoricalFeatures(ctx, oomStore, joinHistoricalFeatureOpt, *joinHistoricalFeatureOutput); err != nil {
			log.Fatalf("failed joining historical features: %v\n", err)
		}
	},
}

func init() {
	joinCmd.AddCommand(joinHistoricalFeatureCmd)

	flags := joinHistoricalFeatureCmd.Flags()

	flags.StringVar(&joinHistoricalFeatureOpt.InputFilePath, "input-file", "", "file path of training label data set")
	_ = joinHistoricalFeatureCmd.MarkFlagRequired("input-file")

	flags.StringSliceVar(&joinHistoricalFeatureOpt.FeatureNames, "feature", nil, "feature names")
	_ = joinHistoricalFeatureCmd.MarkFlagRequired("feature")

	joinHistoricalFeatureOutput = flags.StringP("output", "o", "", "output format")
}
