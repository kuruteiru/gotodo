package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/kuruteiru/gotodo/renderer"
	"github.com/kuruteiru/gotodo/server"
)

func ViewWrongPage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "wrong page")
}

func ViewIndex(w http.ResponseWriter, r *http.Request) {
    renderer.RenderTemplate(w, "index", nil)
}

func ViewHealtz(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&server.Healthy) == 1 {
        w.WriteHeader(http.StatusNoContent)
        return
    }
    w.WriteHeader(http.StatusServiceUnavailable)
}
