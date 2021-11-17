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

type listGroupOption struct {
	entityName *string
}

var listGroupOpt listGroupOption

var listGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "list feature groups",
	Example: `1. featctl list group
2. featctl list group --entity device
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listGroupOpt.entityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var entityID *int

		if listGroupOpt.entityName != nil {
			entity, err := oomStore.GetEntityByName(ctx, *listGroupOpt.entityName)
			if err != nil {
				log.Fatalf("failed to get entity name='%s': %v", *listGroupOpt.entityName, err)
			}
			entityID = &entity.ID
		}

		groups := oomStore.ListGroup(ctx, entityID)
		if err := printGroups(groups, *listOutput); err != nil {
			log.Fatalf("failed printing feature groups, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listGroupCmd)

	flags := listGroupCmd.Flags()

	listGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func printGroups(groups []*types.Group, output string) error {
	switch output {
	case CSV:
		return printGroupsInCSV(groups)
	case ASCIITable:
		return printGroupsInASCIITable(groups)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printGroupsInCSV(groups types.GroupList) error {
	w := csv.NewWriter(os.Stdout)

	if err := w.Write(groupHeader()); err != nil {
		return err
	}
	for _, g := range groups {
		if err := w.Write(groupRecord(g)); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func printGroupsInASCIITable(groups types.GroupList) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(groupHeader())
	table.SetAutoFormatHeaders(false)

	for _, group := range groups {
		table.Append(groupRecord(group))
	}
	table.Render()
	return nil
}

func groupHeader() []string {
	return []string{"GroupName", "GroupID", "EntityName", "Description", "OnlineRevisionID", "CreateTime", "ModifyTime"}
}

func groupRecord(g *types.Group) []string {
	onlineRevisionID := "<NULL>"
	if g.OnlineRevisionID != nil {
		onlineRevisionID = fmt.Sprint(*g.OnlineRevisionID)
	}
	return []string{g.Name, strconv.Itoa(g.ID), g.Entity.Name, g.Description, onlineRevisionID,
		g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}
}
