package utils

import (
	database2 "github.com/onestore-ai/onestore/featctl/pkg/database"
	"github.com/onestore-ai/onestore/pkg/database"
)

func BuildSqlxDBOption(option database2.Option) *database.Option {
	return &database.Option{
		Host:   option.Host,
		Port:   option.Port,
		User:   option.User,
		Pass:   option.Pass,
		DbName: option.DbName,
	}
}
