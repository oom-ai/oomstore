package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/spf13/cobra"
)

var updateEntityOpt metadata.UpdateEntityOpt

var updateEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "update a specified entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		entityName := args[0]
		entity, err := oomStore.GetEntityByName(ctx, entityName)
		if err != nil {
			log.Fatalf("failed to get entity by name=%s: %v", entityName, err)
		}
		updateEntityOpt.EntityID = entity.ID
		if err := oomStore.UpdateEntity(ctx, updateEntityOpt); err != nil {
			log.Fatalf("failed to update entity id=%d, err %v\n", updateEntityOpt.EntityID, err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateEntityCmd)

	flags := updateEntityCmd.Flags()

	flags.StringVarP(&updateEntityOpt.NewDescription, "description", "d", "", "new entity description")
	_ = updateEntityCmd.MarkFlagRequired("description")

}
