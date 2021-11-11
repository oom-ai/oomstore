package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/spf13/cobra"
)

type listFeatureOption struct {
	metadatav2.ListFeatureOpt
	entityName *string
	groupName  *string
}

var listFeatureOpt listFeatureOption

var listFeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "list all existing features given a specific group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureOpt.EntityID = nil
		}
		if !cmd.Flags().Changed("group") {
			listFeatureOpt.GroupID = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if listFeatureOpt.entityName != nil {
			entity, err := oomStore.GetEntityByName(ctx, *listFeatureOpt.entityName)
			if err != nil {
				log.Fatalf("failed to get entity name=%s: %v", *listFeatureOpt.entityName, err)
			}
			listFeatureOpt.EntityID = &entity.ID
		}

		if listFeatureOpt.groupName != nil {
			group, err := oomStore.GetFeatureGroupByName(ctx, *listFeatureOpt.groupName)
			if err != nil {
				log.Fatalf("failed to get feature group name=%s: %v", *listFeatureOpt.groupName, err)
			}
			listFeatureOpt.GroupID = &group.ID
		}

		features := oomStore.ListFeature(ctx, listFeatureOpt.ListFeatureOpt)

		// print features to stdout
		if err := printFeatures(features, *listOutput); err != nil {
			log.Fatalf("failed printing features, error %v\n", err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureCmd)

	flags := listFeatureCmd.Flags()
	listFeatureOpt.entityName = flags.StringP("entity", "e", "", "entity")
	listFeatureOpt.groupName = flags.StringP("group", "g", "", "feature group")
}

func printFeatures(features typesv2.FeatureList, output string) error {
	switch output {
	case CSV:
		return printFeaturesInCSV(features)
	case ASCIITable:
		return printFeaturesInASCIITable(features)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printFeaturesInCSV(features typesv2.FeatureList) error {
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

func printFeaturesInASCIITable(features typesv2.FeatureList) error {
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
	return []string{"Name", "GroupID", "EntityID", "DBValueType", "ValueType", "Description", "OnlineRevision", "CreateTime", "ModifyTime"}
}

func featureRecord(f *typesv2.Feature) []string {
	onlineRevision := "<NULL>"

	if f.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*f.OnlineRevision())
	}

	return []string{f.Name, serializeInt16(f.GroupID), serializeInt16(f.Entity().ID), f.DBValueType, f.ValueType, f.Description, onlineRevision, f.CreateTime.Format(time.RFC3339), f.ModifyTime.Format(time.RFC3339)}
}
