package cmd

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a feature store",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := oomstore.Create(context.Background(), oomStoreCfg); err != nil {
			exit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
