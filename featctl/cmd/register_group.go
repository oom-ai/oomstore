package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/spf13/cobra"
)

type registerGroupOption struct {
	metadatav2.CreateFeatureGroupOpt
	entityName string
}

var registerGroupOpt registerGroupOption

var registerGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "register a new feature group",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		registerGroupOpt.Name = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		entity, err := oomStore.GetEntityByName(ctx, registerGroupOpt.entityName)
		if err != nil {
			log.Fatalf("failed to get entity name=%s: %v", registerGroupOpt.entityName, err)
		}
		registerGroupOpt.EntityID = entity.ID

		if _, err := oomStore.CreateFeatureGroup(ctx, registerGroupOpt.CreateFeatureGroupOpt); err != nil {
			log.Fatalf("failed registering new group: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerGroupCmd)

	flags := registerGroupCmd.Flags()

	flags.StringVarP(&registerGroupOpt.entityName, "entity", "e", "", "entity name")
	_ = registerGroupCmd.MarkFlagRequired("entity")

	flags.StringVarP(&registerGroupOpt.Description, "description", "d", "", "group description")
}
