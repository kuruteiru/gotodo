package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kuruteiru/gotodo/models"
	"github.com/kuruteiru/gotodo/renderer"
)

func ViewTaskPage(w http.ResponseWriter, r *http.Request) {
	p := fmt.Sprintf("task/%s", r.PathValue("page"))
	renderer.RenderTemplate(w, p, nil)
}

//todo: get task by id from db and display it's detail
func ViewTaskDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		renderer.RenderTemplate(w, "task/bad-id", nil)
		return
	}

	//dummy task, todo: switch for something like db.GetTask()
	t := models.GetTask(uint64(id))

	pd := &renderer.PageData{
		Title: fmt.Sprintf("task %v", t.ID),
		Data: map[string]models.Task{
			"Task": *t,
		},
	}

	renderer.RenderTemplate(w, "task/detail", pd)
}

//todo: save task to db
func CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusConflict)
	}
	f := r.PostForm

	//create task in memory
	name := f.Get("name")
	description := f.Get("description")

	var priority models.TaskPriority
	fmt.Printf("default tp: %v\n", priority)
	ps, err := strconv.Atoi(f.Get("priority"))
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	}
	priority = models.TaskPriority(uint8(ps))
	fmt.Printf("given tp: %v\n", priority)

	t := models.NewTask(name, description, priority, 0)
	
	//upload task to db
	//todo: implement permanent storage for tasks

	//get pageData and render detail
	pd := &renderer.PageData{
		Title: fmt.Sprintf("task %v", t.ID),
		Data: map[string]models.Task{
			"Task": t,
		},
	}

	renderer.RenderTemplate(w, "task/detail", pd)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) { }

func ViewTodolist(w http.ResponseWriter, r *http.Request) {
	pd := &renderer.PageData{
		Title: "todolist",
		Data: map[string][]models.Task{
			"Tasks": models.GenerateTasks(10),
		},
	}
	renderer.RenderTemplate(w, "todolist", pd)
}
