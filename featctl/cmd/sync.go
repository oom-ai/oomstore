package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var syncOpt types.SyncOpt

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync feature values from offline store to online store",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		syncOpt.GroupName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		log.Println("syncing features ...")
		if err := oomStore.Sync(ctx, syncOpt); err != nil {
			log.Fatalf("failed sync features: %v\n", err)
		}
		log.Println("succeeded.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	flags := syncCmd.Flags()

	flags.Int32VarP(&syncOpt.RevisionId, "revision", "r", 0, "group revision id")
	_ = syncCmd.MarkFlagRequired("revision")
}
