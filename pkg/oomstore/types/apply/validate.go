package apply

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/qri-io/jsonschema"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

var (
	entityValidator  jsonschema.Schema
	groupValidator   jsonschema.Schema
	featureValidator jsonschema.Schema

	//go:embed entity.json
	entitySchema string
	//go:embed group.json
	groupSchema string
	//go:embed feature.json
	featureSchema string
)

func init() {
	if err := json.Unmarshal([]byte(entitySchema), &entityValidator); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(groupSchema), &groupValidator); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(featureSchema), &featureValidator); err != nil {
		panic(err)
	}
}

func (e *Entity) Validate() error {
	data, err := json.Marshal(e)
	if err != nil {
		return errdefs.WithStack(err)
	}

	return validate(fmt.Sprintf("Entity %s", e.Name), data, &entityValidator)
}

func (g *Group) Validate() error {
	data, err := json.Marshal(g)
	if err != nil {
		return errdefs.WithStack(err)
	}

	return validate(fmt.Sprintf("Group %s", g.Name), data, &groupValidator)
}

func (f *Feature) Validate() error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}

	return validate(fmt.Sprintf("Feature %s", f.Name), data, &featureValidator)
}

func validate(errPrefix string, data []byte, schema *jsonschema.Schema) error {
	errs, err := schema.ValidateBytes(context.Background(), data)
	if err != nil {
		return errdefs.WithStack(err)
	}

	if len(errs) > 0 {
		err := fmt.Errorf("%s: %v", errPrefix, errs[0])
		for i := 1; i < len(errs); i++ {
			err = fmt.Errorf("%v\n%v", err, errs[i])
		}
		return errdefs.WithStack(err)
	}
	return nil
}
