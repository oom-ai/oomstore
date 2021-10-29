package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type initOption struct {
	types.OomStoreOptV2
}

var initOpt initOption

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a feature store",
	Run: func(cmd *cobra.Command, args []string) {
		initOpt.OomStoreOptV2 = oomStoreOpt
		if _, err := oomstore.Create(context.Background(), initOpt.OomStoreOptV2); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
