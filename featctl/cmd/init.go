package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/_init"
	"github.com/spf13/cobra"
)

var initOpt _init.Option

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the feature store",
	Run: func(cmd *cobra.Command, args []string) {
		initOpt.DBOption = sqlxDbOption
		_init.Init(context.Background(), &initOpt)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
