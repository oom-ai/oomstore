package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var importOpt types.ImportOpt
var importCSVFileDataSource types.CsvFileDataSource
var importTableLinkDataSource types.TableLinkDataSource

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import feature data from a data source",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("revision") {
			importOpt.Revision = nil
		}
		if importCSVFileDataSource.InputFilePath == "" && importTableLinkDataSource.TableName == "" {
			return fmt.Errorf(`required flag(s) "input-file" or "table-link" not set`)
		} else if importCSVFileDataSource.InputFilePath != "" && importTableLinkDataSource.TableName != "" {
			return fmt.Errorf(`"input-file" and "table-link" can not be set both`)
		} else if importCSVFileDataSource.InputFilePath != "" {
			importOpt.DataSource = importCSVFileDataSource
		} else if importTableLinkDataSource.TableName != "" {
			importOpt.DataSource = importTableLinkDataSource
		}
		return nil
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

	flags.StringVar(&importCSVFileDataSource.InputFilePath, "input-file", "", "input csv file")
	flags.StringVar(&importTableLinkDataSource.TableName, "table-link", "", "link to a existing data table")

	flags.StringVar(&importCSVFileDataSource.Delimiter, "delimiter", ",", "specify field delimiter")
	importOpt.Revision = flags.Int64P("revision", "r", 0, "user-defined revision")
}
