package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type editFeatureOption struct {
	types.ListFeatureOpt
}

var editFeatureOpt editFeatureOption

var editFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "Edit feature resources",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			editFeatureOpt.EntityName = nil
		}
		if !cmd.Flags().Changed("group") {
			editFeatureOpt.GroupName = nil
		}

		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		} else if len(args) == 1 {
			editFeatureOpt.FeatureFullNames = &[]string{args[0]}
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := queryFeatures(ctx, oomStore, editFeatureOpt.ListFeatureOpt)
		if err != nil {
			log.Fatal(err)
		}

		fileName, err := writeFeaturesToTempFile(ctx, oomStore, features)
		if err != nil {
			log.Fatal(err)
		}

		if err = edit(ctx, oomStore, fileName); err != nil {
			log.Fatalf("apply failed: %v", err)
		}
		log.Println("applied")
	},
}

func init() {
	editCmd.AddCommand(editFeatureCmd)

	flags := editFeatureCmd.Flags()
	editFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	editFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
}

func writeFeaturesToTempFile(ctx context.Context, oomStore *oomstore.OomStore, features types.FeatureList) (string, error) {
	tempFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if err = serializeFeatureToWriter(tempFile, features, YAML); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}
