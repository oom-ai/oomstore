package cmd

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var onlineGetOpt types.OnlineMultiGetOpt

var getOnlineCmd = &cobra.Command{
	Use:   "online",
	Short: "Get online feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		featureValues, err := oomStore.OnlineMultiGet(ctx, onlineGetOpt)
		if err != nil {
			exitf("failed getting online features: %+v", err)
		}

		if err := printOnlineFeatures(featureValues, *getOutput); err != nil {
			exitf("failed printing online feature values, error %+v\n", err)
		}
	},
}

func init() {
	getCmd.AddCommand(getOnlineCmd)

	flags := getOnlineCmd.Flags()

	flags.StringSliceVarP(&onlineGetOpt.EntityKeys, "entity-key", "k", nil, "entity keys")
	_ = getOnlineCmd.MarkFlagRequired("entity-key")

	flags.StringSliceVar(&onlineGetOpt.FeatureFullNames, "feature", nil, "feature full names")
	_ = getOnlineCmd.MarkFlagRequired("feature")
}

func printOnlineFeatures(featureValues map[string]*types.FeatureValues, output string) error {
	switch output {
	case CSV:
		return printOnlineFeaturesInCSV(featureValues)
	case ASCIITable:
		return printOnlineFeaturesInASCIITable(featureValues)
	default:
		return errors.Errorf("unsupported output format %s", output)
	}
}

func printOnlineFeaturesInCSV(featureValues map[string]*types.FeatureValues) error {
	header := getHeader(featureValues)
	records := getRecord(featureValues)

	w := csv.NewWriter(os.Stdout)
	if err := w.Write(header); err != nil {
		return err
	}
	if err := w.WriteAll(records); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func printOnlineFeaturesInASCIITable(featureValues map[string]*types.FeatureValues) error {
	header := getHeader(featureValues)
	records := getRecord(featureValues)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(header)
	table.AppendBulk(records)
	table.Render()
	return nil
}

func getHeader(featureValues map[string]*types.FeatureValues) []string {
	for _, v := range featureValues {
		return append([]string{v.EntityName}, v.FeatureFullNames...)
	}
	return nil
}

func getRecord(featureValues map[string]*types.FeatureValues) (rs [][]string) {
	for _, v := range featureValues {
		rs = append(rs, append([]string{v.EntityKey}, cast.ToStringSlice(v.FeatureValueSlice())...))
	}
	return
}
