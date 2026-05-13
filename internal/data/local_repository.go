package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"todo_list/internal/app"
)

type localTodoList struct {
	Name string  `json:"name"`
	Date string  `json:"date"`
	List []*Todo `json:"list"`
}

type LocalRepository struct {
	dataFile string
	list     *localTodoList
}

func NewLocalRepository() *LocalRepository {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	appDataDir += fmt.Sprintf("/%v/todos/", app.APP_NAME)
	if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(appDataDir, 0755); err != nil {
			panic(err)
		}
	}
	todoDataFile := fmt.Sprintf("%sdefault.json", appDataDir)
	todoDataList := &localTodoList{
		Name: "default",
		Date: time.Now().String(),
		List: []*Todo{},
	}
	if _, err := os.Stat(todoDataFile); !os.IsNotExist(err) {
		if todoData, err := os.ReadFile(todoDataFile); err == nil {
			json.Unmarshal(todoData, todoDataList)
		}
	}
	return &LocalRepository{
		list:     todoDataList,
		dataFile: todoDataFile,
	}
}

// CreateAndAddTodo implements [Repository].
func (l *LocalRepository) CreateAndAddTodo(content string, done bool) (*Todo, error) {
	todo := &Todo{
		ID:         len(l.list.List),
		Content:    content,
		CreateTime: time.Now().Format(time.RFC3339),
		Done:       done,
	}
	if done {
		todo.FinishTime = &todo.CreateTime
	}

	_, err := l.AddTodos([]*Todo{todo})
	return todo, err
}

// AddTodos implements [Repository].
func (l *LocalRepository) AddTodos(todos []*Todo) ([]*Todo, error) {
	l.list.List = append(l.list.List, todos...)
	return todos, l.flushToLocal()
}

// GetTodos implements [Repository].
func (l *LocalRepository) GetTodos() ([]*Todo, error) {
	return l.list.List, nil
}

// ModifyTodo implements [Repository].
func (l *LocalRepository) ModifyTodo(id int, todo *Todo) error {
	panic("unimplemented")
}

// RemoveTodo implements [Repository].
func (l *LocalRepository) RemoveTodos(ids []int) ([]*Todo, error) {
	idsMap := make(map[int]interface{}, len(ids))
	for _, id := range ids {
		idsMap[id] = struct{}{}
	}

	toRemoveTodos := make([]*Todo, 0, len(ids))
	lastTodos := make([]*Todo, 0, len(l.list.List))

	for _, todo := range l.list.List {
		if _, exist := idsMap[todo.ID]; exist {
			toRemoveTodos = append(toRemoveTodos, todo)
		} else {
			lastTodos = append(lastTodos, todo)
		}
	}

	if len(lastTodos) > 0 {
		for i, todo := range lastTodos {
			todo.ID = i
		}
	}
	l.list.List = lastTodos

	return toRemoveTodos, l.flushToLocal()
}

func (l *LocalRepository) Size() int {
	return len(l.list.List)
}

func (l *LocalRepository) flushToLocal() error {
	data, err := json.Marshal(l.list)
	if err != nil {
		return err
	}

	return os.WriteFile(l.dataFile, data, 0644)
}

var _ Repository = (*LocalRepository)(nil)
