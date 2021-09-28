package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/list_revision"
	"github.com/spf13/cobra"
)

var listRevisionOpt list_revision.Option

var listRevisionCmd = &cobra.Command{
	Use:   "revision",
	Short: "list historical revisions given a specific group",
	Run: func(cmd *cobra.Command, args []string) {
		listRevisionOpt.DBOption = dbOption
		list_revision.ListRevision(context.Background(), &listRevisionOpt)
	},
}

func init() {
	listCmd.AddCommand(listRevisionCmd)

	flags := listRevisionCmd.Flags()

	flags.StringVarP(&listRevisionOpt.Group, "group", "g", "", "feature group")
	_ = listRevisionCmd.MarkFlagRequired("group")
}
