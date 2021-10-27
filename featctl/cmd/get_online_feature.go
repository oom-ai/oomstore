package cmd

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"sort"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

var getOnlineFeatureOpt types.GetOnlineFeatureValuesOpt

var getOnlineFeatureCmd = &cobra.Command{
	Use:   "online-feature",
	Short: "get online feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		featureValueMap, err := oneStore.GetOnlineFeatureValues(ctx, getOnlineFeatureOpt)
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
	getCmd.AddCommand(getOnlineFeatureCmd)

	flags := getOnlineFeatureCmd.Flags()

	flags.StringVarP(&getOnlineFeatureOpt.EntityKey, "entity-key", "k", "", "entity keys")
	_ = getOnlineFeatureCmd.MarkFlagRequired("entity")

	flags.StringSliceVar(&getOnlineFeatureOpt.FeatureNames, "feature", nil, "feature names")
	_ = getOnlineFeatureCmd.MarkFlagRequired("feature")
}
