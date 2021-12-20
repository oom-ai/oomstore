package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type getMetaGroupOption struct {
	entityName *string
	groupName  *string
}

var getMetaGroupOpt getMetaGroupOption

var getMetaGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "get existing group given specific conditions",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaGroupOpt.entityName = nil
		}

		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		} else if len(args) == 1 {
			getMetaGroupOpt.groupName = &args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		groups, err := queryGroups(ctx, oomStore, getMetaGroupOpt.entityName, getMetaGroupOpt.groupName)
		if err != nil {
			log.Fatal(err)
		}

		if err = serializeGroupToWriter(ctx, os.Stdout, oomStore, groups, *getMetaOutput); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaGroupCmd)

	flags := getMetaGroupCmd.Flags()

	getMetaGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func queryGroups(ctx context.Context, oomStore *oomstore.OomStore, entityName, groupName *string) (types.GroupList, error) {
	var entityID *int

	if groupName != nil {
		group, err := oomStore.GetGroupByName(ctx, *groupName)
		if err != nil {
			return nil, err
		}

		if entityName != nil && group.Entity.Name != *entityName {
			return nil, fmt.Errorf("group '%s' entityName is '%s' not '%s'", *groupName, group.Entity.Name, *entityName)
		}
		return types.GroupList{group}, err
	}

	if entityName != nil {
		entity, err := oomStore.GetEntityByName(ctx, *entityName)
		if err != nil {
			return nil, fmt.Errorf("failed to get entity name='%s': %v", *entityName, err)
		}
		entityID = &entity.ID
	}

	return oomStore.ListGroup(ctx, entityID)
}

func serializeGroupToWriter(ctx context.Context, w io.Writer, oomStore *oomstore.OomStore,
	groups types.GroupList, outputOpt string) error {

	switch outputOpt {
	case YAML:
		if items, err := groupsToApplyGroupItems(ctx, oomStore, groups); err != nil {
			return err
		} else {
			return serializeInYaml(w, *items)
		}
	default:
		return serializeMetadata(w, groups, outputOpt, *getMetaWide)
	}
}
