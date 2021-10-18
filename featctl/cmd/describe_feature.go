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
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		featureName := args[0]
		richFeature, err := oneStore.GetRichFeature(ctx, featureName)
		if err != nil {
			log.Fatalf("failed getting feature %s, err %v\n", featureName, err)
		}
		fmt.Println(richFeature.String())
	},
}

func init() {
	describeCmd.AddCommand(describeFeatureCmd)
}
