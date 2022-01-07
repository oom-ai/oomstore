package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var syncOpt types.SyncOpt

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync feature values from offline store to online store",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		fmt.Fprintln(os.Stderr, "syncing features ...")
		if err := oomStore.Sync(ctx, syncOpt); err != nil {
			exitf("failed sync features: %+v\n", err)
		}

		fmt.Fprintln(os.Stderr, "succeeded.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	flags := syncCmd.Flags()

	flags.IntVarP(&syncOpt.RevisionID, "revision-id", "r", 0, "group revision id")
	_ = syncCmd.MarkFlagRequired("revision-id")

	flags.IntVarP(&syncOpt.PurgeDelay, "purge-delay", "", 0, "wait time in seconds before purging the old revision")
}
