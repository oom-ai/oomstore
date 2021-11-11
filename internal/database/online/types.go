package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type GetOpt struct {
	Entity      *typesv2.Entity
	RevisionID  int32
	EntityKey   string
	FeatureList typesv2.FeatureList
}

type MultiGetOpt struct {
	Entity      *typesv2.Entity
	RevisionID  int32
	EntityKeys  []string
	FeatureList typesv2.FeatureList
}

type ImportOpt struct {
	Revision    *typesv2.Revision
	Entity      *typesv2.Entity
	Stream      <-chan *typesv2.RawFeatureValueRecord
	FeatureList typesv2.FeatureList
}
