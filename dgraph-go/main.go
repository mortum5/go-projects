package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mortum5/go-projects/dgraph-go/dgraph"
	server "github.com/mortum5/go-projects/dgraph-go/http"
)

func main() {

	dg, cancel := dgraph.New()
	defer cancel()
	dg.Setup()

	server := server.New(dg)

	server.Run()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)

	select {
	case <-exit:
		server.Stop()
		dg.DropAll()
	case err := <-server.Error():
		log.Fatalf("server error: %v", err)
	}
}
