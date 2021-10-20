package cmd

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var updateGroupOpt types.UpdateFeatureGroupOpt

var updateGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "update a specified group",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		updateGroupOpt.GroupName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		if err := oneStore.UpdateFeatureGroup(ctx, updateGroupOpt); err != nil {
			log.Fatalf("failed updating group %s, err %v\n", updateGroupOpt.GroupName, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateGroupCmd)

	flags := updateGroupCmd.Flags()

	flags.StringVarP(&updateGroupOpt.NewDescription, "description", "d", "", "new group description")
	_ = updateGroupCmd.MarkFlagRequired("description")
}
