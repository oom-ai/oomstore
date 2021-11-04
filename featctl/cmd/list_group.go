package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type listFeatureGroupOption struct {
	EntityName *string
}

var listFeatureGroupOpt listFeatureGroupOption
var listFeatureGroupOutput *string

var listFeatureGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "list feature groups",
	Example: `1. featctl list group
2. featctl list group --entity device
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureGroupOpt.EntityName = nil
		}
		if !cmd.Flags().Changed("output") {
			listFeatureGroupOutput = stringPtr(ASCIITable)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		groups, err := oomStore.ListFeatureGroup(ctx, listFeatureGroupOpt.EntityName)
		if err != nil {
			log.Fatal(err)
		}
		if err := printFeatureGroups(groups, *listFeatureGroupOutput); err != nil {
			log.Fatalf("failed printing feature groups, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureGroupCmd)

	flags := listFeatureGroupCmd.Flags()

	listFeatureGroupOpt.EntityName = flags.StringP("entity", "", "", "use to filter groups")
	listFeatureGroupOutput = flags.StringP("output", "o", "", "output format")
}

func printFeatureGroups(groups []*types.FeatureGroup, output string) error {
	switch output {
	case CSV:
		return printFeatureGroupsInCSV(groups)
	case ASCIITable:
		return printFeatureGroupsInASCIITable(groups)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printFeatureGroupsInCSV(groups []*types.FeatureGroup) error {
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

func printFeatureGroupsInASCIITable(groups []*types.FeatureGroup) error {
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
	return []string{"Name", "Entity", "Description", "OnlineRevision", "OfflineLatestRevision", "OfflineLatestDataTable", "CreateTime", "ModifyTime"}
}

func groupRecord(g *types.FeatureGroup) []string {
	var onlineRevision, offlineRevision, offlineDataTable string

	onlineRevision = "<NULL>"
	offlineRevision = "<NULL>"
	offlineDataTable = "<NULL>"
	if g.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*g.OnlineRevision)
	}
	if g.OfflineRevision != nil {
		offlineRevision = fmt.Sprint(*g.OfflineRevision)
	}

	if g.OfflineDataTable != nil {
		offlineDataTable = *g.OfflineDataTable
	}
	return []string{g.Name, g.EntityName, g.Description, onlineRevision, offlineRevision, offlineDataTable,
		g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}
}
