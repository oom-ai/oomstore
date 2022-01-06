package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var updateEntityOpt types.UpdateEntityOpt

var updateEntityCmd = &cobra.Command{
	Use:   "entity <entity_name>",
	Short: "Update a particular entity",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		updateEntityOpt.EntityName = args[0]
		if !cmd.Flags().Changed("description") {
			updateEntityOpt.NewDescription = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := oomStore.UpdateEntity(ctx, updateEntityOpt); err != nil {
			log.Fatalf("failed to update entity id=%s, err %+v\n", updateEntityOpt.EntityName, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateEntityCmd)

	flags := updateEntityCmd.Flags()

	updateEntityOpt.NewDescription = flags.StringP("description", "d", "", "new entity description")
}
