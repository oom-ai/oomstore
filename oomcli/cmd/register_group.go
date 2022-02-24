package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type RegisterGroupOption struct {
	types.CreateGroupOpt

	snapshotInterval string
}

var registerGroupOpt RegisterGroupOption
var registerGroupCmd = &cobra.Command{
	Use:   "group <group_name>",
	Short: "Register a new feature group",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if registerGroupOpt.Category != types.CategoryBatch && registerGroupOpt.Category != types.CategoryStream {
			exitf("illegal category '%s', should be either 'stream' or 'batch'", registerGroupOpt.Category)
		}

		if registerGroupOpt.Category == types.CategoryStream {
			t, err := time.ParseDuration(registerGroupOpt.snapshotInterval)
			if err != nil {
				exitf("invalid snapshot_interval %s: %v", registerGroupOpt.snapshotInterval, err)
			}
			registerGroupOpt.CreateGroupOpt.SnapshotInterval = int(t.Seconds())
		}

		registerGroupOpt.GroupName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if _, err := oomStore.CreateGroup(ctx, registerGroupOpt.CreateGroupOpt); err != nil {
			exitf("failed registering new group: %+v\n", err)
		}
	},
}

func init() {
	registerCmd.AddCommand(registerGroupCmd)

	flags := registerGroupCmd.Flags()

	flags.StringVarP(&registerGroupOpt.EntityName, "entity", "e", "", "entity name")
	_ = registerGroupCmd.MarkFlagRequired("entity")

	flags.StringVarP(&registerGroupOpt.Category, "category", "c", "batch", "group category")

	flags.StringVarP(&registerGroupOpt.snapshotInterval, "snapshot_interval", "s", "24h", "stream group snapshot interval")

	flags.StringVarP(&registerGroupOpt.Description, "description", "d", "", "group description")
}
