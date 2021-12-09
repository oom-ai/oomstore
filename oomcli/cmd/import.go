package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var importOpt types.ImportOpt
var importDataSource types.CsvFileDataSource

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import feature data from a csv file",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("revision") {
			importOpt.Revision = nil
		}
		importOpt.DataSource = importDataSource
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		log.Println("importing features ...")
		revisionID, err := oomStore.Import(ctx, importOpt)
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

	flags.StringVar(&importDataSource.InputFilePath, "input-file", "", "input csv file")
	_ = importCmd.MarkFlagRequired("input-file")

	flags.StringVar(&importDataSource.Delimiter, "delimiter", ",", "specify field delimiter")
	importOpt.Revision = flags.Int64P("revision", "r", 0, "user-defined revision")
}
