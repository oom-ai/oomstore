package cmd

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var updateGroupOpt types.UpdateGroupOpt

var updateGroupCmd = &cobra.Command{
	Use:   "group <group_name>",
	Short: "Update a particular group",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		updateGroupOpt.GroupName = args[0]
		if !cmd.Flags().Changed("description") {
			updateGroupOpt.NewDescription = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := oomStore.UpdateGroup(ctx, updateGroupOpt); err != nil {
			exitf("failed updating group %s, err %+v\n", updateGroupOpt.GroupName, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateGroupCmd)

	flags := updateGroupCmd.Flags()

	updateGroupOpt.NewDescription = flags.StringP("description", "d", "", "new group description")
}
