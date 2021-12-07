package cmd

import (
	"context"
	"log"
	"os"

	"github.com/ethhte88/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
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

		tempFile, err := getTempFile()
		if err != nil {
			log.Fatal(err)
		}
		if err := outputEntity(ctx, tempFile, YAML, oomStore, editEntityOpt.entityName); err != nil {
			log.Fatal(err)
		}
		tempFileName := tempFile.Name()
		tempFile.Close()

		if err = openFileByEditor(ctx, tempFileName); err != nil {
			log.Fatal(err)
		}

		tempFile, err = os.OpenFile(tempFileName, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer tempFile.Close()

		if err := oomStore.Apply(ctx, apply.ApplyOpt{R: tempFile}); err != nil {
			log.Fatalf("apply failed: %v", err)
		}
		log.Println("applied")
	},
}

func init() {
	editCmd.AddCommand(editEntityCmd)
}
