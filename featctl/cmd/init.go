package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type initOption struct {
	types.OomStoreOpt
}

var initOpt initOption

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a feature store",
	Run: func(cmd *cobra.Command, args []string) {
		initOpt.OomStoreOpt = oomStoreOpt
		if _, err := oomstore.Create(context.Background(), initOpt.OomStoreOpt); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
