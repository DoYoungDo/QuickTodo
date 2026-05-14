package data

import (
	"sync"
	"todo_list/internal/setting"
)

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
	ModifyTodo(id int, todo Todo) error
	RemoveTodos(ids []int) ([]*Todo, error)
	ClearTodos() ([]*Todo, error)
	Size() int
}

var CreateRepository = sync.OnceValue(func() Repository {
	st, err := setting.Get()
	if err != nil {
		panic(err)
	}
	repository, err := NewRepositoryWithSetting(st)
	if err != nil {
		panic(err)
	}
	return repository
})

func NewRepositoryWithSetting(st setting.Setting) (Repository, error) {
	switch st.Get(setting.KeyRepositoryName) {
	case setting.RepositoryLocal:
		return NewLocalRepositoryWithSetting(st)
	default:
		return NewLocalRepositoryWithSetting(st)
	}
}
