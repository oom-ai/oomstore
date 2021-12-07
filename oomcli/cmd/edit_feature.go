package cmd

import (
	"context"
	"log"
	"os"

	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/ethhte88/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
)

type editFeatureOption struct {
	types.ListFeatureOpt
}

var editFeatureOpt editFeatureOption

var editFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "edit feature resources",
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
			editFeatureOpt.FeatureNames = &[]string{args[0]}
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

		if err = openFileByEditor(ctx, fileName); err != nil {
			log.Fatal(err)
		}

		file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()

		if err := oomStore.Apply(ctx, apply.ApplyOpt{R: file}); err != nil {
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
