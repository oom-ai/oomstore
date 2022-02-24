package cmd

import (
	"context"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var importOpt types.ImportOpt
var importCSVFileDataSource types.CsvFileDataSource
var importTableLinkDataSource types.TableLinkDataSource
var delimitre string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import feature data from a data source",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		importCSVFileDataSource.Delimiter, _ = utf8.DecodeRuneInString(delimitre)

		if !cmd.Flags().Changed("revision") {
			importOpt.Revision = nil
		}
		if importCSVFileDataSource.InputFilePath == "" && importTableLinkDataSource.TableName == "" {
			return fmt.Errorf(`required flag(s) "input-file" or "table-link" not set`)
		} else if importCSVFileDataSource.InputFilePath != "" && importTableLinkDataSource.TableName != "" {
			return fmt.Errorf(`"input-file" and "table-link" can not be set both`)
		} else if importCSVFileDataSource.InputFilePath != "" {
			importOpt.DataSourceType = types.CSV_FILE
			importOpt.CsvFileDataSource = &importCSVFileDataSource
		} else if importTableLinkDataSource.TableName != "" {
			importOpt.DataSourceType = types.TABLE_LINK
			importOpt.TableLinkDataSource = &importTableLinkDataSource
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		fmt.Println("importing features ...")
		revisionID, err := oomStore.Import(ctx, importOpt)
		if err != nil {
			exitf("failed importing features: %+v\n", err)
		}
		fmt.Fprintln(os.Stderr, "succeeded")
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

	flags.StringVar(&delimitre, "delimiter", ",", "specify field delimiter")
	importOpt.Revision = flags.Int64P("revision", "r", 0, "user-defined revision")
}
