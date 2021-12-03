package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var getMetaEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "get existing entity given specific conditions",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		}

		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		entities, err := oomStore.ListEntity(ctx)
		if err != nil {
			log.Fatalf("failed getting entities, error %v\n", err)
		}

		if len(args) > 0 {
			entities = entities.Filter(func(e *types.Entity) bool {
				return e.Name == args[0]
			})
		}
		// print entities to stdout
		if err := serializeMetadataList(entities, *getMetaOutput, *getMetaWide); err != nil {
			log.Fatalf("failed printing entities, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}
