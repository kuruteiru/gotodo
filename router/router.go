package router

import (
	"net/http"

	"github.com/kuruteiru/gotodo/handlers"
)

func Route() http.Handler {
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("GET /", handlers.ViewNoContent)
	router.HandleFunc("GET /{$}", handlers.ViewIndex)
	router.HandleFunc("GET /{page}", handlers.ViewPage)

	router.HandleFunc("GET /healtz", handlers.ViewHealtz)
	router.HandleFunc("GET /todolist", handlers.ViewTodolist)

	router.HandleFunc("GET /task/{page}", handlers.ViewTaskPage)
	router.HandleFunc("GET /task/detail/{id}", handlers.ViewTaskDetail)

	router.HandleFunc("POST /task/create", handlers.CreateTask)
	// router.HandleFunc("POST /task/update", handlers.UpdateTask)

	return router
}
