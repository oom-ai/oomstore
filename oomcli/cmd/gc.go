package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "gc temporary table",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := oomStore.Gc(ctx); err != nil {
			exitf("gc failed: %+v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gcCmd)
}
