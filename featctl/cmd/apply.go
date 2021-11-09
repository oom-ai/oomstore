package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

var applyOpt types.ApplyOpt

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a change",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := oomStore.Apply(ctx, applyOpt); err != nil {
			log.Fatalf("apply failed: %v", err)
		}

		log.Println("applied")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	flags := applyCmd.Flags()

	flags.StringVarP(&applyOpt.Filepath, "filepath", "f", "", "filepath")
}
