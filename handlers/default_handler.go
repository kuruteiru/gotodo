package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/kuruteiru/gotodo/renderer"
	"github.com/kuruteiru/gotodo/server"
)


func ViewPage(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("page")
	renderer.RenderTemplate(w, p, nil)
}

func ViewNoContent(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "nocontent", nil)
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
