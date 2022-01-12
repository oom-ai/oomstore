package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type PushOption struct {
	EntityKey   string
	GroupName   string
	FeaturePair map[string]string

	featureNames  []string
	featureValues []interface{}
}

var pushOpt PushOption

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push stream feature",
	PreRun: func(cmd *cobra.Command, args []string) {
		for feature, value := range pushOpt.FeaturePair {
			pushOpt.featureNames = append(pushOpt.featureNames, feature)
			pushOpt.featureValues = append(pushOpt.featureValues, value)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := oomStore.Push(ctx, types.PushOpt{
			EntityKey:     pushOpt.EntityKey,
			GroupName:     pushOpt.GroupName,
			FeatureNames:  pushOpt.featureNames,
			FeatureValues: pushOpt.featureValues,
		}); err != nil {
			exitf("failed push features: %+v\n", err)
		}

		fmt.Fprintln(os.Stderr, "succeeded")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	flags := pushCmd.Flags()

	flags.StringVarP(&pushOpt.EntityKey, "entity-key", "k", "", "entity key")
	_ = getOnlineCmd.MarkFlagRequired("entity-key")

	flags.StringVarP(&pushOpt.GroupName, "group", "g", "", "feature group")
	_ = pushCmd.MarkFlagRequired("group")

	flags.StringToStringVarP(&pushOpt.FeaturePair, "features", "f", nil, "features")
	_ = pushCmd.MarkFlagRequired("features")
}
