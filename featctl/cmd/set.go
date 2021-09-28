package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/featctl/pkg/set"
	"github.com/spf13/cobra"
)

var setOpt set.Option

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "configure a specified feature",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		setOpt.RevisionChanged = flags.Changed("revision")
		setOpt.DescriptionChanged = flags.Changed("description")
		setOpt.StatusChanged = flags.Changed("status")
		setOpt.RevisionsLimitChanged = flags.Changed("revisions-limit")
		setOpt.DBOption = dbOption

		if setOpt.StatusChanged {
			if err := validateStatus(setOpt.Status); err != nil {
				return err
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		set.Set(ctx, &setOpt)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	flags := setCmd.Flags()

	flags.StringVarP(&setOpt.Group, "group", "g", "", "feature group")
	_ = setCmd.MarkFlagRequired("group")

	flags.StringVarP(&setOpt.Name, "name", "n", "", "feature name")
	_ = setCmd.MarkFlagRequired("name")

	flags.StringVar(&setOpt.Revision, "revision", "", "current revision")
	flags.StringVar(&setOpt.Description, "description", "", "feature description")
	flags.StringVarP(&setOpt.Status, "status", "s", "", "feature status")
	flags.IntVar(&setOpt.RevisionsLimit, "revisions-limit", 0, "feature history revisions limit")
}
