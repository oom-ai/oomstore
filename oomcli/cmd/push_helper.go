package cmd

import (
	"strconv"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func parsePushFeatureArgument(args map[string]string, features types.FeatureList) (map[string]interface{}, error) {
	featureValues := make(map[string]interface{})
	for k, v := range args {
		feature := features.Find(func(f *types.Feature) bool {
			return f.Name == k
		})
		if feature != nil {
			value, err := serializerFeatureValue(k, v, feature.ValueType)
			if err != nil {
				return nil, err
			}
			featureValues[k] = value
		}
	}
	return featureValues, nil
}

func serializerFeatureValue(featureName, featureValue string, valueType types.ValueType) (interface{}, error) {
	switch valueType {
	case types.String:
		return featureValue, nil
	case types.Int64:
		ret, err := strconv.Atoi(featureValue)
		if err != nil {
			return nil, errdefs.Errorf("feature %s is of type int64 and cannot be assigned to %s", featureName, featureValue)
		}
		return int64(ret), nil
	case types.Float64:
		ret, err := strconv.ParseFloat(featureValue, 64)
		return ret, errdefs.WithStack(err)
	case types.Bool:
		if featureValue == "true" || featureValue == "1" {
			return true, nil
		}
		if featureValue == "false" || featureValue == "0" {
			return false, nil
		}
		return nil, errdefs.Errorf("feature %s is of type bool and cannot be assigned to %s", featureName, featureValue)
	case types.Time:
		intValue, err := strconv.Atoi(featureValue)
		if err != nil {
			return nil, errdefs.WithStack(err)
		}
		return strconv.FormatInt(int64(intValue), 36), nil
	case types.Bytes:
		return []byte(featureValue), nil
	default:
		return nil, errdefs.Errorf("unable to serialize feature %s, type %s", featureName, featureValue)
	}
}
