package cmd

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var updateEntityOpt types.UpdateEntityOpt

var updateEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "update a specified entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		updateEntityOpt.EntityName = args[0]
		if err := oneStore.UpdateEntity(ctx, updateEntityOpt); err != nil {
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
