package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var getMetaEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "get existing entity given specific conditions",
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
		if err := printEntities(entities, *getMetaOutput, *getMetaWide); err != nil {
			log.Fatalf("failed printing entities, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func printEntities(entities types.EntityList, output string, wide bool) error {
	switch output {
	case CSV:
		return printEntitiesInCSV(entities, wide)
	case ASCIITable:
		return printEntitiesInASCIITable(entities, true, wide)
	case Column:
		return printEntitiesInASCIITable(entities, false, wide)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printEntitiesInCSV(entities types.EntityList, wide bool) error {
	header, err := serializeHeader(types.Entity{}, wide)
	if err != nil {
		return err
	}
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(header); err != nil {
		return err
	}
	for _, entity := range entities {
		record, err := serializeRecord(*entity, wide)
		if err != nil {
			return err
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printEntitiesInASCIITable(entities types.EntityList, border, wide bool) error {
	header, err := serializeHeader(types.Entity{}, wide)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if !border {
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetNoWhiteSpace(true)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetTablePadding("  ")
	}

	for _, entity := range entities {
		record, err := serializeRecord(*entity, wide)
		if err != nil {
			return err
		}
		table.Append(record)
	}
	table.Render()
	return nil
}
