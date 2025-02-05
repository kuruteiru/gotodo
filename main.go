package main

import (
	"log"
	"os"

	"github.com/kuruteiru/gotodo/db"
	"github.com/kuruteiru/gotodo/router"
	"github.com/kuruteiru/gotodo/server"
)

func initServer() {
	db.Main()
	return

	logger := log.New(os.Stdout, "gotodo: ", log.LstdFlags)
	
	r := router.Route()
	r = server.Logging(logger)(r)
	r = server.Tracing(server.NextRequestID)(r)

	server.Serve(&r, logger)
}

func main() {
	initServer()
	os.Exit(0)
}
