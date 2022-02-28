package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type GcOption struct {
	Force     bool
	UnixMilli int64
}

var gcOpt GcOption

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "gc temporary table",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		tableNames, err := oomStore.GetTemporaryTables(ctx, gcOpt.UnixMilli)
		if err != nil {
			exitf("gc failed: %+v", err)
		}
		if len(tableNames) == 0 {
			return
		}

		if !gcOpt.Force {
			fmt.Fprintln(os.Stderr, "The following tables will be deleted:")
			for _, name := range tableNames {
				fmt.Fprintln(os.Stderr, name)
			}
		} else {
			if err := oomStore.DropTemporaryTables(ctx, tableNames); err != nil {
				exitf("gc failed: %+v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(gcCmd)

	flags := gcCmd.Flags()

	flags.Int64VarP(&gcOpt.UnixMilli, "unix-milli", "u", 0, "any temporary tables before this time will be deleted")
	_ = gcCmd.MarkFlagRequired("unix-milli")

	flags.BoolVar(&gcOpt.Force, "force", false, "force run gc")
}
