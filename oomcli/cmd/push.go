package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type PushOption struct {
	EntityKey     string
	GroupName     string
	FeatureValues map[string]string
}

var pushOpt PushOption

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push stream feature",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := oomStore.ListFeature(ctx, types.ListFeatureOpt{
			GroupNames: &[]string{pushOpt.GroupName},
		})
		if err != nil {
			exitf("failed push features: %+v\n", err)
		}

		featureValues, err := parsePushFeatureArgument(pushOpt.FeatureValues, features)
		if err != nil {
			exitf("failed push features: %+v\n", err)
		}

		if err := oomStore.Push(ctx, types.PushOpt{
			EntityKey:     pushOpt.EntityKey,
			GroupName:     pushOpt.GroupName,
			FeatureValues: featureValues,
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

	flags.StringToStringVarP(&pushOpt.FeatureValues, "feature", "f", nil, "feature name-value pairs")
	_ = pushCmd.MarkFlagRequired("feature")
}
