package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var getEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "get entity resource",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

	},
}

func init() {
	metaCmd.AddCommand(getEntityCmd)
}
