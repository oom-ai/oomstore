package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var importOpt types.ImportBatchFeaturesOpt

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import feature data from a csv file",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		log.Println("importing features ...")
		if err := oomStore.ImportBatchFeatures(ctx, importOpt); err != nil {
			log.Fatalf("failed importing features: %v\n", err)
		}

		log.Println("succeeded.")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	flags := importCmd.Flags()

	flags.StringVarP(&importOpt.GroupName, "group", "g", "", "feature group")
	_ = importCmd.MarkFlagRequired("group")

	flags.StringVar(&importOpt.Description, "description", "", "revision description")
	_ = importCmd.MarkFlagRequired("description")

	flags.StringVar(&importOpt.DataSource.FilePath, "input-file", "", "input csv file")
	_ = importCmd.MarkFlagRequired("input-file")
	flags.StringVar(&importOpt.DataSource.Delimiter, "delimiter", ",", "specify field delimiter")
}
