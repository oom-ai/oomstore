package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var describeFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "show details of a specific feature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreOpt)
		defer oomStore.Close()

		featureName := args[0]
		richFeature, err := oomStore.GetRichFeature(ctx, featureName)
		if err != nil {
			log.Fatalf("failed getting feature %s, err %v\n", featureName, err)
		}
		fmt.Println(richFeature.String())
	},
}

func init() {
	describeCmd.AddCommand(describeFeatureCmd)
}
