package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/ethhte88/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
)

type getMetaEntityOption struct {
	entityName *string
}

var getMetaEntityOpt getMetaEntityOption

var getMetaEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "get existing entity given specific conditions",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		} else if len(args) == 1 {
			getMetaEntityOpt.entityName = &args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := outputEntity(ctx, os.Stdout, *getMetaOutput, oomStore, getMetaEntityOpt.entityName); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func outputEntity(ctx context.Context, out io.Writer, outputOpt string, oomStore *oomstore.OomStore, entityName *string) error {
	entities, err := oomStore.ListEntity(ctx)
	if err != nil {
		return fmt.Errorf("failed getting entities, error %v\n", err)
	}

	if entityName != nil {
		if entities = entities.Filter(func(e *types.Entity) bool {
			return e.Name == *entityName
		}); len(entities) == 0 {
			return fmt.Errorf("entity '%s' not found", *entityName)
		}
	}

	switch outputOpt {
	case YAML:
		err = serializeEntitiesInYaml(ctx, out, oomStore, entities)
	default:
		err = serializeMetadata(out, entities, outputOpt, *getMetaWide)
	}
	if err != nil {
		err = fmt.Errorf("failed printing entities, error: %v\n", err)
	}
	return err
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
