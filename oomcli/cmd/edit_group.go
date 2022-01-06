package cmd

import (
	"context"
	"log"

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
	Use:   "group",
	Short: "Edit group resources",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			editGroupOpt.entityName = nil
		}

		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		} else if len(args) == 1 {
			editGroupOpt.groupName = &args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		groups, err := queryGroups(ctx, oomStore, editGroupOpt.entityName, editGroupOpt.groupName)
		if err != nil {
			log.Fatal(err)
		}

		fileName, err := writeGroupsToTempFile(ctx, oomStore, groups)
		if err != nil {
			log.Fatal(err)
		}

		if err = edit(ctx, oomStore, fileName); err != nil {
			log.Fatalf("apply failed: %+v", err)
		}
		log.Println("applied")
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
	defer tempFile.Close()

	if err = serializeGroupToWriter(ctx, tempFile, oomStore, groups, YAML); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}
