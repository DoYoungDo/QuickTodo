package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"todo_list/internal/setting"
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

func NewLocalRepositoryWithSetting(st setting.Setting) (*LocalRepository, error) {
	table := st.Get(setting.KeyRepositoryLocalTable)
	todoDataFile := filepath.Join(st.AppDataDir(), "todos", table+".json")
	return newLocalRepository(todoDataFile, table)
}

func NewLocalRepositoryWithPath(todoDataFile string) (*LocalRepository, error) {
	return newLocalRepository(todoDataFile, listNameFromDataFile(todoDataFile))
}

func listNameFromDataFile(todoDataFile string) string {
	base := filepath.Base(todoDataFile)
	if ext := filepath.Ext(base); ext != "" {
		return strings.TrimSuffix(base, ext)
	}
	return base
}

func newLocalRepository(todoDataFile string, listName string) (*LocalRepository, error) {
	appDataDir := filepath.Dir(todoDataFile)
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return nil, err
	}

	todoDataList := &localTodoList{
		Name: listName,
		Date: time.Now().String(),
		List: []*Todo{},
	}
	if _, err := os.Stat(todoDataFile); !os.IsNotExist(err) {
		if todoData, err := os.ReadFile(todoDataFile); err == nil {
			if err := json.Unmarshal(todoData, todoDataList); err != nil {
				return nil, fmt.Errorf("load todo data: %w", err)
			}
		} else {
			return nil, err
		}
	}
	return &LocalRepository{
		list:     todoDataList,
		dataFile: todoDataFile,
	}, nil
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
	return l.cloneTodo(todo), err
}

// AddTodos implements [Repository].
func (l *LocalRepository) AddTodos(todos []*Todo) ([]*Todo, error) {
	l.list.List = append(l.list.List, slices.Collect(func() iter.Seq[*Todo] {
		return func(yield func(*Todo) bool) {
			for _, t := range todos {
				if !yield(l.cloneTodo(t)) {
					return
				}
			}
		}
	}())...)
	return todos, l.flushToLocal()
}

// GetTodos implements [Repository].
func (l *LocalRepository) GetTodos() ([]*Todo, error) {
	return slices.Collect(func() iter.Seq[*Todo] {
		return func(yield func(*Todo) bool) {
			for _, t := range l.list.List {
				if !yield(l.cloneTodo(t)) {
					return
				}
			}
		}
	}()), nil
}

func (l *LocalRepository) GetTodoById(id int) (*Todo, error) {
	for _, todo := range l.list.List {
		if todo.ID == id {
			return l.cloneTodo(todo), nil
		}
	}
	return nil, errors.New("todo not found")
}

// ModifyTodo implements [Repository].
func (l *LocalRepository) ModifyTodo(id int, todo Todo) (*Todo, error) {
	index := slices.IndexFunc(l.list.List, func(t *Todo) bool {
		return t.ID == id
	})
	if index == -1 {
		return nil, errors.New("todo not found")
	}
	l.list.List[index] = func() *Todo {
		newTodo := l.cloneTodo(l.list.List[index])
		newTodo.Content = todo.Content
		newTodo.CreateTime = todo.CreateTime
		newTodo.Priority = todo.Priority

		if newTodo.Done != todo.Done {
			newTodo.Done = todo.Done

			if todo.Done {
				if todo.FinishTime == nil {
					timeNow := time.Now().Format(time.RFC3339)
					newTodo.FinishTime = &timeNow
				} else {
					newTodo.FinishTime = todo.FinishTime
				}
			} else {
				newTodo.FinishTime = nil
			}
		}

		return newTodo
	}()
	return l.cloneTodo(l.list.List[index]), l.flushToLocal()
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

func (l *LocalRepository) ClearTodos() ([]*Todo, error) {
	clearedTodos := l.list.List
	l.list.List = []*Todo{}
	if err := os.Remove(l.dataFile); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return clearedTodos, nil
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

// cloneTodo ensures returned todo pointers cannot mutate repository-owned data.
func (l *LocalRepository) cloneTodo(todo *Todo) *Todo {
	var finishTime *string
	if todo.FinishTime != nil {
		val := *todo.FinishTime
		finishTime = &val
	}
	var priority *int
	if todo.Priority != nil {
		val := *todo.Priority
		priority = &val
	}
	return &Todo{
		ID:         todo.ID,
		Content:    todo.Content,
		CreateTime: todo.CreateTime,
		FinishTime: finishTime,
		Priority:   priority,
		Done:       todo.Done,
	}
}

var _ Repository = (*LocalRepository)(nil)
