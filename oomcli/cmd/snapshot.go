package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var snapshotGroupName string
var snapshotCmd = &cobra.Command{
	Use:   "snapshot [group_name]",
	Short: "Generate snapshots for the group",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			snapshotGroupName = args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := oomStore.Snapshot(ctx, snapshotGroupName); err != nil {
			exitf("failed to take snapshot for the group %s: %+v\n", snapshotGroupName, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
