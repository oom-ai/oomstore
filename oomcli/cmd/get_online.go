package cmd

import (
	"context"
	"encoding/csv"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var multiGetOpt types.OnlineMultiGetOpt
var getOnlineOutput *string

var getOnlineCmd = &cobra.Command{
	Use:   "online",
	Short: "Get online feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		featureValues, err := oomStore.OnlineMultiGet(ctx, multiGetOpt)
		if err != nil {
			exitf("failed getting online features: %+v", err)
		}

		if err := printOnlineFeatures(featureValues, *getOnlineOutput); err != nil {
			exitf("failed printing online feature values, error %+v\n", err)
		}
	},
}

func init() {
	getCmd.AddCommand(getOnlineCmd)

	flags := getOnlineCmd.Flags()

	flags.StringSliceVarP(&multiGetOpt.EntityKeys, "entity-key", "k", nil, "entity keys")
	_ = getOnlineCmd.MarkFlagRequired("entity-key")

	flags.StringSliceVar(&multiGetOpt.FeatureNames, "feature", nil, "feature full names")
	_ = getOnlineCmd.MarkFlagRequired("feature")

	getOnlineOutput = flags.StringP("output", "o", ASCIITable, "output format [csv,ascii_table]")
}

func printOnlineFeatures(featureValues map[string]*types.FeatureValues, output string) error {
	switch output {
	case CSV:
		return printOnlineFeaturesInCSV(featureValues)
	case ASCIITable:
		return printOnlineFeaturesInASCIITable(featureValues)
	default:
		return errdefs.Errorf("unsupported output format %s", output)
	}
}

func printOnlineFeaturesInCSV(featureValues map[string]*types.FeatureValues) error {
	if len(featureValues) == 0 {
		return nil
	}

	var (
		w = csv.NewWriter(os.Stdout)

		keys   = entityKeys(featureValues)
		header = append([]string{featureValues[keys[0]].EntityName}, featureValues[keys[0]].FeatureNames...)
	)

	if err := w.Write(header); err != nil {
		return err
	}
	for _, key := range keys {
		value := featureValues[key]
		record := append([]string{value.EntityKey}, cast.ToStringSlice(value.FeatureValueSlice())...)
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func printOnlineFeaturesInASCIITable(featureValues map[string]*types.FeatureValues) error {
	if len(featureValues) == 0 {
		return nil
	}

	var (
		table = tablewriter.NewWriter(os.Stdout)

		keys   = entityKeys(featureValues)
		header = append([]string{featureValues[keys[0]].EntityName}, featureValues[keys[0]].FeatureNames...)
	)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(header)

	for _, key := range keys {
		value := featureValues[key]
		record := append([]string{value.EntityKey}, cast.ToStringSlice(value.FeatureValueSlice())...)
		table.Append(record)
	}

	table.Render()
	return nil
}

func entityKeys(featureValues map[string]*types.FeatureValues) []string {
	keys := make([]string, 0, len(featureValues))
	for key := range featureValues {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
