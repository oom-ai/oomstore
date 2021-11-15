package cmd

import (
	"context"
	"log"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var syncOpt types.SyncOpt

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync feature values from offline store to online store",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		i, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			log.Fatalf("illegal revisionID: '%s' cannot be parsed into int32", args[0])
		}
		syncOpt.RevisionId = int32(i)
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
}
