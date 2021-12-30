package sqlutil

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func BuildPushCondition(opt online.PushOpt, backend types.BackendType) (string, string, string, []interface{}, error) {
	qt, err := dbutil.QuoteFn(backend)
	if err != nil {
		return "", "", "", nil, err
	}

	insertColumns := append([]string{opt.Entity.Name}, opt.FeatureNames...)
	insertValues := append([]interface{}{opt.EntityKey}, opt.FeatureValues...)

	updatePlaceholders := make([]string, 0, len(opt.FeatureNames))
	for _, name := range opt.FeatureNames {
		updatePlaceholders = append(updatePlaceholders, fmt.Sprintf("%s=?", qt(name)))
	}

	return qt(insertColumns...),
		dbutil.Fill(len(insertColumns), "?", ","),
		strings.Join(updatePlaceholders, ","),
		append(insertValues, opt.FeatureValues...),
		nil

}
