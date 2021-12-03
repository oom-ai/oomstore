package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		// print entities to stdout
		switch *getMetaOutput {
		case YAML:
			err = printEntitiesInYaml(ctx, oomStore, entities)
		default:
			err = serializeMetadata(entities, *getMetaOutput, *getMetaWide)
		}
		if err != nil {
			log.Fatalf("failed printing entities, error: %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func printEntitiesInYaml(ctx context.Context, store *oomstore.OomStore, entities types.EntityList) error {
	var (
		out   []byte
		items = apply.EntityItems{
			Items: make([]apply.Entity, 0, entities.Len()),
		}
	)

	// TODO: Use entitys ids to filter, rather than taking them all out
	groups, err := store.ListGroup(ctx, nil)
	if err != nil {
		return err
	}
	groupItems, err := groupsToApplyGroupItems(ctx, store, groups)
	if err != nil {
		return err
	}

	items = apply.FromEntityList(entities, groupItems)
	if len(items.Items) > 1 {
		if out, err = yaml.Marshal(items); err != nil {
			return err
		}
	} else if len(items.Items) == 1 {
		if out, err = yaml.Marshal(items.Items[0]); err != nil {
			return err
		}
	}
	fmt.Println(strings.Trim(string(out), "\n"))
	return nil
}
