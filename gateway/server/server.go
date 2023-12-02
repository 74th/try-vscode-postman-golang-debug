package server

import (
	"github.com/74th/vscode-book-r2-golang/domain/usecase"
	"github.com/74th/vscode-book-r2-golang/gateway/memdb"
	"github.com/74th/vscode-book-r2-golang/gateway/server/grpc"
	"github.com/74th/vscode-book-r2-golang/gateway/server/openapi"
)

// サーバAPI
type Server struct {
	openapi    *openapi.Server
	grpc       *grpc.Server
	interactor *usecase.Interactor
}

// サーバAPIのインスタンスを作成する
func New(restAddr string, grpcAddr string, webroot string) *Server {
	interactor := &usecase.Interactor{
		Database: memdb.NewDB(),
	}

	sv1 := openapi.New(restAddr, webroot, interactor)
	sv2 := grpc.New(grpcAddr, interactor)

	s := &Server{
		openapi:    sv1,
		grpc:       sv2,
		interactor: interactor,
	}

	return s
}

// サーバを開始する
func (s *Server) Serve() error {
	if err := s.openapi.Serve(); err != nil {
		return err
	}

	if err := s.grpc.Serve(); err != nil {
		return err
	}

	return nil
}

// サーバを開始する
func (s *Server) Close() error {
	if err := s.openapi.Close(); err != nil {
		return err
	}

	if err := s.openapi.Close(); err != nil {
		return err
	}

	return nil
}
