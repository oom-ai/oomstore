package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

var getMetaFeatureOpt types.ListFeatureOpt

var getMetaFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "get existing features given specific conditions",
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

		features, err := oomStore.ListFeature(ctx, getMetaFeatureOpt)
		if err != nil {
			log.Fatalf("failed getting features, error %v\n", err)
		}

		if len(args) != 0 && len(features) == 0 {
			log.Fatalf("feature '%s' not found", args[0])
		}

		w := os.Stdout
		switch *getMetaOutput {
		case YAML:
			err = serializeFeaturesInYaml(w, features)
		default:
			err = serializeMetadata(w, features, *getMetaOutput, *getMetaWide)
		}
		if err != nil {
			log.Fatalf("failed printing features, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaFeatureCmd)

	flags := getMetaFeatureCmd.Flags()
	getMetaFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	getMetaFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
}

func serializeFeaturesInYaml(w io.Writer, features types.FeatureList) error {
	var (
		out   []byte
		err   error
		items = apply.FromFeatureList(features)
	)

	if len(items.Items) > 1 {
		if out, err = yaml.Marshal(items); err != nil {
			return err
		}
	} else if len(items.Items) == 1 {
		if out, err = yaml.Marshal(items.Items[0]); err != nil {
			return err
		}
	}
	fmt.Fprintln(w, strings.Trim(string(out), "\n"))
	return nil
}
