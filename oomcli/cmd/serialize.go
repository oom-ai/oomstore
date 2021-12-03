package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
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

func serializeMetadata(w io.Writer, i interface{}, format string, wide bool) error {
	truncate := !wide
	headerTokens, dataTokens, err := parseTokenLists(i)
	if err != nil {
		return err
	}
	header := parseHeader(headerTokens, wide, truncate)
	records, err := parseRecords(dataTokens, wide, truncate)
	if err != nil {
		return err
	}
	switch format {
	case CSV:
		return serializeInCSV(w, header, records)
	case ASCIITable:
		return serializeInASCIITable(w, header, records, true)
	case Column:
		return serializeInASCIITable(w, header, records, false)
	default:
		return fmt.Errorf("unsupported output format %s", format)
	}
}

func serializeInCSV(w io.Writer, header []string, records [][]string) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(header); err != nil {
		return err
	}
	if err := cw.WriteAll(records); err != nil {
		return err
	}
	cw.Flush()
	return nil
}

func serializeInASCIITable(w io.Writer, header []string, records [][]string, border bool) error {
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

func parseHeader(tokens TokenList, wide, truncate bool) []string {
	if !wide {
		return tokens.Brief().SerializeHeader(truncate)
	}
	return tokens.SerializeHeader(truncate)
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
