package cmd

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/ethhte88/oomstore/pkg/oomstore/types/apply"
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
			if entities = entities.Filter(func(e *types.Entity) bool {
				return e.Name == args[0]
			}); len(entities) == 0 {
				log.Fatalf("entity '%s' not found", args[0])
			}
		}

		w := os.Stdout
		switch *getMetaOutput {
		case YAML:
			err = serializeEntitiesInYaml(ctx, w, oomStore, entities)
		default:
			err = serializeMetadata(w, entities, *getMetaOutput, *getMetaWide)
		}
		if err != nil {
			log.Fatalf("failed printing entities, error: %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func serializeEntitiesInYaml(ctx context.Context, w io.Writer, store *oomstore.OomStore, entities types.EntityList) error {
	// TODO: Use entitys ids to filter, rather than taking them all out
	groups, err := store.ListGroup(ctx, nil)
	if err != nil {
		return err
	}

	groupItems, err := groupsToApplyGroupItems(ctx, store, groups)
	if err != nil {
		return err
	}

	return serializeInYaml(w, apply.FromEntityList(entities, groupItems))
}
