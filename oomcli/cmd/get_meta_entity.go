package cmd

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

type getMetaEntityOption struct {
	entityName *string
}

var getMetaEntityOpt getMetaEntityOption

var getMetaEntityCmd = &cobra.Command{
	Use:   "entity [entity_name]",
	Short: "Get existing entity given specific conditions",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			getMetaEntityOpt.entityName = &args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var listEntityOpt types.ListEntityOpt
		if getMetaEntityOpt.entityName != nil {
			listEntityOpt.EntityNames = &[]string{*getMetaEntityOpt.entityName}
		}
		entities, err := oomStore.ListEntity(ctx, listEntityOpt)
		if err != nil {
			exit(err)
		}

		if err = outputEntity(ctx, entities, outputParams{
			writer:    os.Stdout,
			oomStore:  oomStore,
			outputOpt: *getMetaOutput,
		}); err != nil {
			exitf("failed printing entities, error: %+v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

type outputParams struct {
	writer    io.Writer
	oomStore  *oomstore.OomStore
	outputOpt string
}

func outputEntity(ctx context.Context, entities types.EntityList, params outputParams) error {
	switch params.outputOpt {
	case YAML:
		entityNames := entities.Names()
		groups, err := params.oomStore.ListGroup(ctx, types.ListGroupOpt{
			EntityNames: &entityNames,
		})
		if err != nil {
			return err
		}
		features, err := params.oomStore.ListFeature(ctx, types.ListFeatureOpt{
			EntityNames: &entityNames,
		})
		if err != nil {
			return err
		}
		groupItems := apply.BuildGroupItems(groups, features)
		entityItems := apply.BuildEntityItems(entities, groupItems)
		return serializeInYaml(params.writer, entityItems)
	default:
		return serializeMetadata(params.writer, entities, params.outputOpt, *getMetaWide)
	}
}
