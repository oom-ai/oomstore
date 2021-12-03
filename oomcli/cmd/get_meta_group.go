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

type getMetaGroupOption struct {
	entityName *string
}

var getMetaGroupOpt getMetaGroupOption

var getMetaGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "get existing group given specific conditions",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaGroupOpt.entityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		}

		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var entityID *int

		if getMetaGroupOpt.entityName != nil {
			entity, err := oomStore.GetEntityByName(ctx, *getMetaGroupOpt.entityName)
			if err != nil {
				log.Fatalf("failed to get entity name='%s': %v", *getMetaGroupOpt.entityName, err)
			}
			entityID = &entity.ID
		}

		groups, err := oomStore.ListGroup(ctx, entityID)
		if err != nil {
			log.Fatalf("failed getting feature groups, error %v\n", err)
		}

		if len(args) > 0 {
			if groups = groups.Filter(func(g *types.Group) bool {
				return g.Name == args[0]
			}); len(groups) == 0 {
				log.Fatalf("group '%s' not found", args[0])
			}
		}

		// print groups to stdout
		switch *getMetaOutput {
		case YAML:
			err = printGroupInYaml(ctx, oomStore, groups)
		default:
			err = serializeMetadata(groups, *getMetaOutput, *getMetaWide)
		}
		if err != nil {
			log.Fatalf("failed printing groups, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaGroupCmd)

	flags := getMetaGroupCmd.Flags()

	getMetaGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func printGroupInYaml(ctx context.Context, oomStore *oomstore.OomStore, groups types.GroupList) error {
	var (
		out   []byte
		err   error
		items *apply.GroupItems
	)

	if items, err = groupsToApplyGroupItems(ctx, oomStore, groups); err != nil {
		return err
	}

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

func groupsToApplyGroupItems(ctx context.Context, store *oomstore.OomStore, groups types.GroupList) (*apply.GroupItems, error) {
	// TODO: Use group ids to filter, rather than taking them all out
	features, err := store.ListFeature(ctx, types.ListFeatureOpt{})
	if err != nil {
		return nil, err
	}

	return apply.FromGroupList(groups, features), nil
}
