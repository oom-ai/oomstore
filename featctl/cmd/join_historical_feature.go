package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

var joinHistoricalFeatureOpt JoinHistoricalFeaturesOpt

var joinHistoricalFeatureCmd = &cobra.Command{
	Use:   "historical-feature",
	Short: "join training label data set with historical feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		if err := joinHistoricalFeatures(ctx, oneStore, joinHistoricalFeatureOpt); err != nil {
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
}
