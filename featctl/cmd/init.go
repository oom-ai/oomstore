package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/onestore"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

type initOption struct {
	types.OneStoreOpt
}

var initOpt initOption

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a feature store",
	Run: func(cmd *cobra.Command, args []string) {
		initOpt.OneStoreOpt = oneStoreOpt
		if _, err := onestore.Create(context.Background(), initOpt.OneStoreOpt); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
