package data

import "sync"

type Todo struct {
	ID         int     `json:"ID"`
	Content    string  `json:"content"`
	CreateTime string  `json:"createTime"`
	FinishTime *string `json:"finishTime"`
	Priority   *int    `json:"Priority"`
	Done       bool    `json:"done"`
}

type Repository interface {
	CreateAndAddTodo(content string, done bool) (*Todo, error)
	AddTodos([]*Todo) ([]*Todo, error)
	GetTodos() ([]*Todo, error)
	GetTodoById(id int) (*Todo, error)
	ModifyTodo(id int, todo *Todo) error
	RemoveTodos(ids []int) ([]*Todo, error)
	Size() int
}

var CreateRepository = sync.OnceValue(func() Repository {
	return NewLocalRepository()
})
