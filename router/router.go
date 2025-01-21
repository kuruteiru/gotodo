package router

import (
    "github.com/kuruteiru/gotodo/handlers"
    "net/http"
)

func Route() http.Handler {
    router := http.NewServeMux()

    fs := http.FileServer(http.Dir("static/"))
    router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("GET /", handlers.ViewWrongPage)
	router.HandleFunc("GET /{$}", handlers.ViewIndex)
	router.HandleFunc("GET /healtz", handlers.ViewHealtz)
    router.HandleFunc("GET /todolist", handlers.ViewTodolist)

    return router
}