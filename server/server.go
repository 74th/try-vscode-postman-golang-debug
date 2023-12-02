package server

import (
	"github.com/74th/vscode-book-r2-golang/domain/usecase"
	"github.com/74th/vscode-book-r2-golang/memdb"
	"github.com/74th/vscode-book-r2-golang/server/openapi"
)

// サーバAPI
type Server struct {
	openapi    *openapi.Server
	interactor *usecase.Interactor
}

// サーバAPIのインスタンスを作成する
func New(addr string, webroot string) *Server {
	interactor := &usecase.Interactor{
		Database: memdb.NewDB(),
	}

	sv1 := openapi.New(addr, webroot, interactor)

	s := &Server{
		openapi:    sv1,
		interactor: interactor,
	}

	return s
}

// サーバを開始する
func (s *Server) Serve() error {
	if err := s.openapi.Serve(); err != nil {
		return err
	}
	return nil
}
