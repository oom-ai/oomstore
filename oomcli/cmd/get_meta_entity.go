package cmd

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
)

type getMetaEntityOption struct {
	entityName *string
}

var getMetaEntityOpt getMetaEntityOption

var getMetaEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Get existing entity given specific conditions",
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

		entities, err := queryEntities(ctx, oomStore, getMetaEntityOpt.entityName)
		if err != nil {
			log.Fatal(err)
		}

		if err = serializeEntitiesToWriter(ctx, os.Stdout, oomStore, entities, *getMetaOutput); err != nil {
			log.Fatalf("failed printing entities, error: %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func queryEntities(ctx context.Context, oomStore *oomstore.OomStore, entityName *string) (types.EntityList, error) {
	if entityName != nil {
		entity, err := oomStore.GetEntityByName(ctx, *entityName)
		return types.EntityList{entity}, err
	}
	return oomStore.ListEntity(ctx)
}

func serializeEntitiesToWriter(ctx context.Context, w io.Writer, oomStore *oomstore.OomStore,
	entities types.EntityList, outputOpt string) error {

	switch outputOpt {
	case YAML:
		// TODO: Use entity ids to filter, rather than taking them all out
		groups, err := oomStore.ListGroup(ctx, nil)
		if err != nil {
			return err
		}

		groupItems, err := groupsToApplyGroupItems(ctx, oomStore, groups)
		if err != nil {
			return err
		}

		return serializeInYaml(w, apply.FromEntityList(entities, groupItems))
	default:
		return serializeMetadata(w, entities, outputOpt, *getMetaWide)
	}
}
