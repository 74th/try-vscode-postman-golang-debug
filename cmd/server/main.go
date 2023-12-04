package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/74th/vscode-book-r2-golang/domain/usecase"
	"github.com/74th/vscode-book-r2-golang/gateway/memdb"
	"github.com/74th/vscode-book-r2-golang/gateway/server"
	"github.com/74th/vscode-book-r2-golang/gateway/token"
)

func main() {
	var (
		webroot       string
		restAddr      string
		grpcAddr      string
		validatorHost string
	)

	flag.StringVar(&webroot, "webroot", "./public", "web root path")
	flag.StringVar(&restAddr, "restAddr", "0.0.0.0:8000", "server addr")
	flag.StringVar(&grpcAddr, "grpcAddr", "0.0.0.0:8001", "server addr")
	flag.StringVar(&validatorHost, "validatorHost", "https://a274ebfe-ee67-4c03-a2db-f278c2535a83.mock.pstmn.io", "server addr")
	flag.Parse()

	interactor := &usecase.Interactor{
		Database: memdb.NewDB(),
	}

	tokenValidator := token.New(validatorHost)

	svr := server.New(restAddr, grpcAddr, webroot, interactor, tokenValidator)
	err := svr.Serve()
	if err != nil {
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGTERM, syscall.SIGQUIT:
				exitChan <- 0
			default:
				exitChan <- 1
			}
		}
	}()
	code := <-exitChan
	os.Exit(code)

}
