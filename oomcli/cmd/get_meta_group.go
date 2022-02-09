package cmd

import (
	"context"
	"io"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"

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
	Use:   "group [group_name]",
	Short: "Get existing group given specific conditions",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaGroupOpt.entityName = nil
		}
		if len(args) == 1 {
			getMetaGroupOpt.groupName = &args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		listGroupOpt := types.ListGroupOpt{}
		if getMetaGroupOpt.entityName != nil {
			listGroupOpt.EntityNames = &[]string{*getMetaGroupOpt.entityName}
		}
		if getMetaGroupOpt.groupName != nil {
			listGroupOpt.GroupNames = &[]string{*getMetaGroupOpt.groupName}
		}
		groups, err := oomStore.ListGroup(ctx, listGroupOpt)
		if err != nil {
			exitf("%+v", err)
		}

		if err = outputGroup(ctx, os.Stdout, oomStore, groups, *getMetaOutput); err != nil {
			exitf("%+v", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaGroupCmd)

	flags := getMetaGroupCmd.Flags()

	getMetaGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func outputGroup(
	ctx context.Context,
	w io.Writer,
	oomStore *oomstore.OomStore,
	groups types.GroupList,
	outputOpt string) error {
	switch outputOpt {
	case YAML:
		groupNames := groups.Names()
		features, err := oomStore.ListFeature(ctx, types.ListFeatureOpt{
			GroupNames: &groupNames,
		})
		if err != nil {
			return err
		}
		groupItems := apply.BuildGroupItems(groups, features)
		if err != nil {
			return err
		}
		return serializeInYaml(w, groupItems)
	default:
		return serializeMetadata(w, groups, outputOpt, *getMetaWide)
	}
}
