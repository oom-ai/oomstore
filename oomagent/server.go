package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/oom-ai/oomstore/oomagent/codegen"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

type server struct {
	codegen.UnimplementedOomAgentServer

	oomstore *oomstore.OomStore
}

func (s *server) HealthCheck(ctx context.Context, req *codegen.HealthCheckRequest) (*codegen.HealthCheckResponse, error) {
	if err := s.oomstore.Ping(ctx); err != nil {
		return nil, status.Errorf(codes.Unavailable, "oomstore is currently unavailable")
	}
	return &codegen.HealthCheckResponse{}, nil
}

func (s *server) OnlineGet(ctx context.Context, req *codegen.OnlineGetRequest) (*codegen.OnlineGetResponse, error) {
	result, err := s.oomstore.OnlineGet(ctx, types.OnlineGetOpt{
		FeatureNames: req.Features,
		EntityKey:    req.EntityKey,
	})
	if err != nil {
		return nil, internalError(err.Error())
	}

	valueMap, err := convertToValueMap(result.FeatureValueMap)
	if err != nil {
		return nil, internalError(err.Error())
	}
	return &codegen.OnlineGetResponse{
		Result: &codegen.FeatureValueMap{
			Map: valueMap,
		},
	}, nil
}

func (s *server) OnlineMultiGet(ctx context.Context, req *codegen.OnlineMultiGetRequest) (*codegen.OnlineMultiGetResponse, error) {
	result, err := s.oomstore.OnlineMultiGet(ctx, types.OnlineMultiGetOpt{
		FeatureNames: req.Features,
		EntityKeys:   req.EntityKeys,
	})
	if err != nil {
		return nil, internalError(err.Error())
	}

	resultMap := make(map[string]*codegen.FeatureValueMap)
	for entityKey, featureValues := range result {
		valueMap, err := convertToValueMap(featureValues.FeatureValueMap)
		if err != nil {
			return nil, internalError(err.Error())
		}
		resultMap[entityKey] = &codegen.FeatureValueMap{
			Map: valueMap,
		}
	}
	return &codegen.OnlineMultiGetResponse{
		Result: resultMap,
	}, nil
}

func (s *server) Sync(ctx context.Context, req *codegen.SyncRequest) (*codegen.SyncResponse, error) {
	var revisionID int
	if req.RevisionId != nil {
		revisionID = int(*req.RevisionId)
	}
	if err := s.oomstore.Sync(ctx, types.SyncOpt{
		GroupName:  req.Group,
		RevisionID: &revisionID,
		PurgeDelay: int(req.PurgeDelay),
	}); err != nil {
		return nil, internalError(err.Error())
	}

	return &codegen.SyncResponse{}, nil
}

func (s *server) ChannelImport(stream codegen.OomAgent_ChannelImportServer) error {
	firstReq, err := stream.Recv()
	if err != nil {
		return internalError(err.Error())
	}
	if firstReq.Group == nil {
		return status.Errorf(codes.InvalidArgument, "group is required in first request")
	}
	var description string
	var delimiter rune
	if firstReq.Description != nil {
		description = *firstReq.Description
	}
	delimiter = ','
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
		GroupName:      *firstReq.Group,
		Revision:       firstReq.Revision,
		Description:    description,
		DataSourceType: types.CSV_READER,
		CsvReaderDataSource: &types.CsvReaderDataSource{
			Reader:    reader,
			Delimiter: delimiter,
		},
	})
	if err != nil {
		return internalError(err.Error())
	}
	return stream.SendAndClose(&codegen.ImportResponse{
		RevisionId: int32(revisionID),
	})
}

func (s *server) Push(ctx context.Context, req *codegen.PushRequest) (*codegen.PushResponse, error) {
	if err := s.oomstore.Push(ctx, types.PushOpt{
		EntityKey:     req.EntityKey,
		GroupName:     req.Group,
		FeatureValues: convertToInterfaceMap(req.FeatureValues),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &codegen.PushResponse{}, nil
}

func (s *server) Snapshot(ctx context.Context, re *codegen.SnapshotRequest) (*codegen.SnapshotResponse, error) {
	if err := s.oomstore.Snapshot(ctx, re.Group); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &codegen.SnapshotResponse{}, nil
}

func (s *server) Import(ctx context.Context, req *codegen.ImportRequest) (*codegen.ImportResponse, error) {
	var description string
	var delimiter rune

	if req.Description != nil {
		description = *req.Description
	}
	if req.Delimiter != nil {
		delimiter = []rune(*req.Delimiter)[0]
	} else {
		delimiter = ','
	}
	revisionID, err := s.oomstore.Import(ctx, types.ImportOpt{
		GroupName:      req.Group,
		Description:    description,
		Revision:       req.Revision,
		DataSourceType: types.CSV_FILE,
		CsvFileDataSource: &types.CsvFileDataSource{
			InputFilePath: req.InputFile,
			Delimiter:     delimiter,
		},
	})
	if err != nil {
		return nil, internalError(err.Error())
	}

	return &codegen.ImportResponse{
		RevisionId: int32(revisionID),
	}, nil
}

func (s *server) ChannelJoin(stream codegen.OomAgent_ChannelJoinServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We need to read the first request to get the feature names and value names
	firstReq, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return status.Error(codes.InvalidArgument, "invalid request: empty feature")
		}
		return wrapErr(err)
	}
	if firstReq.GetEntityRow() == nil {
		return nil
	}

	// This channel receives requests from the client.
	entityRows := make(chan types.EntityRow, 1)

	go func() {
		defer close(entityRows)

		entityRows <- types.EntityRow{
			EntityKey: firstReq.EntityRow.EntityKey,
			UnixMilli: firstReq.EntityRow.UnixMilli,
			Values:    firstReq.EntityRow.Values,
		}
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				select {
				case entityRows <- types.EntityRow{Error: err}:
					return
				case <-ctx.Done():
					return
				}
			}
			if req.GetEntityRow() == nil {
				select {
				case entityRows <- types.EntityRow{Error: internalError("cannot process nil entity row")}:
					return
				case <-ctx.Done():
					return
				}
			}

			select {
			case entityRows <- types.EntityRow{
				EntityKey: req.EntityRow.EntityKey,
				UnixMilli: req.EntityRow.UnixMilli,
				Values:    req.EntityRow.Values,
				Error:     nil,
			}:
				// nothing to do
			case <-ctx.Done():
				return
			}
		}
	}()

	// This goroutine runs the join operation, and send whatever joined as the response
	joinResult, err := s.oomstore.ChannelJoin(ctx, types.ChannelJoinOpt{
		JoinFeatureNames:    firstReq.JoinFeatures,
		EntityRows:          entityRows,
		ExistedFeatureNames: firstReq.ExistedFeatures,
	})
	if err != nil {
		return wrapErr(err)
	}

	header := joinResult.Header
	for row := range joinResult.Data {
		if row.Error != nil {
			return row.Error
		}
		joinedRow, err := convertJoinedRow(row.Record)
		if err != nil {
			return wrapErr(err)
		}
		resp := &codegen.ChannelJoinResponse{
			Header:    header,
			JoinedRow: joinedRow,
		}
		if err = stream.Send(resp); err != nil {
			return wrapErr(err)
		}
		// Only need to send header upon the first response
		header = nil
	}
	return nil
}

func (s *server) Join(ctx context.Context, req *codegen.JoinRequest) (*codegen.JoinResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	joinResult, err := s.oomstore.Join(ctx, types.JoinOpt{
		FeatureNames:  req.Features,
		InputFilePath: req.InputFile,
	})
	if err != nil {
		return nil, internalError(err.Error())
	}

	if err := writeJoinResultToFile(req.OutputFile, joinResult); err != nil {
		return nil, wrapErr(err)
	}

	return &codegen.JoinResponse{}, nil
}

func writeJoinResultToFile(outputFilePath string, joinResult *types.JoinResult) error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	if err := w.Write(joinResult.Header); err != nil {
		return err
	}
	for row := range joinResult.Data {
		if row.Error != nil {
			return row.Error
		}
		if err := w.Write(joinRecord(row.Record)); err != nil {
			return err
		}
	}
	return nil
}

func joinRecord(row []interface{}) []string {
	record := make([]string, 0, len(row))
	for _, value := range row {
		record = append(record, cast.ToString(value))
	}
	return record
}

func (s *server) ChannelExport(req *codegen.ChannelExportRequest, stream codegen.OomAgent_ChannelExportServer) error {
	if len(req.Features) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exportResult, err := s.oomstore.ChannelExport(ctx, types.ChannelExportOpt{
		FeatureNames: req.Features,
		UnixMilli:    req.UnixMilli,
		Limit:        req.Limit,
	})
	if err != nil {
		return internalError(err.Error())
	}

	header := exportResult.Header
	for row := range exportResult.Data {
		if row.Error != nil {
			return row.Error
		}

		valueRow, err := convertToValueSlice(row.Record)
		if err != nil {
			return internalError(err.Error())
		}
		if err := stream.Send(&codegen.ChannelExportResponse{
			Header: header,
			Row:    valueRow,
		}); err != nil {
			return err
		}
		// Only need to send header upon the first response
		header = nil
	}
	return nil
}

func (s *server) Export(ctx context.Context, req *codegen.ExportRequest) (*codegen.ExportResponse, error) {
	err := s.oomstore.Export(ctx, types.ExportOpt{
		FeatureNames:   req.Features,
		UnixMilli:      req.UnixMilli,
		Limit:          req.Limit,
		OutputFilePath: req.OutputFile,
	})
	if err != nil {
		return nil, internalError(err.Error())
	}
	return &codegen.ExportResponse{}, nil
}

func convertToValueMap(m map[string]interface{}) (map[string]*codegen.Value, error) {
	valueMap := make(map[string]*codegen.Value, len(m))
	for key, i := range m {
		value, err := convertInterfaceToValue(i)
		if err != nil {
			return nil, err
		}
		valueMap[key] = value
	}
	return valueMap, nil
}

func convertToInterfaceMap(m map[string]*codegen.Value) map[string]interface{} {
	rs := make(map[string]interface{}, len(m))
	for k, v := range m {
		rs[k] = convertValueToInterface(v)
	}
	return rs
}

func convertToValueSlice(s []interface{}) ([]*codegen.Value, error) {
	valueSlice := make([]*codegen.Value, 0, len(s))
	for _, i := range s {
		value, err := convertInterfaceToValue(i)
		if err != nil {
			return nil, err
		}
		valueSlice = append(valueSlice, value)
	}
	return valueSlice, nil
}

func convertInterfaceToValue(i interface{}) (*codegen.Value, error) {
	switch s := i.(type) {
	case nil:
		return nil, nil
	case int64:
		return &codegen.Value{
			Value: &codegen.Value_Int64{
				Int64: s,
			},
		}, nil
	case float64:
		return &codegen.Value{
			Value: &codegen.Value_Double{
				Double: s,
			},
		}, nil
	case string:
		return &codegen.Value{
			Value: &codegen.Value_String_{
				String_: s,
			},
		}, nil
	case bool:
		return &codegen.Value{
			Value: &codegen.Value_Bool{
				Bool: s,
			},
		}, nil
	case time.Time:
		return &codegen.Value{
			Value: &codegen.Value_UnixMilli{
				UnixMilli: s.UnixMilli(),
			},
		}, nil
	case []byte:
		return &codegen.Value{
			Value: &codegen.Value_Bytes{
				Bytes: s,
			},
		}, nil
	default:
		return nil, errdefs.Errorf("unsupported value type %T", i)
	}
}

func convertJoinedRow(row []interface{}) ([]*codegen.Value, error) {
	res := make([]*codegen.Value, 0, len(row))
	for _, value := range row {
		v, err := convertInterfaceToValue(value)
		if err != nil {
			return nil, errdefs.Errorf("failed to marshal %v", value)
		}
		res = append(res, v)
	}
	return res, nil
}

func convertValueToInterface(i *codegen.Value) interface{} {
	if i == nil {
		return nil
	}
	kind := i.GetValue()
	switch kind.(type) {
	case *codegen.Value_Int64:
		return i.GetInt64()
	case *codegen.Value_Double:
		return i.GetDouble()
	case *codegen.Value_String_:
		return i.GetString_()
	case *codegen.Value_Bool:
		return i.GetBool()
	case *codegen.Value_UnixMilli:
		return i.GetUnixMilli()
	case *codegen.Value_Bytes:
		return i.GetBytes()
	default:
		panic(fmt.Sprintf("unsupported value type: %T, value=%v", i, i))
	}
}

func internalError(msg string) error {
	return status.Errorf(codes.Internal, msg)
}

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Internal, err.Error())
}
