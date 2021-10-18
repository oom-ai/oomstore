package cmd

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var registerGroupOpt types.CreateFeatureGroupOpt

var registerGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "register a new feature group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		registerGroupOpt.Name = args[0]
		ctx := context.Background()
		onestore := mustOpenOneStore(ctx, oneStoreOpt)
		if err := onestore.CreateFeatureGroup(ctx, registerGroupOpt); err != nil {
			log.Fatalf("failed registering new group: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerGroupCmd)

	flags := registerGroupCmd.Flags()

	flags.StringVarP(&registerGroupOpt.EntityName, "entity", "e", "", "entity name")
	_ = registerGroupCmd.MarkFlagRequired("entity")

	flags.StringVar(&registerGroupOpt.Category, "category", "", "group category")
	_ = registerGroupCmd.MarkFlagRequired("category")

	flags.StringVarP(&registerGroupOpt.Description, "description", "d", "", "group description")
}
