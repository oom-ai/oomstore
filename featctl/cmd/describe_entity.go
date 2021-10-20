package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var describeEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "show details of a specific entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		entityName := args[0]
		entity, err := oneStore.GetEntity(ctx, entityName)
		if err != nil {
			log.Fatalf("failed getting entity %s, err %v\n", entityName, err)
		}
		fmt.Println(entity.String())
	},
}

func init() {
	describeCmd.AddCommand(describeEntityCmd)
}
