package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type editEntityOption struct {
	entityName *string
}

var editEntityOpt editEntityOption

var editEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Edit entity resources",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			exitf("argument at most one, got %d", len(args))
		} else if len(args) == 1 {
			editEntityOpt.entityName = &args[0]
		}
	},
	Run: func(execCmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		entities, err := queryEntities(ctx, oomStore, editEntityOpt.entityName)
		if err != nil {
			exit(err)
		}

		fileName, err := writeEntitiesToTempFile(ctx, oomStore, entities)
		if err != nil {
			exit(err)
		}

		if err = edit(ctx, oomStore, fileName); err != nil {
			exitf("apply failed: %+v", err)
		}
		fmt.Fprintf(os.Stderr, "applied")
	},
}

func init() {
	editCmd.AddCommand(editEntityCmd)
}

func writeEntitiesToTempFile(ctx context.Context, oomStore *oomstore.OomStore, entities types.EntityList) (string, error) {
	tempFile, err := getTempFile()
	if err != nil {
		exit(err)
	}
	defer tempFile.Close()

	if err := serializeEntitiesToWriter(ctx, tempFile, oomStore, entities, YAML); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
