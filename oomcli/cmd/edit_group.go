package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type editGroupOption struct {
	entityName *string
	groupName  *string
}

var editGroupOpt editGroupOption

var editGroupCmd = &cobra.Command{
	Use:   "group [group_name]",
	Short: "Edit group resources",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			editGroupOpt.entityName = nil
		}
		if len(args) == 1 {
			editGroupOpt.groupName = &args[0]
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
			exit(err)
		}

		fileName, err := writeGroupsToTempFile(ctx, oomStore, groups)
		if err != nil {
			exit(err)
		}

		if err = edit(ctx, oomStore, fileName); err != nil {
			exitf("apply failed: %+v", err)
		}
		fmt.Fprintln(os.Stderr, "applied")
	},
}

func init() {
	editCmd.AddCommand(editGroupCmd)

	flags := editGroupCmd.Flags()
	editGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func writeGroupsToTempFile(ctx context.Context, oomStore *oomstore.OomStore, groups types.GroupList) (string, error) {
	tempFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tempFile.Close()
	}()

	if err := outputGroup(ctx, groups, outputParams{
		writer:    tempFile,
		oomStore:  oomStore,
		outputOpt: YAML,
	}); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}
