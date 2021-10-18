package cmd

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var listEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "list all existing entities",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		entities, err := oneStore.ListEntity(ctx)
		if err != nil {
			log.Fatalf("failed listing entities, error %v\n", err)
		}

		// print csv to stdout
		fmt.Println(entityCsvHeader())
		for _, entity := range entities {
			recordStr, err := entityCsvRecord(entity)
			if err != nil {
				log.Fatalf("failed writing entity %s, error %v\n", entity.Name, err)
			}
			fmt.Println(recordStr)
		}
	},
}

func init() {
	listCmd.AddCommand(listEntityCmd)
}

func entityCsvHeader() string {
	return strings.Join([]string{"Name", "Length", "Description", "CreateTime", "ModifyTime"}, ",")
}

func entityCsvRecord(entity *types.Entity) (string, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	w := csv.NewWriter(buf)
	if err := w.Write([]string{entity.Name, strconv.Itoa(entity.Length), entity.Description, entity.CreateTime.Format(time.RFC3339), entity.ModifyTime.Format(time.RFC3339)}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
