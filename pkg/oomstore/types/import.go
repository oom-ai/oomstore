package types

import "io"

type DataSourceType int

const (
	CSV_FILE DataSourceType = iota
	CSV_READER
	TABLE_LINK
)

type ImportOpt struct {
	GroupName   string
	Description string
	Revision    *int64

	DataSourceType      DataSourceType
	CsvFileDataSource   *CsvFileDataSource
	CsvReaderDataSource *CsvReaderDataSource
	TableLinkDataSource *TableLinkDataSource
}

type CsvReaderDataSource struct {
	Reader    io.Reader
	Delimiter string
}

type CsvFileDataSource struct {
	InputFilePath string
	Delimiter     string
}

type TableLinkDataSource struct {
	TableName string
}

type ImportStreamOpt struct {
	GroupName   string
	Description string

	DataSourceType      DataSourceType
	CsvFileDataSource   *CsvFileDataSource
	CsvReaderDataSource *CsvReaderDataSource
	TableLinkDataSource *TableLinkDataSource
}
