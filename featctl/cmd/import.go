package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type importOption struct {
	types.ImportBatchFeaturesOpt
	FilePath string
}

var importOpt importOption

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import feature data from a csv file",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("revision") {
			importOpt.Revision = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		file, err := os.Open(importOpt.FilePath)
		if err != nil {
			log.Fatalf("read file %s failed: %v", importOpt.FilePath, err)
		}
		defer file.Close()

		importOpt.DataSource.Reader = file

		log.Println("importing features ...")
		revisionID, err := oomStore.ImportBatchFeatures(ctx, importOpt.ImportBatchFeaturesOpt)
		if err != nil {
			log.Fatalf("failed importing features: %v\n", err)
		}
		log.Println("succeeded")
		fmt.Printf("RevisionID: %d\n", revisionID)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	flags := importCmd.Flags()

	flags.StringVarP(&importOpt.GroupName, "group", "g", "", "feature group")
	_ = importCmd.MarkFlagRequired("group")

	flags.StringVar(&importOpt.Description, "description", "", "revision description")
	_ = importCmd.MarkFlagRequired("description")

	flags.StringVar(&importOpt.FilePath, "input-file", "", "input csv file")
	_ = importCmd.MarkFlagRequired("input-file")

	flags.StringVar(&importOpt.DataSource.Delimiter, "delimiter", ",", "specify field delimiter")
	importOpt.Revision = flags.Int64P("revision", "r", 0, "user-defined revision")
}
