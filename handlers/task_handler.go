package handlers

import (
	"net/http"

	"github.com/kuruteiru/gotodo/models"
	"github.com/kuruteiru/gotodo/renderer"
)

func ViewTodolist(w http.ResponseWriter, r *http.Request) {
	pd := &renderer.PageData{
		Title: "todolist",
		Data: map[string][]models.Task{
			"Tasks": models.GenerateTasks(),
		},
	}
	renderer.RenderTemplate(w, "todolist", pd)
}
