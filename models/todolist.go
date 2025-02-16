package models

import (
	"time"
)

type Todolist struct {
	ID        uint64
	Name      string
	Created   time.Time
	Completed *time.Time
	UserID    uint64
}

func NewTodolist(name string, userID uint64) Todolist {
	//todo: check if userID exists?
	return Todolist{
		Name:    name,
		Created: time.Now().UTC(),
		UserID:  userID,
	}
}
