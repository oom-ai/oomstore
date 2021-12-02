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

type getMetaGroupOption struct {
	entityName *string
}

var getMetaGroupOpt getMetaGroupOption

var getMetaGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "get existing group given specific conditions",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaGroupOpt.entityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		}

		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var entityID *int

		if getMetaGroupOpt.entityName != nil {
			entity, err := oomStore.GetEntityByName(ctx, *getMetaGroupOpt.entityName)
			if err != nil {
				log.Fatalf("failed to get entity name='%s': %v", *getMetaGroupOpt.entityName, err)
			}
			entityID = &entity.ID
		}

		groups, err := oomStore.ListGroup(ctx, entityID)
		if err != nil {
			log.Fatalf("failed getting feature groups, error %v\n", err)
		}

		if len(args) > 0 {
			groups = groups.Filter(func(g *types.Group) bool {
				return g.Name == args[0]
			})
		}
		if err := printGroups(groups, *getMetaOutput); err != nil {
			log.Fatalf("failed printing feature groups, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaGroupCmd)

	flags := getMetaGroupCmd.Flags()

	getMetaGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
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
	return []string{"ID", "NAME", "ENTITY", "DESCRIPTION", "ONLINE-REVISION-ID", "CREATE-TIME", "MODIFY-TIME"}
}

func groupRecord(g *types.Group) []string {
	onlineRevisionID := "<NULL>"
	if g.OnlineRevisionID != nil {
		onlineRevisionID = fmt.Sprint(*g.OnlineRevisionID)
	}
	return []string{strconv.Itoa(g.ID), g.Name, g.Entity.Name, g.Description, onlineRevisionID,
		g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}
}
