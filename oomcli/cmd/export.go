package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var exportOutput *string
var exportOpt types.ChannelExportOpt

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export historical features in a group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("limit") {
			exportOpt.Limit = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := export(ctx, oomStore, exportOpt, *exportOutput); err != nil {
			log.Fatalf("failed exporting features: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	flags := exportCmd.Flags()
	exportOutput = flags.StringP("output", "o", ASCIITable, "output format [csv,ascii_table]")

	flags.StringSliceVar(&exportOpt.FeatureNames, "feature", nil, "select feature names")

	flags.IntVarP(&exportOpt.RevisionID, "revision-id", "r", 0, "group revision id")
	_ = exportCmd.MarkFlagRequired("revision-id")

	exportOpt.Limit = flags.Uint64P("limit", "l", 0, "max records to export")
}
