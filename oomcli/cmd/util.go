package cmd

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/fatih/structtag"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

const (
	CSV        = "csv"
	ASCIITable = "ascii_table"
	Column     = "column"
	YAML       = "yaml"
)

const (
	MetadataFieldTruncateAt = 40
)

func mustOpenOomStore(ctx context.Context, opt types.OomStoreConfig) *oomstore.OomStore {
	store, err := oomstore.Open(ctx, oomStoreCfg)
	if err != nil {
		log.Fatalf("failed opening OomStore: %v", err)
	}
	return store
}

func stringPtr(s string) *string {
	return &s
}

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
		tableTag, err := tags.Get("table")
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

func serializeHeader(i interface{}, wide bool) ([]string, error) {
	tags, err := parseHeaderTag(i)
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
