package handlers

import (
	"fmt"
	"net/http"
	"strings"
    "github.com/kuruteiru/gotodo/models"
)

func ViewTodolist(w http.ResponseWriter, r *http.Request) {
    tasks := models.GenerateTasks()
    // PrintTasks(tasks)

    var sb strings.Builder
    for i, task := range tasks {
        fmt.Fprintf(&sb, "%v: %+v\n", i, task)
    }

    fmt.Fprintf(w, "%v", sb.String())
}
