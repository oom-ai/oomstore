package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var updateEntityOpt types.UpdateEntityOpt

var updateEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "update a specified entity",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		updateEntityOpt.EntityName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreOpt)
		defer oomStore.Close()

		if err := oomStore.UpdateEntity(ctx, updateEntityOpt); err != nil {
			log.Fatalf("failed updating entity %s, err %v\n", updateEntityOpt.EntityName, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateEntityCmd)

	flags := updateEntityCmd.Flags()

	flags.StringVarP(&updateEntityOpt.NewDescription, "description", "d", "", "new entity description")
	_ = updateEntityCmd.MarkFlagRequired("description")

}
