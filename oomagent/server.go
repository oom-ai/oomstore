package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/oom-ai/oomstore/oomagent/codegen"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	code "google.golang.org/genproto/googleapis/rpc/code"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

type server struct {
	codegen.UnimplementedOomAgentServer

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
	if err := s.oomstore.Sync(ctx, types.SyncOpt{
		RevisionID: int(req.RevisionId),
		PurgeDelay: int(req.PurgeDelay),
	}); err != nil {
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

func (s *server) ChannelImport(stream codegen.OomAgent_ChannelImportServer) error {
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

	revisionID, err := s.oomstore.ChannelImport(context.Background(), types.ChannelImport{
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

func (s *server) Import(ctx context.Context, req *codegen.ImportRequest) (*codegen.ImportResponse, error) {
	revisionID, err := s.oomstore.Import(ctx, types.ImportOpt{
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

func (s *server) ChannelJoin(stream codegen.OomAgent_ChannelJoinServer) error {
	// We need to read the first request to get the feature names and value names
	firstReq, err := stream.Recv()
	if err != nil {
		return err
	}

	// A global error
	var globalErr error

	// This channel indicates when the the ChannelJoin oomstore operation is finished, whether succeeded or failed.
	done := make(chan struct{})
	// This channel receives requests from the client.
	entityRows := make(chan types.EntityRow)

	// This goroutine runs the join operation, and send whatever joined as the response
	go func() {
		joinResult, err := s.oomstore.ChannelJoin(context.Background(), types.ChannelJoinOpt{
			FeatureNames: firstReq.FeatureNames,
			EntityRows:   entityRows,
			ValueNames:   firstReq.ValueNames,
		})
		if err != nil {
			globalErr = err
		} else {
			for row := range joinResult.Data {
				joinedRow, err := convertJoinedRow(row)
				if err != nil {
					globalErr = err
					break
				}
				err = stream.Send(&codegen.ChannelJoinResponse{
					Status:    buildStatus(code.Code_OK, ""),
					Header:    joinResult.Header,
					JoinedRow: joinedRow,
				})
				if err != nil {
					globalErr = err
					break
				}
			}
		}
		done <- struct{}{}
	}()

	// DO NOT move it before the goroutine starts,
	// otherwise it blocks since the channel `entityRows` is not being consumed
	entityRows <- types.EntityRow{
		EntityKey: firstReq.EntityRow.EntityKey,
		UnixMilli: firstReq.EntityRow.UnixMilli,
		Values:    firstReq.EntityRow.Values,
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			globalErr = err
			break
		}
		if globalErr != nil {
			break
		}
		if req.GetEntityRow() == nil {
			globalErr = fmt.Errorf("cannot process nil entity row")
			break
		}
		entityRows <- types.EntityRow{
			EntityKey: req.EntityRow.EntityKey,
			UnixMilli: req.EntityRow.UnixMilli,
			Values:    req.EntityRow.Values,
		}
	}

	close(entityRows)
	// wait until oomstore ChannelJoin is done, whether succeeded or failed
	<-done

	return globalErr
}

func (s *server) Join(ctx context.Context, req *codegen.JoinRequest) (*codegen.JoinResponse, error) {
	err := s.oomstore.Join(ctx, types.JoinOpt{
		FeatureNames:   req.FeatureNames,
		InputFilePath:  req.InputFilePath,
		OutputFilePath: req.OutputFilePath,
	})
	if err != nil {
		return &codegen.JoinResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}

	return &codegen.JoinResponse{
		Status: buildStatus(code.Code_OK, ""),
	}, nil
}

func (s *server) ChannelExport(req *codegen.ChannelExportRequest, stream codegen.OomAgent_ChannelExportServer) error {
	ctx := context.Background()
	exportResult, err := s.oomstore.ChannelExport(ctx, types.ChannelExportOpt{
		FeatureNames: req.FeatureNames,
		RevisionID:   int(req.RevisionId),
		Limit:        req.Limit,
	})
	if err != nil {
		return err
	}
	for row := range exportResult.Data {
		valueRow, err := convertToValueSlice(row)
		if err != nil {
			return err
		}
		if err := stream.Send(&codegen.ChannelExportResponse{
			Status: buildStatus(code.Code_OK, ""),
			Header: exportResult.Header,
			Row:    valueRow,
		}); err != nil {
			return err
		}
	}
	exportErr := exportResult.CheckStreamError()
	if exportErr != nil {
		if err := stream.Send(&codegen.ChannelExportResponse{
			Status: buildStatus(code.Code_INTERNAL, exportErr.Error()),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Export(ctx context.Context, req *codegen.ExportRequest) (*codegen.ExportResponse, error) {
	err := s.oomstore.Export(ctx, types.ExportOpt{
		FeatureNames:   req.FeatureNames,
		RevisionID:     int(req.RevisionId),
		Limit:          req.Limit,
		OutputFilePath: req.OutputFilePath,
	})
	if err != nil {
		return &codegen.ExportResponse{
			Status: buildStatus(code.Code_INTERNAL, err.Error()),
		}, err
	}
	return &codegen.ExportResponse{
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

func convertToValueSlice(s []interface{}) ([]*codegen.Value, error) {
	valueSlice := make([]*codegen.Value, 0, len(s))
	for i := range s {
		value, err := convertInterfaceToValue(i)
		if err != nil {
			return nil, err
		}
		valueSlice = append(valueSlice, value)
	}
	return valueSlice, nil
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
			Kind: &codegen.Value_UnixMilliValue{
				UnixMilliValue: s.UnixMilli(),
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

func convertJoinedRow(row []interface{}) ([]*codegen.Value, error) {
	res := make([]*codegen.Value, 0, len(row))
	for _, value := range row {
		v, err := convertInterfaceToValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %v", value)
		}
		res = append(res, v)
	}
	return res, nil
}
