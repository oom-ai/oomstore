package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

var joinHistoricalFeaturesOpt JoinHistoricalFeaturesOpt

var joinHistoricalFeaturesCmd = &cobra.Command{
	Use:   "historical-features",
	Short: "join training label data set with historical feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		if err := joinHistoricalFeatures(ctx, oneStore, joinHistoricalFeaturesOpt); err != nil {
			log.Fatalf("failed joining historical features: %v\n", err)
		}
	},
}

func init() {
	joinCmd.AddCommand(joinHistoricalFeaturesCmd)

	flags := joinHistoricalFeaturesCmd.Flags()

	flags.StringVar(&joinHistoricalFeaturesOpt.InputFilePath, "input-file", "", "file path of training label data set")
	_ = joinHistoricalFeaturesCmd.MarkFlagRequired("input-file")

	flags.StringSliceVar(&joinHistoricalFeaturesOpt.FeatureNames, "feature-names", nil, "feature names")
	_ = joinHistoricalFeaturesCmd.MarkFlagRequired("feature-names")
}
