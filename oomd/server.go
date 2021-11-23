package oomd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oom-ai/oomstore/oomd/codegen"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	code "google.golang.org/genproto/googleapis/rpc/code"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type server struct {
	codegen.UnimplementedOomDServer

	oomstore *oomstore.OomStore
}

func (s *server) OnlineGet(ctx context.Context, req *codegen.OnlineGetRequest) (*codegen.OnlineGetResponse, error) {
	result, err := s.oomstore.OnlineGet(ctx, types.OnlineGetOpt{
		FeatureNames: req.FeatureNames,
		EntityKey:    req.EntityKey,
	})
	if err != nil {
		return &codegen.OnlineGetResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}

	anyMap, err := convertFeatureValueMap(result.FeatureValueMap)
	if err != nil {
		return &codegen.OnlineGetResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}
	return &codegen.OnlineGetResponse{
		Status: buildStatus(code.Code_OK, ""),
		Result: &codegen.FeatureValueMap{
			Map: anyMap,
		},
	}, nil
}

func (s *server) OnlineMultiGet(ctx context.Context, req *codegen.OnlineMultiGetRequest) (*codegen.OnlineMultiGetResponse, error) {
	result, err := s.oomstore.OnlineMultiGet(ctx, types.OnlineMultiGetOpt{
		FeatureNames: req.FeatureNames,
		EntityKeys:   req.EntityKeys,
	})
	if err != nil {
		return &codegen.OnlineMultiGetResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}

	resultMap := make(map[string]*codegen.FeatureValueMap)
	for entityKey, featureValues := range result {
		anyMap, err := convertFeatureValueMap(featureValues.FeatureValueMap)
		if err != nil {
			return &codegen.OnlineMultiGetResponse{
				Status: buildStatus(code.Code_INTERNAL, err.Error()),
			}, err
		}
		resultMap[entityKey] = &codegen.FeatureValueMap{
			Map: anyMap,
		}
	}
	return &codegen.OnlineMultiGetResponse{
		Status: buildStatus(code.Code_OK, ""),
		Result: resultMap,
	}, nil
}

func (s *server) Sync(ctx context.Context, req *codegen.SyncRequest) (*codegen.SyncResponse, error) {
	if err := s.oomstore.Sync(ctx, types.SyncOpt{RevisionID: int(req.RevisionId)}); err != nil {
		return &codegen.SyncResponse{
			Status: &status.Status{
				Code:    int32(code.Code_INTERNAL),
				Message: err.Error(),
			},
		}, err
	}

	return &codegen.SyncResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

func (s *server) Import(codegen.OomD_ImportServer) error {
	panic("implement me")
}

func (s *server) Join(*codegen.JoinRequest, codegen.OomD_JoinServer) error {
	panic("implement me")
}

func (s *server) ImportByFile(context.Context, *codegen.ImportByFileRequest) (*codegen.ImportResponse, error) {
	panic("implement me")
}

func (s *server) JoinByFile(context.Context, *codegen.JoinRequest) (*codegen.JoinByFileResponse, error) {
	panic("implement me")
}

func convertFeatureValueMap(m map[string]interface{}) (map[string]*anypb.Any, error) {
	anyMap := make(map[string]*anypb.Any)
	for key, value := range m {
		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %v", value)
		}
		anyMap[key] = &anypb.Any{
			Value: bytes,
		}
	}
	return anyMap, nil
}

func buildStatus(code code.Code, message string) *status.Status {
	return &status.Status{
		Code:    int32(code),
		Message: message,
	}
}
