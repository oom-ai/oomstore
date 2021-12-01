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

var getMetaFeatureOpt types.ListFeatureOpt

var getMetaFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "get existing features given specific conditions",
	Args:  cobra.RangeArgs(0, 1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaFeatureOpt.EntityName = nil
		}
		if !cmd.Flags().Changed("group") {
			getMetaFeatureOpt.GroupName = nil
		}
		if len(args) == 1 {
			getMetaFeatureOpt.FeatureNames = &[]string{args[0]}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		features, err := oomStore.ListFeature(ctx, getMetaFeatureOpt)
		if err != nil {
			log.Fatalf("failed getting features, error %v\n", err)
		}

		// print features to stdout
		if err := printFeatures(features, *getMetaOutput); err != nil {
			log.Fatalf("failed printing features, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaFeatureCmd)

	flags := getMetaFeatureCmd.Flags()
	getMetaFeatureOpt.EntityName = flags.StringP("entity", "e", "", "entity")
	getMetaFeatureOpt.GroupName = flags.StringP("group", "g", "", "feature group")
}

func printFeatures(features types.FeatureList, output string) error {
	switch output {
	case CSV:
		return printFeaturesInCSV(features)
	case ASCIITable:
		return printFeaturesInASCIITable(features, true)
	case Column:
		return printFeaturesInASCIITable(features, false)
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

func printFeaturesInASCIITable(features types.FeatureList, border bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(featureHeader())
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
