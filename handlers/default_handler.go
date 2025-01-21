package handlers

import (
	"net/http"
	"sync/atomic"
    "strings"
    "fmt"

	"github.com/kuruteiru/gotodo/renderer"
	"github.com/kuruteiru/gotodo/server"
)

func ViewWrongPage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "wrong page")
}

func ViewIndex(w http.ResponseWriter, r *http.Request) {
    renderer.Init()
    if renderer.Tmpl == nil {
        fmt.Printf("empty tmpl\n")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    renderer.Tmpl.Execute(w, nil)
    return

    var sb strings.Builder
    for i, tmp := range renderer.Tmpl.Templates() {
        fmt.Fprintf(&sb, "tmpl %v: %+v\n", i, tmp.Name())
    }

    fmt.Fprintf(w, "%v\n", sb.String())
}

func ViewHealtz(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&server.Healthy) == 1 {
        w.WriteHeader(http.StatusNoContent)
        return
    }
    w.WriteHeader(http.StatusServiceUnavailable)
}
