package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/74th/vscode-book-r2-golang/server"
)

func main() {
	var (
		webroot string
		addr    string
	)

	flag.StringVar(&webroot, "webroot", "./public", "web root path")
	flag.StringVar(&addr, "addr", "0.0.0.0:8000", "server addr")
	flag.Parse()

	svr := server.New(addr, webroot)
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
