package main

import (
    "log"
    "os"
)

func main() {
    logger := log.New(os.Stdout, "gotodo: ", log.LstdFlags)

    router := route()
    router = logging(logger)(router)
    router = tracing(nextRequestID)(router)

    serve(&router, logger)
}
