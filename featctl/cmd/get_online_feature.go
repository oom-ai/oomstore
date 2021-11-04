package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var getOnlineFeatureOpt types.GetOnlineFeatureValuesOpt
var getOnlineFeatureOutput *string

var getOnlineFeatureCmd = &cobra.Command{
	Use:   "online-feature",
	Short: "get online feature values",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("output") {
			getOnlineFeatureOutput = stringPtr(ASCIITable)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		featureValueMap, err := oomStore.GetOnlineFeatureValues(ctx, getOnlineFeatureOpt)
		if err != nil {
			log.Fatalf("failed getting online features: %v", err)
		}

		if err := printOnlineFeatures(featureValueMap, *getOnlineFeatureOutput); err != nil {
			log.Fatalf("failed printing online feature values, error %v\n", err)
		}
	},
}

func init() {
	getCmd.AddCommand(getOnlineFeatureCmd)

	flags := getOnlineFeatureCmd.Flags()

	flags.StringVarP(&getOnlineFeatureOpt.EntityKey, "entity-key", "k", "", "entity keys")
	_ = getOnlineFeatureCmd.MarkFlagRequired("entity")

	flags.StringSliceVar(&getOnlineFeatureOpt.FeatureNames, "feature", nil, "feature names")
	_ = getOnlineFeatureCmd.MarkFlagRequired("feature")

	getOnlineFeatureOutput = flags.StringP("output", "o", "", "output format")
}

func printOnlineFeatures(featureValueMap types.FeatureValueMap, output string) error {
	switch output {
	case CSV:
		return printOnlineFeaturesInCSV(featureValueMap)
	case ASCIITable:
		return printOnlineFeaturesInASCIITable(featureValueMap)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printOnlineFeaturesInCSV(featureValueMap types.FeatureValueMap) error {
	w := csv.NewWriter(os.Stdout)
	header := onlineFeatureHeader(featureValueMap)
	if err := w.Write(header); err != nil {
		return err
	}

	record := onlineFeatureRecord(featureValueMap, header)

	if err := w.Write(record); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func printOnlineFeaturesInASCIITable(featureValueMap types.FeatureValueMap) error {
	table := tablewriter.NewWriter(os.Stdout)
	header := onlineFeatureHeader(featureValueMap)
	table.SetHeader(header)
	table.SetAutoFormatHeaders(false)

	record := onlineFeatureRecord(featureValueMap, header)
	table.Append(record)
	table.Render()
	return nil
}

func onlineFeatureHeader(featureValueMap types.FeatureValueMap) []string {
	header := make([]string, 0, len(featureValueMap))
	for featureNames := range featureValueMap {
		header = append(header, featureNames)
	}
	sort.Strings(header)
	return header
}

func onlineFeatureRecord(featureValueMap types.FeatureValueMap, header []string) []string {
	record := make([]string, 0, len(header))
	for _, featureName := range header {
		record = append(record, cast.ToString(featureValueMap[featureName]))
	}
	return record
}
