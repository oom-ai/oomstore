package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

var getMetaFeatureOpt types.ListFeatureOpt

var getMetaFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "Get existing features given specific conditions",
	Args:  cobra.RangeArgs(0, 1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaFeatureOpt.EntityName = nil
		}
		if !cmd.Flags().Changed("group") {
			getMetaFeatureOpt.GroupName = nil
		}
		if len(args) == 1 {
			getMetaFeatureOpt.FeatureNames = &[]string{args[0]}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := queryFeatures(ctx, oomStore, getMetaFeatureOpt)
		if err != nil {
			log.Fatal(err)
		}

		if err := serializeFeatureToWriter(os.Stdout, features, *getMetaOutput); err != nil {
			log.Fatalf("failed printing features: %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaFeatureCmd)

	flags := getMetaFeatureCmd.Flags()
	getMetaFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	getMetaFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
}

func queryFeatures(ctx context.Context, oomStore *oomstore.OomStore, opt types.ListFeatureOpt) (types.FeatureList, error) {
	features, err := oomStore.ListFeature(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed getting features, error %v\n", err)
	}

	if opt.FeatureNames != nil && len(features) == 0 {
		return nil, fmt.Errorf("feature '%s' not found", (*opt.FeatureNames)[0])
	}

	return features, nil
}

func serializeFeatureToWriter(w io.Writer, features types.FeatureList, outputOpt string) error {
	switch outputOpt {
	case YAML:
		return serializeInYaml(w, apply.FromFeatureList(features))
	default:
		return serializeMetadata(w, features, outputOpt, *getMetaWide)
	}
}
