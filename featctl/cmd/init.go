package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/onestore-ai/onestore/pkg/onestore"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type initOption struct {
	types.OneStoreOpt
}

var initOpt initOption

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the feature store",
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
