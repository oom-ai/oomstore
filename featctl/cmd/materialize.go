package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var materializeOpt types.MaterializeOpt

var materializeCmd = &cobra.Command{
	Use:   "materialize",
	Short: "materialize feature values from offline store to online store",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		materializeOpt.GroupName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreOpt)
		defer oomStore.Close()

		log.Println("materializing features ...")
		if err := oomStore.Materialize(ctx, materializeOpt); err != nil {
			log.Fatalf("failed materializing features: %v\n", err)
		}

		log.Println("succeeded.")
	},
}

func init() {
	rootCmd.AddCommand(materializeCmd)

	flags := materializeCmd.Flags()

	flags.StringVarP(&importOpt.Description, "revision", "r", "", "group revision")
}
