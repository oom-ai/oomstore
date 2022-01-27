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
	entityValidate  jsonschema.Schema
	groupValidate   jsonschema.Schema
	featureValidate jsonschema.Schema

	//go:embed entity.json
	entitySchema string
	//go:embed group.json
	groupSchema string
	//go:embed feature.json
	featureSchema string
)

func init() {
	if err := json.Unmarshal([]byte(entitySchema), &entityValidate); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(groupSchema), &groupValidate); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(featureSchema), &featureValidate); err != nil {
		panic(err)
	}
}

func (e *Entity) Validate() error {
	data, err := json.Marshal(e)
	if err != nil {
		return errdefs.WithStack(err)
	}

	errs, err := entityValidate.ValidateBytes(context.Background(), data)
	if err != nil {
		return errdefs.WithStack(err)
	}

	if len(errs) > 0 {
		err = fmt.Errorf("Entity Validate: %v", errs[0])
		for i := 1; i < len(errs); i++ {
			err = fmt.Errorf("%v\n%v", err, errs[i])
		}
		return errdefs.WithStack(err)
	}
	return nil

}

func (g *Group) Validate() error {
	data, err := json.Marshal(g)
	if err != nil {
		return errdefs.WithStack(err)
	}

	errs, err := groupValidate.ValidateBytes(context.Background(), data)
	if err != nil {
		return errdefs.WithStack(err)
	}

	if len(errs) > 0 {
		err = fmt.Errorf("Group Validate: %v", errs[0])
		for i := 1; i < len(errs); i++ {
			err = fmt.Errorf("%v\n%v", err, errs[i])
		}
		return errdefs.WithStack(err)
	}
	return nil
}

func (f *Feature) Validate() error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}

	errs, err := featureValidate.ValidateBytes(context.Background(), data)
	if err != nil {
		return errdefs.WithStack(err)
	}

	if len(errs) > 0 {
		err = fmt.Errorf("Feature Validate: %v", errs[0])
		for i := 1; i < len(errs); i++ {
			err = fmt.Errorf("%v\n%v", err, errs[i])
		}
		return errdefs.WithStack(err)
	}
	return nil
}
