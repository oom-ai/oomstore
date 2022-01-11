package cmd

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var onlineGetOpt types.OnlineGetOpt

var getOnlineCmd = &cobra.Command{
	Use:   "online",
	Short: "Get online feature values",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		featureValues, err := oomStore.OnlineGet(ctx, onlineGetOpt)
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

	flags.StringVarP(&onlineGetOpt.EntityKey, "entity-key", "k", "", "entity key")
	_ = getOnlineCmd.MarkFlagRequired("entity-key")

	flags.StringSliceVar(&onlineGetOpt.FeatureFullNames, "feature", nil, "feature full names")
	_ = getOnlineCmd.MarkFlagRequired("feature")
}

func printOnlineFeatures(featureValues *types.FeatureValues, output string) error {
	switch output {
	case CSV:
		return printOnlineFeaturesInCSV(featureValues)
	case ASCIITable:
		return printOnlineFeaturesInASCIITable(featureValues)
	default:
		return errdefs.Errorf("unsupported output format %s", output)
	}
}

func printOnlineFeaturesInCSV(featureValues *types.FeatureValues) error {
	header := append([]string{featureValues.EntityName}, featureValues.FeatureFullNames...)
	record := append([]string{featureValues.EntityKey}, cast.ToStringSlice(featureValues.FeatureValueSlice())...)

	w := csv.NewWriter(os.Stdout)
	if err := w.Write(header); err != nil {
		return err
	}
	if err := w.Write(record); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func printOnlineFeaturesInASCIITable(featureValues *types.FeatureValues) error {
	header := append([]string{featureValues.EntityName}, featureValues.FeatureFullNames...)
	record := append([]string{featureValues.EntityKey}, cast.ToStringSlice(featureValues.FeatureValueSlice())...)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(header)
	table.Append(record)
	table.Render()
	return nil
}
