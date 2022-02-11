package sqlutil

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type PushQueryParams struct {
	TableName          string
	EntityName         string
	Fields             string
	InsertPlaceholders string
	UpdatePlaceholders string
	InsertValues       []interface{}
	UpdateValues       []interface{}
	Backend            types.BackendType
}

func BuildPushQueryParams(opt online.PushOpt, backend types.BackendType) PushQueryParams {
	qt := dbutil.QuoteFn(backend)
	params := PushQueryParams{
		TableName:    dbutil.OnlineStreamTableName(opt.GroupID),
		EntityName:   opt.EntityName,
		Fields:       qt(append([]string{opt.EntityName}, opt.Features.Names()...)...),
		InsertValues: append([]interface{}{opt.EntityKey}, opt.FeatureValues...),
		UpdateValues: opt.FeatureValues,
		Backend:      backend,
	}

	params.InsertPlaceholders = dbutil.Fill(len(params.InsertValues), "?", ",")
	updatePlaceholders := make([]string, 0, opt.Features.Len())
	for _, name := range opt.Features.Names() {
		updatePlaceholders = append(updatePlaceholders, fmt.Sprintf("%s=?", qt(name)))
	}
	params.UpdatePlaceholders = strings.Join(updatePlaceholders, ",")
	return params
}

func BuildPushQuery(params PushQueryParams, queryTemplate string) (string, error) {
	qt := dbutil.QuoteFn(params.Backend)
	t := template.Must(template.New("push").Funcs(template.FuncMap{
		"qt": qt,
	}).Parse(queryTemplate))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", errdefs.WithStack(err)
	}
	return buf.String(), nil
}
