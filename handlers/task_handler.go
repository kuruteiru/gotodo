package handlers

import (
	"net/http"
    "github.com/kuruteiru/gotodo/models"
    "github.com/kuruteiru/gotodo/renderer"
)

func ViewTodolist(w http.ResponseWriter, r *http.Request) {
    renderer.RenderTemplate(w, "todolist", struct {
        Tasks []models.Task
    }{
        Tasks: models.GenerateTasks(),
    })
}
