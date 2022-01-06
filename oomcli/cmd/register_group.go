package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var registerGroupOpt types.CreateGroupOpt
var registerGroupCmd = &cobra.Command{
	Use:   "group <group_name>",
	Short: "Register a new feature group",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if registerGroupOpt.Category != types.CategoryBatch && registerGroupOpt.Category != types.CategoryStream {
			log.Fatalf("illegal category '%s', should be either 'stream' or 'batch'", registerGroupOpt.Category)
		}

		registerGroupOpt.GroupName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if _, err := oomStore.CreateGroup(ctx, registerGroupOpt); err != nil {
			log.Fatalf("failed registering new group: %+v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerGroupCmd)

	flags := registerGroupCmd.Flags()

	flags.StringVarP(&registerGroupOpt.EntityName, "entity", "e", "", "entity name")
	_ = registerGroupCmd.MarkFlagRequired("entity")

	flags.StringVarP(&registerGroupOpt.Category, "category", "c", "batch", "group category")

	flags.StringVarP(&registerGroupOpt.Description, "description", "d", "", "group description")
}
