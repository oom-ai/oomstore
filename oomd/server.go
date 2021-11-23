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
			Status: &status.Status{
				Code:    int32(code.Code_INTERNAL),
				Message: fmt.Sprintf("failed at OnlineGet, err=%+v", err),
			},
		}, fmt.Errorf("failed at OnlineGet, err=%+v", err)
	}

	resultMap := make(map[string]*anypb.Any)
	for key, value := range result.FeatureValueMap {
		bytes, err := json.Marshal(value)
		if err != nil {
			return &codegen.OnlineGetResponse{
				Status: &status.Status{
					Code:    int32(code.Code_INTERNAL),
					Message: fmt.Sprintf("failed to marshal %v", value),
				},
			}, fmt.Errorf("failed to marshal %v", value)
		}
		resultMap[key] = &anypb.Any{
			Value: bytes,
		}
	}
	return &codegen.OnlineGetResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
		Result: &codegen.FeatureValueMap{
			Map: resultMap,
		},
	}, nil
}

func (s *server) OnlineMultiGet(ctx context.Context, req *codegen.OnlineMultiGetRequest) (*codegen.OnlineMultiGetResponse, error) {
	panic("implement me")
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
