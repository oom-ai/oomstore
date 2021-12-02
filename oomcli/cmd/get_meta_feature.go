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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
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
		if err := printFeatures(features, *getMetaOutput, *getMetaWide); err != nil {
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

func printFeatures(features types.FeatureList, output string, wide bool) error {
	switch output {
	case CSV:
		return printFeaturesInCSV(features, wide)
	case ASCIITable:
		return printFeaturesInASCIITable(features, true, wide)
	case Column:
		return printFeaturesInASCIITable(features, false, wide)
	case YAML:
		return printFeatureInYaml(features)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printFeatureInYaml(features types.FeatureList) error {
	var (
		out   []byte
		err   error
		items = apply.FromFeatureList(features)
	)

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

func printFeaturesInCSV(features types.FeatureList, wide bool) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(featureHeader(wide)); err != nil {
		return err
	}
	for _, feature := range features {
		if err := w.Write(featureRecord(feature, wide)); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printFeaturesInASCIITable(features types.FeatureList, border, wide bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(featureHeader(wide))
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
		table.Append(featureRecord(feature, wide))
	}
	table.Render()
	return nil
}

func featureHeader(wide bool) []string {
	if wide {
		return []string{"NAME", "GROUP", "ENTITY", "CATEGORY", "DB-VALUE-TYPE", "VALUE-TYPE", "DESCRIPTION", "ONLINE-REVISION-ID", "CREATE-TIME", "MODIFY-TIME"}
	}
	return []string{"NAME", "GROUP", "ENTITY", "CATEGORY", "VALUE-TYPE"}
}

func featureRecord(f *types.Feature, wide bool) []string {
	onlineRevisionID := "<NULL>"

	if f.OnlineRevisionID() != nil {
		onlineRevisionID = strconv.FormatInt(int64(*f.OnlineRevisionID()), 10)
	}

	if wide {
		return []string{f.Name, f.Group.Name, f.Entity().Name, f.Group.Category, f.DBValueType, f.ValueType, f.Description, onlineRevisionID, f.CreateTime.Format(time.RFC3339), f.ModifyTime.Format(time.RFC3339)}
	}
	return []string{f.Name, f.Group.Name, f.Entity().Name, f.Group.Category, f.ValueType}
}
