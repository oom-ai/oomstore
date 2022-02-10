package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type editFeatureOption struct {
	types.ListFeatureOpt
}

var editFeatureOpt editFeatureOption
var editFeatureEntityName, editFeatureGroupName *string

var editFeatureCmd = &cobra.Command{
	Use:   "feature [feature_name]",
	Short: "Edit feature resources",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("entity") {
			editFeatureOpt.EntityNames = &[]string{*editFeatureEntityName}
		}
		if cmd.Flags().Changed("group") {
			editFeatureOpt.GroupNames = &[]string{*editFeatureGroupName}
		}

		if len(args) == 1 {
			editFeatureOpt.FeatureNames = &[]string{args[0]}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := queryFeatures(ctx, oomStore, editFeatureOpt.ListFeatureOpt)
		if err != nil {
			exit(err)
		}

		fileName, err := writeFeaturesToTempFile(features)
		if err != nil {
			exit(err)
		}

		if err = edit(ctx, oomStore, fileName); err != nil {
			exitf("apply failed: %+v", err)
		}
		fmt.Fprintln(os.Stderr, "applied")
	},
}

func init() {
	editCmd.AddCommand(editFeatureCmd)

	flags := editFeatureCmd.Flags()
	editFeatureEntityName = flags.StringP("entity", "e", "", "entity")
	editFeatureGroupName = flags.StringP("group", "g", "", "feature group")
}

func writeFeaturesToTempFile(features types.FeatureList) (string, error) {
	tempFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if err = outputFeature(features, outputParams{
		writer:    tempFile,
		outputOpt: YAML,
	}); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}
