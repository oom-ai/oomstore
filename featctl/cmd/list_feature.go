package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/oom-ai/oomstore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var listFeatureOpt types.ListFeatureOpt

var listFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "list all existing features given a specific group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureOpt.EntityName = nil
		}
		if !cmd.Flags().Changed("group") {
			listFeatureOpt.GroupName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		features, err := oneStore.ListRichFeature(ctx, listFeatureOpt)
		if err != nil {
			log.Fatalf("failed listing features given option %v, error %v\n", listFeatureOpt, err)
		}

		// print csv to stdout
		fmt.Println(types.RichFeatureCsvHeader())
		for _, feature := range features {
			fmt.Println(feature.ToCsvRecord())
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureCmd)

	flags := listFeatureCmd.Flags()

	listFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	listFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
}
