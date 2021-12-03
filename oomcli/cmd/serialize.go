package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/fatih/structtag"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
)

type Token struct {
	Value    interface{}
	Name     string
	Wide     bool
	Truncate bool
}

type TokenList []Token

func parseTokens(i interface{}) (TokenList, error) {
	v := reflect.ValueOf(i)
	var rs TokenList
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag
		tags, err := structtag.Parse(string(tag))
		if err != nil {
			return rs, err
		}
		tableTag, err := tags.Get("oomcli")
		if err != nil {
			return rs, err
		}
		name := tableTag.Name
		wide := tableTag.HasOption("wide")
		truncate := tableTag.HasOption("truncate")
		rs = append(rs, Token{
			Value: v.Field(i).Interface(), Name: name, Wide: wide, Truncate: truncate})
	}
	return rs, nil
}

func (l TokenList) Brief() TokenList {
	var rs TokenList
	for _, t := range l {
		if !t.Wide {
			rs = append(rs, t)
		}
	}
	return rs
}

func (l TokenList) SerializeHeader(truncate bool) []string {
	var rs []string
	for _, t := range l {
		if t.Truncate && len(t.Name) > MetadataFieldTruncateAt {
			rs = append(rs, t.Name[:MetadataFieldTruncateAt-3]+"...")
		}
		rs = append(rs, t.Name)
	}
	return rs
}

func (l TokenList) SerializeRecord(truncate bool) ([]string, error) {
	var rs []string
	for _, t := range l {
		s := serializeValue(t.Value)
		if t.Truncate && len(s) > MetadataFieldTruncateAt {
			s = s[:MetadataFieldTruncateAt-3] + "..."
		}
		rs = append(rs, s)
	}
	return rs, nil
}

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

func serializeMetadata(i interface{}, output string, wide bool) error {
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
