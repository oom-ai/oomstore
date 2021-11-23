package oomd

import (
	"context"

	"github.com/oom-ai/oomstore/oomd/codegen"
)

type server struct {
	codegen.UnimplementedOomDServer
}

func (s *server) OnlineGet(ctx context.Context, req *codegen.OnlineGetRequest) (*codegen.OnlineGetResponse, error) {
	panic("implement me")
}

func (s *server) OnlineMultiGet(ctx context.Context, req *codegen.OnlineMultiGetRequest) (*codegen.OnlineMultiGetResponse, error) {
	panic("implement me")
}

func (s *server) Sync(context.Context, *codegen.SyncRequest) (*codegen.SyncResponse, error) {
	panic("implement me")
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
