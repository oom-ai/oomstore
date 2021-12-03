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

type HeaderTag struct {
	Value    interface{}
	Header   string
	Core     bool
	Truncate bool
}

type HeaderTagList []HeaderTag

func (l HeaderTagList) Core() HeaderTagList {
	var rs HeaderTagList
	for _, t := range l {
		if t.Core {
			rs = append(rs, t)
		}
	}
	return rs
}

func (l HeaderTagList) SerializeHeader(truncate bool) []string {
	var rs []string
	for _, t := range l {
		rs = append(rs, t.Header)
	}
	return rs
}

func (l HeaderTagList) SerializeRecord(truncate bool) ([]string, error) {
	var rs []string
	for _, t := range l {
		s, err := tableFormatSerialize(t.Value)
		if err != nil {
			return nil, err
		}
		if t.Truncate && len(s) > MetadataFieldTruncateAt {
			s = s[:MetadataFieldTruncateAt-3] + "..."
		}
		rs = append(rs, s)
	}
	return rs, nil
}

func parseHeaderTag(st interface{}) (HeaderTagList, error) {
	v := reflect.ValueOf(st)
	var rs HeaderTagList
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
		header := tableTag.Name
		core := tableTag.HasOption("core")
		truncate := tableTag.HasOption("truncate")
		rs = append(rs, HeaderTag{
			Value: v.Field(i).Interface(), Header: header, Core: core, Truncate: truncate})
	}
	return rs, nil
}

func tableFormatSerialize(i interface{}) (string, error) {
	switch v := i.(type) {
	case time.Time:
		return v.Format(time.RFC3339), nil
	default:
		return cast.ToStringE(v)
	}
}

func serializeHeader(e interface{}, wide bool) ([]string, error) {
	tags, err := parseHeaderTag(e)
	if err != nil {
		return nil, err
	}
	if wide {
		return tags.SerializeHeader(false), nil
	}
	return tags.Core().SerializeHeader(true), nil
}

func serializeRecord(i interface{}, wide bool) ([]string, error) {
	tags, err := parseHeaderTag(i)
	if err != nil {
		return nil, err
	}
	if wide {
		return tags.SerializeRecord(false)
	}
	return tags.Core().SerializeRecord(true)
}

// TODO: how to get the inner element type of the slice ? to that we don't need the
// parameter `e`
func serializeList(s interface{}, e interface{}, output string, wide bool) error {
	slice, ok := s.([]*interface{})
	if !ok {
		return fmt.Errorf("expect slice, got %T", slice)
	}
	switch output {
	case CSV:
		return serializeSliceInCSV(slice, e, wide)
	case ASCIITable:
		return serializeSliceInASCIITable(slice, e, true, wide)
	case Column:
		return serializeSliceInASCIITable(slice, e, false, wide)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func serializeSliceInCSV(slice []*interface{}, e interface{}, wide bool) error {
	header, err := serializeHeader(e, wide)
	if err != nil {
		return err
	}
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(header); err != nil {
		return err
	}
	for _, elem := range slice {
		record, err := serializeRecord(*elem, wide)
		if err != nil {
			return err
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func serializeSliceInASCIITable(slice []*interface{}, e interface{}, border, wide bool) error {
	header, err := serializeHeader(e, wide)
	if err != nil {
		return err
	}
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

	for _, elem := range slice {
		record, err := serializeRecord(*elem, wide)
		if err != nil {
			return err
		}
		table.Append(record)
	}
	table.Render()
	return nil
}
