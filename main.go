package main

import (
    "log"
    "os"
)

func initServer() {
    logger := log.New(os.Stdout, "gotodo: ", log.LstdFlags)

    router := route()
    router = logging(logger)(router)
    router = tracing(nextRequestID)(router)

    serve(&router, logger)
}

func main() {
    initServer()
    os.Exit(0)
}
