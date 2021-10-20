package cmd

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"sort"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

var getOnlineFeaturesOpt types.GetOnlineFeatureValuesOpt

var getOnlineFeaturesCmd = &cobra.Command{
	Use:   "online-features",
	Short: "get online feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		featureValueMap, err := oneStore.GetOnlineFeatureValues(ctx, getOnlineFeaturesOpt)
		if err != nil {
			log.Fatalf("failed getting online features: %v", err)
		}

		header := []string{}
		for featureNames := range featureValueMap {
			header = append(header, featureNames)
		}
		sort.Strings(header)
		data := []string{}
		for _, featureName := range header {
			data = append(data, cast.ToString(featureValueMap[featureName]))
		}

		w := csv.NewWriter(os.Stdout)
		_ = w.Write(header)
		_ = w.Write(data)
		w.Flush()
	},
}

func init() {
	getCmd.AddCommand(getOnlineFeaturesCmd)

	flags := getOnlineFeaturesCmd.Flags()

	flags.StringVarP(&getOnlineFeaturesOpt.EntityKey, "entity-key", "k", "", "entity keys")
	_ = getOnlineFeaturesCmd.MarkFlagRequired("entity")

	flags.StringSliceVar(&getOnlineFeaturesOpt.FeatureNames, "feature", nil, "feature names")
	_ = getOnlineFeaturesCmd.MarkFlagRequired("feature")
}
