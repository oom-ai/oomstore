package cmd

import (
	"context"
	"os"

	"github.com/onestore-ai/onestore/featctl/pkg/export"
	"github.com/spf13/cobra"
)

var exportOpt export.Option

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

	os.Stdout.Name()
	flags.StringVarP(&exportOpt.Group, "group", "g", "", "feature group")
	flags.StringSliceVarP(&exportOpt.Features, "name", "n", nil, "feature name")
	flags.StringVarP(&exportOpt.OutputFile, "output-file", "o", os.Stdout.Name(), "output file")
	_ = exportCmd.MarkFlagRequired("group")
}
