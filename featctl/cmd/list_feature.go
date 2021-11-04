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

var listFeatureOpt types.ListFeatureOpt
var listFeatureOutput *string

var listFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "list all existing features given a specific group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureOpt.EntityName = nil
		}
		if !cmd.Flags().Changed("group") {
			listFeatureOpt.GroupName = nil
		}
		if !cmd.Flags().Changed("output") {
			listFeatureOutput = stringPtr(ASCIITable)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := oomStore.ListFeature(ctx, listFeatureOpt)
		if err != nil {
			log.Fatalf("failed listing features given option %v, error %v\n", listFeatureOpt, err)
		}

		// print features to stdout
		if err := printFeatures(features, *listFeatureOutput); err != nil {
			log.Fatalf("failed printing features, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureCmd)

	flags := listFeatureCmd.Flags()

	listFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	listFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
	listFeatureOutput = flags.StringP("output", "o", "", "output format")
}

func printFeatures(features types.FeatureList, output string) error {
	switch output {
	case CSV:
		return printFeaturesInCSV(features)
	case ASCIITable:
		return printFeaturesInASCIITable(features)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printFeaturesInCSV(features types.FeatureList) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(featureHeader()); err != nil {
		return err
	}
	for _, feature := range features {
		if err := w.Write(featureRecord(feature)); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printFeaturesInASCIITable(features types.FeatureList) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(featureHeader())
	table.SetAutoFormatHeaders(false)

	for _, feature := range features {
		table.Append(featureRecord(feature))
	}
	table.Render()
	return nil
}

func featureHeader() []string {
	return []string{"Name", "Group", "Entity", "Category", "DBValueType", "ValueType", "Description", "OnlineRevision", "OfflineLatestRevision", "OfflineLatestDataTable", "CreateTime", "ModifyTime"}
}

func featureRecord(f *types.Feature) []string {
	onlineRevision := "<NULL>"
	offlineRevision := "<NULL>"
	offlineDataTable := "<NULL>"

	if f.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*f.OnlineRevision)
	}
	if f.OfflineRevision != nil {
		offlineRevision = fmt.Sprint(*f.OfflineRevision)
	}
	if f.OfflineDataTable != nil {
		offlineDataTable = *f.OfflineDataTable
	}

	return []string{f.Name, f.GroupName, f.EntityName, f.Category, f.DBValueType, f.ValueType, f.Description, onlineRevision, offlineRevision, offlineDataTable, f.CreateTime.Format(time.RFC3339), f.ModifyTime.Format(time.RFC3339)}
}
