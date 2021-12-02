package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type getMetaGroupOption struct {
	entityName *string
}

var getMetaGroupOpt getMetaGroupOption

var getMetaGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "get existing group given specific conditions",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			getMetaGroupOpt.entityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatalf("argument at most one, got %d", len(args))
		}

		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		var entityID *int

		if getMetaGroupOpt.entityName != nil {
			entity, err := oomStore.GetEntityByName(ctx, *getMetaGroupOpt.entityName)
			if err != nil {
				log.Fatalf("failed to get entity name='%s': %v", *getMetaGroupOpt.entityName, err)
			}
			entityID = &entity.ID
		}

		groups, err := oomStore.ListGroup(ctx, entityID)
		if err != nil {
			log.Fatalf("failed getting feature groups, error %v\n", err)
		}

		if len(args) > 0 {
			groups = groups.Filter(func(g *types.Group) bool {
				return g.Name == args[0]
			})
		}

		if err := printGroups(ctx, oomStore, groups, *getMetaOutput, *getMetaWide); err != nil {
			log.Fatalf("failed printing feature groups, error %v\n", err)
		}
	},
}

func init() {
	getMetaCmd.AddCommand(getMetaGroupCmd)

	flags := getMetaGroupCmd.Flags()

	getMetaGroupOpt.entityName = flags.StringP("entity", "", "", "use to filter groups")
}

func printGroups(ctx context.Context, oomStore *oomstore.OomStore, groups []*types.Group, output string, wide bool) error {
	switch output {
	case CSV:
		return printGroupsInCSV(groups, wide)
	case ASCIITable:
		return printGroupsInASCIITable(groups, true, wide)
	case Column:
		return printGroupsInASCIITable(groups, false, wide)
	case YAML:
		return printGroupInYaml(ctx, oomStore, groups)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printGroupInYaml(ctx context.Context, oomStore *oomstore.OomStore, groups types.GroupList) error {
	var (
		out      []byte
		err      error
		features types.FeatureList
		items    apply.GroupItems
	)

	// TODO: Use group ids to filter, rather than taking them all out
	if features, err = oomStore.ListFeature(ctx, types.ListFeatureOpt{}); err != nil {
		return err
	}

	items = apply.FromGroupList(groups, features)

	if len(items.Items) > 1 {
		if out, err = yaml.Marshal(items); err != nil {
			return err
		}
	} else if len(items.Items) == 1 {
		if out, err = yaml.Marshal(items.Items[0]); err != nil {
			return err
		}
	}
	fmt.Println(strings.Trim(string(out), "\n"))
	return nil
}

func printGroupsInCSV(groups types.GroupList, wide bool) error {
	w := csv.NewWriter(os.Stdout)

	if err := w.Write(groupHeader(wide)); err != nil {
		return err
	}
	for _, g := range groups {
		if err := w.Write(groupRecord(g, wide)); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func printGroupsInASCIITable(groups types.GroupList, border, wide bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(groupHeader(wide))
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if !border {
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetNoWhiteSpace(true)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetTablePadding("  ")
	}

	for _, group := range groups {
		table.Append(groupRecord(group, wide))
	}
	table.Render()
	return nil
}

func groupHeader(wide bool) []string {
	if wide {
		return []string{"ID", "NAME", "ENTITY", "DESCRIPTION", "ONLINE-REVISION-ID", "CREATE-TIME", "MODIFY-TIME"}
	}
	return []string{"ID", "NAME", "ENTITY", "DESCRIPTION"}
}

func groupRecord(g *types.Group, wide bool) []string {
	onlineRevisionID := "<NULL>"
	if g.OnlineRevisionID != nil {
		onlineRevisionID = fmt.Sprint(*g.OnlineRevisionID)
	}
	if wide {
		return []string{strconv.Itoa(g.ID), g.Name, g.Entity.Name, g.Description, onlineRevisionID,
			g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}
	}
	desc := g.Description
	if len(desc) > MaxDescriptionLen {
		desc = fmt.Sprintf("%s...", desc[0:MaxDescriptionLen])
	}
	return []string{strconv.Itoa(g.ID), g.Name, g.Entity.Name, desc}
}
