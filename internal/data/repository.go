package data

import "sync"

type Todo struct {
	ID         int     `json:"ID"`
	Content    string  `json:"content"`
	CreateTime string  `json:"createTime"`
	FinishTime *string `json:"finishTime"`
	Done       bool    `json:"done"`
}

type Repository interface {
	CreateAndAddTodo(content string, done bool) (*Todo, error)
	AddTodos([]*Todo) error
	GetTodos() ([]*Todo, error)
	ModifyTodo(id int, todo *Todo) error
	RemoveTodo(id int) error
	Size() int
}

var CreateRepository = sync.OnceValue(func() Repository {
	return NewLocalRepository()
})
