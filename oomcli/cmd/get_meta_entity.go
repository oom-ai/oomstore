package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	MaxDescriptionLen = 40
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
		if err := printEntities(ctx, oomStore, entities, *getMetaOutput, *getMetaWide); err != nil {
			log.Fatalf("failed printing entities, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaEntityCmd)
}

func printEntities(ctx context.Context, store *oomstore.OomStore, entities types.EntityList, output string, wide bool) error {
	switch output {
	case CSV:
		return printEntitiesInCSV(entities, wide)
	case ASCIITable:
		return printEntitiesInASCIITable(entities, true, wide)
	case Column:
		return printEntitiesInASCIITable(entities, false, wide)
	case YAML:
		return printEntitiesInYaml(ctx, store, entities)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printEntitiesInYaml(ctx context.Context, store *oomstore.OomStore, entities types.EntityList) error {
	var (
		out   []byte
		items = apply.EntityItems{
			Items: make([]apply.Entity, 0, entities.Len()),
		}
	)

	// TODO: Use entitys ids to filter, rather than taking them all out
	groups, err := store.ListGroup(ctx, nil)
	if err != nil {
		return err
	}
	groupItems, err := groupsToApplyGroupItems(ctx, store, groups)
	if err != nil {
		return err
	}

	items = apply.FromEntityList(entities, groupItems)
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

func printEntitiesInCSV(entities types.EntityList, wide bool) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(entityHeader(wide)); err != nil {
		return err
	}
	for _, entity := range entities {
		if err := w.Write(entityRecord(entity, wide)); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printEntitiesInASCIITable(entities types.EntityList, border, wide bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(entityHeader(wide))
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
		table.Append(entityRecord(entity, wide))
	}
	table.Render()
	return nil
}

func entityRecord(entity *types.Entity, wide bool) []string {
	if wide {
		return []string{strconv.Itoa(entity.ID), entity.Name, strconv.Itoa(entity.Length), entity.Description, entity.CreateTime.Format(time.RFC3339),
			entity.ModifyTime.Format(time.RFC3339)}
	}
	desc := entity.Description
	if len(desc) > MaxDescriptionLen {
		desc = fmt.Sprintf("%s...", desc[0:MaxDescriptionLen])
	}
	return []string{strconv.Itoa(entity.ID), entity.Name, strconv.Itoa(entity.Length), desc}
}

func entityHeader(wide bool) []string {
	if wide {
		return []string{"ID", "NAME", "LENGTH", "DESCRIPTION", "CREATE-TIME", "MODIFY-TIME"}
	}
	return []string{"ID", "NAME", "LENGTH", "DESCRIPTION"}
}
