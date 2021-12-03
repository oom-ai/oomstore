package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
)

func serializeValue(i interface{}) string {
	if reflect.TypeOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil() {
		return "<NULL>"
	}
	switch v := i.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return cast.ToString(v)
	}
}

func serializeMetadataList(i interface{}, output string, wide bool) error {
	truncate := !wide
	lists, err := parseTokenLists(i)
	if err != nil {
		return err
	}
	header := parseHeader(lists, wide, truncate)
	records, err := parseRecords(lists, wide, truncate)
	if err != nil {
		return err
	}
	switch output {
	case CSV:
		return serializeInCSV(header, records)
	case ASCIITable:
		return serializeInASCIITable(header, records, true)
	case Column:
		return serializeInASCIITable(header, records, false)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func serializeInCSV(header []string, records [][]string) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(header); err != nil {
		return err
	}
	if err := w.WriteAll(records); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func serializeInASCIITable(header []string, records [][]string, border bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
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

	table.AppendBulk(records)
	table.Render()
	return nil
}

func parseHeader(tokens []TokenList, wide, truncate bool) []string {
	t := tokens[0]
	if !wide {
		t = t.Brief()
	}
	return t.SerializeHeader(truncate)
}

func parseRecords(tokens []TokenList, wide, truncate bool) ([][]string, error) {
	var rs [][]string
	for _, t := range tokens {
		if !wide {
			t = t.Brief()
		}
		record, err := t.SerializeRecord(truncate)
		if err != nil {
			return nil, err
		}
		rs = append(rs, record)
	}
	return rs, nil
}
