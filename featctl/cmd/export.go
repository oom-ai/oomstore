package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/export"
	"github.com/spf13/cobra"
)

var exportOpt export.Option

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export a group of features",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		exportOpt.DBOption = dbOption
		export.Export(ctx, &exportOpt)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	flags := exportCmd.Flags()

	flags.StringVarP(&exportOpt.Group, "group", "g", "", "feature group")
	flags.StringArrayVarP(&exportOpt.Features, "name", "n", nil, "feature name")
	flags.StringVarP(&exportOpt.OutputFile, "output-file", "o", "", "output file")
	_ = exportCmd.MarkFlagRequired("group")
	_ = exportCmd.MarkFlagRequired("output-file")
}
