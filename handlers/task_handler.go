package handlers

import (
	"fmt"
	"net/http"

	"github.com/kuruteiru/gotodo/models"
	"github.com/kuruteiru/gotodo/renderer"
)

func ViewTask(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("page")
	p = "task/"+p
	renderer.RenderTemplate(w, p, nil)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusConflict)
	}
	f := r.PostForm
	fmt.Printf("%+v", f)
	// renderer.RenderTemplate(w, p, nil)
}

func ViewTodolist(w http.ResponseWriter, r *http.Request) {
	pd := &renderer.PageData{
		Title: "todolist",
		Data: map[string][]models.Task{
			"Tasks": models.GenerateTasks(),
		},
	}
	renderer.RenderTemplate(w, "todolist", pd)
}
