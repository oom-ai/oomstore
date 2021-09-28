package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/_init"
	"github.com/spf13/cobra"
)

var initOpt _init.Option

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the feature store",
	Run: func(cmd *cobra.Command, args []string) {
		initOpt.DBOption = dbOption
		_init.Init(context.Background(), &initOpt)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
