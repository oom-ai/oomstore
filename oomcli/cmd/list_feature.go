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

var listFeatureOpt types.ListFeatureOpt

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
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := oomStore.ListFeature(ctx, listFeatureOpt)
		if err != nil {
			log.Fatalf("failed listing features, error %v\n", err)
		}

		// print features to stdout
		if err := printFeatures(features, *listOutput); err != nil {
			log.Fatalf("failed printing features, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureCmd)

	flags := listFeatureCmd.Flags()
	listFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	listFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
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
	return []string{"Name", "Group", "Entity", "Category", "DBValueType", "ValueType", "Description", "OnlineRevisionID", "CreateTime", "ModifyTime"}
}

func featureRecord(f *types.Feature) []string {
	onlineRevisionID := "<NULL>"

	if f.OnlineRevisionID() != nil {
		onlineRevisionID = strconv.FormatInt(int64(*f.OnlineRevisionID()), 10)
	}

	return []string{f.Name, f.Group.Name, f.Entity().Name, f.Group.Category, f.DBValueType, f.ValueType, f.Description, onlineRevisionID, f.CreateTime.Format(time.RFC3339), f.ModifyTime.Format(time.RFC3339)}
}
