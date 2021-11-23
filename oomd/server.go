package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/oom-ai/oomstore/oomd/codegen"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	code "google.golang.org/genproto/googleapis/rpc/code"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
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

func generateBytesFrom(rows []*anypb.Any) []byte {
	var res []byte
	for _, row := range rows {
		res = append(res, row.Value...)
	}
	return res
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

		if _, err := writer.Write(generateBytesFrom(firstReq.Row)); err != nil {
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
			if _, err := writer.Write(generateBytesFrom(req.Row)); err != nil {
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
	firstReq, err := stream.Recv()
	if err != nil {
		return err
	}

	entityRows := make(chan types.EntityRow)

	joinResult, err := s.oomstore.Join(context.Background(), types.JoinOpt{
		FeatureNames: firstReq.FeatureNames,
		EntityRows:   entityRows,
	})
	if err != nil {
		return err
	}

	var sendErr error
	go func() {
		for row := range joinResult.Data {
			joinedRow, err := convertJoinedRow(row)
			if err != nil {
				sendErr = err
				break
			}
			err = stream.Send(&codegen.JoinResponse{
				Status:    buildStatus(code.Code_OK, ""),
				Header:    joinResult.Header,
				JoinedRow: joinedRow,
			})
			if err != nil {
				sendErr = err
				break
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if sendErr != nil {
			return sendErr
		}
		entityRows <- types.EntityRow{
			EntityKey: req.EntityRow.EntityKey,
			UnixTime:  req.EntityRow.UnixTime,
		}
	}
}

func (s *server) ImportByFile(context.Context, *codegen.ImportByFileRequest) (*codegen.ImportResponse, error) {
	panic("implement me")
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

func convertJoinedRow(row []interface{}) ([]*anypb.Any, error) {
	res := make([]*anypb.Any, 0, len(row))
	for _, value := range row {
		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %v", value)
		}
		res = append(res, &anypb.Any{
			Value: bytes,
		})
	}
	return res, nil
}

// The methods below serves the temporary testing purpose
// We should add more formal tests soon
func newServer() *server {
	store, err := oomstore.Open(context.Background(), buildOomStoreConfig())
	if err != nil {
		log.Fatal(err)
	}
	return &server{oomstore: store}
}

func buildOomStoreConfig() types.OomStoreConfig {
	return types.OomStoreConfig{
		MetadataStore: types.MetadataStoreConfig{
			Backend:  types.POSTGRES,
			Postgres: getPostgresOpt(),
		},
		OnlineStore: types.OnlineStoreConfig{
			Backend:  types.POSTGRES,
			Postgres: getPostgresOpt(),
		},
		OfflineStore: types.OfflineStoreConfig{
			Backend:  types.POSTGRES,
			Postgres: getPostgresOpt(),
		},
	}
}

func getPostgresOpt() *types.PostgresOpt {
	return &types.PostgresOpt{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Database: "oomstore",
	}
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	codegen.RegisterOomDServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
