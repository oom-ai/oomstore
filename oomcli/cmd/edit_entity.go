package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/ethhte88/oomstore/pkg/oomstore/types/apply"
)

type editEntityOption struct {
	entityName *string
}

var editEntityOpt editEntityOption

var editEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "edit entity resources",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
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
			log.Fatal(err)
		}

		fileName, err := writeEntitiesToTempFile(ctx, oomStore, entities)
		if err != nil {
			log.Fatal(err)
		}

		if err = openFileByEditor(ctx, fileName); err != nil {
			log.Fatal(err)
		}

		file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()

		if err := oomStore.Apply(ctx, apply.ApplyOpt{R: file}); err != nil {
			log.Fatalf("apply failed: %v", err)
		}
		log.Println("applied")
	},
}

func init() {
	editCmd.AddCommand(editEntityCmd)
}

func writeEntitiesToTempFile(ctx context.Context, oomStore *oomstore.OomStore, entities types.EntityList) (string, error) {
	tempFile, err := getTempFile()
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()

	if err := serializeEntitiesToWriter(ctx, tempFile, oomStore, entities, YAML); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
