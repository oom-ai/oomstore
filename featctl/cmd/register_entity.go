package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var registerEntityOpt types.CreateEntityOpt

var registerEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "register a new entity",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		registerEntityOpt.EntityName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if _, err := oomStore.CreateEntity(ctx, registerEntityOpt); err != nil {
			log.Fatalf("failed registering new entity: %v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerEntityCmd)

	flags := registerEntityCmd.Flags()

	flags.IntVarP(&registerEntityOpt.Length, "length", "l", 0, "entity value length")
	_ = registerEntityCmd.MarkFlagRequired("length")

	flags.StringVarP(&registerEntityOpt.Description, "description", "d", "", "entity description")
}
