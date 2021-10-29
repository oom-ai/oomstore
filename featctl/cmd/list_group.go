package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cobra"
)

type listFeatureGroupOption struct {
	EntityName *string
}

var listFeatureGroupOpt listFeatureGroupOption

var listFeatureGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "list feature groups",
	Example: `1. featctl list group
2. featctl list group --entity device
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureGroupOpt.EntityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		groups, err := oomStore.ListRichFeatureGroup(ctx, listFeatureGroupOpt.EntityName)
		if err != nil {
			log.Fatal(err)
		}
		if err := printFeatureGroups(groups); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureGroupCmd)

	flags := listFeatureGroupCmd.Flags()

	listFeatureGroupOpt.EntityName = flags.StringP("entity", "", "", "use to filter groups")
}

func printFeatureGroups(featureGroups []*types.RichFeatureGroup) error {
	w := csv.NewWriter(os.Stdout)

	if err := w.Write([]string{"Name", "Entity", "Description", "OnlineRevision", "OfflineLatestRevision", "OfflineLatestDataTable", "CreateTime", "ModifyTime"}); err != nil {
		return err
	}
	var onlineRevision, offlineRevision, offlineDataTable string
	for _, g := range featureGroups {
		onlineRevision = "<NULL>"
		offlineRevision = "<NULL>"
		offlineDataTable = "<NULL>"
		if g.OnlineRevision != nil {
			onlineRevision = fmt.Sprint(*g.OnlineRevision)
		}
		if g.OfflineRevision != nil {
			offlineRevision = fmt.Sprint(*g.OfflineRevision)
		}

		if g.OfflineDataTable != nil {
			offlineDataTable = *g.OfflineDataTable
		}

		if err := w.Write([]string{g.Name, g.EntityName, g.Description, onlineRevision, offlineRevision, offlineDataTable,
			g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
