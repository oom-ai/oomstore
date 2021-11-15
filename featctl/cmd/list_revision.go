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

type listRevisionOption struct {
	groupName *string
}

var listRevisionOpt listRevisionOption

var listRevisionCmd = &cobra.Command{
	Use:   "revision",
	Short: "list historical revisions given a specific group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("group") {
			listRevisionOpt.groupName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var groupID *int16

		if listRevisionOpt.groupName != nil {
			group, err := oomStore.GetFeatureGroupByName(ctx, *listRevisionOpt.groupName)
			if err != nil {
				log.Fatalf("failed to get feature group name=%s: %v", *listRevisionOpt.groupName, err)
			}
			groupID = &group.ID
		}

		revisions := oomStore.ListRevision(ctx, groupID)

		if err := printRevisions(revisions, *listOutput); err != nil {
			log.Fatalf("failed printing revisions, error %v\n", err)
		}

	},
}

func init() {
	listCmd.AddCommand(listRevisionCmd)

	flags := listRevisionCmd.Flags()
	listRevisionOpt.groupName = flags.StringP("group", "g", "", "feature group")
}

func printRevisions(revisions []*types.Revision, output string) error {
	switch output {
	case CSV:
		return printRevisionsInCSV(revisions)
	case ASCIITable:
		return printRevisionsInASCIITable(revisions)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printRevisionsInCSV(revisions []*types.Revision) error {
	w := csv.NewWriter(os.Stdout)

	if err := w.Write(revisionHeader()); err != nil {
		return err
	}
	for _, r := range revisions {
		if err := w.Write(revisionRecord(r)); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printRevisionsInASCIITable(revisions []*types.Revision) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(revisionHeader())
	table.SetAutoFormatHeaders(false)

	for _, revision := range revisions {
		table.Append(revisionRecord(revision))
	}
	table.Render()
	return nil
}

func revisionHeader() []string {
	return []string{"Revision", "RevisionID", "GroupName", "DataTable", "Description", "CreateTime", "ModifyTime"}
}

func revisionRecord(r *types.Revision) []string {
	return []string{strconv.Itoa(int(r.Revision)), serializeInt32(r.ID), r.Group.Name, r.DataTable, r.Description,
		r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)}
}
