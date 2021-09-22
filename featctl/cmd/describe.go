package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/describe"
	"github.com/spf13/cobra"
)

var describeOpt describe.Option

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "show details of a specific feature",
	Run: func(cmd *cobra.Command, args []string) {
		describeOpt.DBOption = dbOption
		describe.Run(context.Background(), &describeOpt)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)

	flags := describeCmd.Flags()

	flags.StringVarP(&describeOpt.Group, "group", "g", "", "feature group")
	_ = describeCmd.MarkFlagRequired("group")

	flags.StringVarP(&describeOpt.Name, "name", "n", "", "feature name")
	_ = describeCmd.MarkFlagRequired("name")
}
