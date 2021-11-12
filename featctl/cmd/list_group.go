package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/spf13/cobra"
)

type listFeatureGroupOption struct {
	entityName *string
}

var listFeatureGroupOpt listFeatureGroupOption

var listFeatureGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "list feature groups",
	Example: `1. featctl list group
2. featctl list group --entity device
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureGroupOpt.entityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var entityID *int16

		if listFeatureGroupOpt.entityName != nil {
			entity, err := oomStore.GetEntityByName(ctx, *listFeatureGroupOpt.entityName)
			if err != nil {
				log.Fatalf("failed to get entity name='%s': %v", *listFeatureGroupOpt.entityName, err)
			}
			entityID = &entity.ID
		}

		groups := oomStore.ListFeatureGroup(ctx, entityID)
		if err := printFeatureGroups(groups, *listOutput); err != nil {
			log.Fatalf("failed printing feature groups, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureGroupCmd)

	flags := listFeatureGroupCmd.Flags()

	listFeatureGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func printFeatureGroups(groups []*typesv2.FeatureGroup, output string) error {
	switch output {
	case CSV:
		return printFeatureGroupsInCSV(groups)
	case ASCIITable:
		return printFeatureGroupsInASCIITable(groups)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printFeatureGroupsInCSV(groups typesv2.FeatureGroupList) error {
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

func printFeatureGroupsInASCIITable(groups typesv2.FeatureGroupList) error {
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
	return []string{"Name", "EntityID", "Description", "OnlineRevision", "CreateTime", "ModifyTime"}
}

func groupRecord(g *typesv2.FeatureGroup) []string {
	onlineRevision := "<NULL>"
	if g.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*g.OnlineRevision)
	}
	return []string{g.Name, serializeInt16(g.EntityID), g.Description, onlineRevision,
		g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}
}
