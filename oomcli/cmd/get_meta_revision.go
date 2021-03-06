package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type getRevisionOption struct {
	groupName *string
}

var getRevisionOpt getRevisionOption

var getMetaRevisionCmd = &cobra.Command{
	Use:   "revision [revision_id]",
	Short: "Get existing revision given specific conditions",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("group") {
			getRevisionOpt.groupName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var groupID *int

		if getRevisionOpt.groupName != nil {
			group, err := oomStore.GetGroupByName(ctx, *getRevisionOpt.groupName)
			if err != nil {
				exitf("failed to get feature group name=%s: %+v", *getRevisionOpt.groupName, err)
			}
			groupID = &group.ID
		}

		revisions, err := oomStore.ListRevision(ctx, groupID)
		if err != nil {
			exitf("failed getting revisions, error %+v\n", err)
		}

		if len(args) > 0 {
			if revisions = revisions.Filter(func(r *types.Revision) bool {
				return strconv.Itoa(r.ID) == args[0]
			}); len(revisions) == 0 {
				exitf("revision '%s' not found", args[0])
			}
		}

		if err := serializeMetadata(os.Stdout, revisions, *getMetaOutput, *getMetaWide); err != nil {
			exitf("failed printing revisions, error %+v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaRevisionCmd)

	flags := getMetaRevisionCmd.Flags()
	getRevisionOpt.groupName = flags.StringP("group", "g", "", "feature group")
}
