package cmd

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/featctl/pkg/export"
	"github.com/spf13/cobra"
)

var exportOpt export.ExportOpt

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export a group of features",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("revision") {
			exportOpt.GroupRevision = nil
		}
		if !cmd.Flags().Changed("limit") {
			exportOpt.Limit = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		exportOpt.GroupName = args[0]
		ctx := context.Background()
		store := mustOpenOneStore(ctx, oneStoreOpt)
		if err := export.Export(ctx, store, exportOpt); err != nil {
			log.Fatalf("failed exporting features: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	flags := exportCmd.Flags()

	flags.StringSliceVarP(&exportOpt.FeatureNames, "select", "s", nil, "select feature names")
	exportOpt.Limit = flags.Uint64P("limit", "l", 0, "max records to export")
	exportOpt.GroupRevision = flags.Int64P("revision", "r", 0, "feature group revision")
}
