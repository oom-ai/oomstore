package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/spf13/cobra"
)

var updateGroupOpt metadata.UpdateFeatureGroupOpt

var updateGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "update a specified group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		groupName := args[0]
		group, err := oomStore.GetFeatureGroupByName(ctx, groupName)
		if err != nil {
			log.Fatalf("failed to get feature group name=%s: %v", groupName, err)
		}
		updateGroupOpt.GroupID = group.ID

		if err := oomStore.UpdateFeatureGroup(ctx, updateGroupOpt); err != nil {
			log.Fatalf("failed updating group %d, err %v\n", group.ID, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateGroupCmd)

	flags := updateGroupCmd.Flags()

	updateGroupOpt.NewDescription = flags.StringP("description", "d", "", "new group description")
	_ = updateGroupCmd.MarkFlagRequired("description")
}
