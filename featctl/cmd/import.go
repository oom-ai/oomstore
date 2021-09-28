package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/_import"
	"github.com/spf13/cobra"
)

var importOpt _import.Option

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import feature data from a csv file",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		importOpt.DBOption = dbOption
		_import.Import(ctx, &importOpt)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	flags := importCmd.Flags()

	flags.StringVarP(&importOpt.Group, "group", "g", "", "feature group")
	_ = importCmd.MarkFlagRequired("group")

	flags.StringVarP(&importOpt.Revision, "revision", "r", "", "data revision")
	_ = importCmd.MarkFlagRequired("revision")

	flags.StringVarP(&importOpt.SchemaTemplate, "schema-template", "s", "", "entity table schema template")
	_ = importCmd.MarkFlagRequired("schema-template")

	flags.StringVar(&importOpt.InputOption.FilePath, "input-file", "", "input csv file")
	_ = importCmd.MarkFlagRequired("input-file")

	flags.StringVar(&importOpt.Description, "description", "", "revision description")
	_ = importCmd.MarkFlagRequired("description")

	flags.StringVar(&importOpt.InputOption.Separator, "separator", ",", "specify field delimiter")
	flags.StringVar(&importOpt.InputOption.Delimiter, "delimiter", "\"", "specify quoting delimiter")
	flags.BoolVar(&importOpt.InputOption.HasHeader, "has-header", false, "indicate that the input has header row")
	flags.StringVar(&importOpt.InputOption.NullLiteral, "null-literal", "\\N", "null value literal")
	flags.BoolVar(&importOpt.InputOption.BackslashEscape, "backslash-escape", false, "whether to interpret backslash")
	flags.BoolVar(&importOpt.InputOption.TrimLastSeparator, "trim-last-separator", false, "whether to trim the last separator")
}
