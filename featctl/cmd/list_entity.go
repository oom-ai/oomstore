package cmd

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var listEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "list all existing entities",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreOpt)
		defer oomStore.Close()

		entities, err := oomStore.ListEntity(ctx)
		if err != nil {
			log.Fatalf("failed listing entities, error %v\n", err)
		}

		// print entities to stdout
		if err := printEntities(entities); err != nil {
			log.Fatalf("failing printing entities, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listEntityCmd)
}

func printEntities(entities []*types.Entity) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write([]string{"Name", "Length", "Description", "CreateTime", "ModifyTime"}); err != nil {
		return err
	}
	for _, entity := range entities {
		if err := w.Write([]string{entity.Name, strconv.Itoa(entity.Length), entity.Description, entity.CreateTime.Format(time.RFC3339),
			entity.ModifyTime.Format(time.RFC3339)}); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}
