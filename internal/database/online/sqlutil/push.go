package sqlutil

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type PushCondition struct {
	Inserts            string
	InsertPlaceholders string
	InsertValues       []interface{}
	UpdateValues       []interface{}
	UpdatePlaceholders string
}

func BuildPushCondition(opt online.PushOpt, backend types.BackendType) *PushCondition {
	qt := dbutil.QuoteFn(backend)
	cond := PushCondition{}

	cond.Inserts = qt(append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)...)
	cond.InsertValues = append([]interface{}{opt.EntityKey}, opt.FeatureValues...)
	cond.InsertPlaceholders = dbutil.Fill(len(cond.InsertValues), "?", ",")

	updatePlaceholders := make([]string, 0, opt.FeatureList.Len())
	for _, name := range opt.FeatureList.Names() {
		updatePlaceholders = append(updatePlaceholders, fmt.Sprintf("%s=?", qt(name)))
	}
	cond.UpdatePlaceholders = strings.Join(updatePlaceholders, ",")
	cond.UpdateValues = opt.FeatureValues

	return &cond
}
