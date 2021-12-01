package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var getMetaEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "get existing entity given specific confitions",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		}

		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		entities, err := oomStore.ListEntity(ctx)
		if err != nil {
			log.Fatalf("failed getting entities, error %v\n", err)
		}

		if len(args) > 0 {
			entities = entities.Filter(func(e *types.Entity) bool {
				return e.Name == args[0]
			})
		}
		// print entities to stdout
		if err := printEntities(entities, *getMetaOutput); err != nil {
			log.Fatalf("failed printing entities, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func printEntities(entities types.EntityList, output string) error {
	switch output {
	case CSV:
		return printEntitiesInCSV(entities)
	case ASCIITable:
		return printEntitiesInASCIITable(entities)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printEntitiesInCSV(entities types.EntityList) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(entityHeader()); err != nil {
		return err
	}
	for _, entity := range entities {
		if err := w.Write(entityRecord(entity)); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printEntitiesInASCIITable(entities types.EntityList) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(entityHeader())
	table.SetAutoFormatHeaders(false)

	for _, entity := range entities {
		table.Append(entityRecord(entity))
	}
	table.Render()
	return nil
}

func entityRecord(entity *types.Entity) []string {
	return []string{entity.Name, strconv.Itoa(entity.Length), entity.Description, entity.CreateTime.Format(time.RFC3339),
		entity.ModifyTime.Format(time.RFC3339)}
}

func entityHeader() []string {
	return []string{"NAME", "LENGTH", "DESCRIPTION", "CREATE-TIME", "MODIFY-TIME"}
}
