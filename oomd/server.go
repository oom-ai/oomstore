package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	code "google.golang.org/genproto/googleapis/rpc/code"
	status "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/oom-ai/oomstore/oomd/codegen"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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

	valueMap, err := convertToValueMap(result.FeatureValueMap)
	if err != nil {
		return &codegen.OnlineGetResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}
	return &codegen.OnlineGetResponse{
		Status: buildStatus(code.Code_OK, ""),
		Result: &codegen.FeatureValueMap{
			Map: valueMap,
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
		valueMap, err := convertToValueMap(featureValues.FeatureValueMap)
		if err != nil {
			return &codegen.OnlineMultiGetResponse{
				Status: buildStatus(code.Code_INTERNAL, err.Error()),
			}, err
		}
		resultMap[entityKey] = &codegen.FeatureValueMap{
			Map: valueMap,
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

func (s *server) Import(stream codegen.OomD_ImportServer) error {
	firstReq, err := stream.Recv()
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()

	go func() {
		defer func() {
			_ = writer.Close()
		}()

		if _, err := writer.Write(firstReq.Row); err != nil {
			return
		}

		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			if _, err := writer.Write(req.Row); err != nil {
				return
			}
		}
	}()

	revisionID, err := s.oomstore.Import(context.Background(), types.ImportOpt{
		GroupName:   firstReq.GroupName,
		Revision:    firstReq.Revision,
		Description: firstReq.Description,
		DataSource: types.CsvDataSource{
			Reader:    reader,
			Delimiter: ",",
		},
	})
	if err != nil {
		return stream.SendAndClose(&codegen.ImportResponse{
			Status: &status.Status{
				Code:    int32(code.Code_INTERNAL),
				Message: err.Error(),
			},
		})
	}
	return stream.SendAndClose(&codegen.ImportResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
		RevisionId: int64(revisionID),
	})
}

func (s *server) Join(stream codegen.OomD_JoinServer) error {
	panic("implement me")
}

func (s *server) ImportByFile(ctx context.Context, req *codegen.ImportByFileRequest) (*codegen.ImportResponse, error) {
	revisionID, err := s.oomstore.ImportByFile(ctx, types.ImportByFileOpt{
		GroupName:   req.GroupName,
		Description: req.Description,
		Revision:    req.Revision,
		DataSource: types.CsvDataSourceWithFile{
			InputFilePath: req.InputFilePath,
			Delimiter:     ",",
		},
	})
	if err != nil {
		return &codegen.ImportResponse{
			Status:     buildStatus(code.Code_INTERNAL, err.Error()),
			RevisionId: int64(revisionID),
		}, err
	}

	return &codegen.ImportResponse{
		Status:     buildStatus(code.Code_OK, ""),
		RevisionId: int64(revisionID),
	}, nil
}

func (s *server) JoinByFile(ctx context.Context, req *codegen.JoinByFileRequest) (*codegen.JoinByFileResponse, error) {
	err := s.oomstore.JoinByFile(ctx, types.JoinByFileOpt{
		FeatureNames:   req.FeatureNames,
		InputFilePath:  req.InputFilePath,
		OutputFilePath: req.OutputFilePath,
	})
	if err != nil {
		return &codegen.JoinByFileResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}

	return &codegen.JoinByFileResponse{
		Status: buildStatus(code.Code_OK, ""),
	}, nil
}

func convertToValueMap(m map[string]interface{}) (map[string]*codegen.Value, error) {
	valueMap := make(map[string]*codegen.Value)
	for key, i := range m {
		value, err := convertInterfaceToValue(i)
		if err != nil {
			return nil, err
		}
		valueMap[key] = value
	}
	return valueMap, nil
}

func buildStatus(code code.Code, message string) *status.Status {
	return &status.Status{
		Code:    int32(code),
		Message: message,
	}
}

func convertInterfaceToValue(i interface{}) (*codegen.Value, error) {
	switch s := i.(type) {
	case int64:
		return &codegen.Value{
			Kind: &codegen.Value_Int64Value{
				Int64Value: s,
			},
		}, nil
	case float64:
		return &codegen.Value{
			Kind: &codegen.Value_DoubleValue{
				DoubleValue: s,
			},
		}, nil
	case string:
		return &codegen.Value{
			Kind: &codegen.Value_StringValue{
				StringValue: s,
			},
		}, nil
	case bool:
		return &codegen.Value{
			Kind: &codegen.Value_BoolValue{
				BoolValue: s,
			},
		}, nil
	case time.Time:
		return &codegen.Value{
			Kind: &codegen.Value_UnixTimestampValue{
				UnixTimestampValue: s.UnixMilli(),
			},
		}, nil
	case []byte:
		return &codegen.Value{
			Kind: &codegen.Value_BytesValue{
				BytesValue: s,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported value type %v", i)
	}
}
