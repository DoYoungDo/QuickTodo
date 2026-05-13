package data

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestRepository(t *testing.T) (*LocalRepository, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "todos", "default.json")
	repo, err := NewLocalRepositoryWithPath(path)
	if err != nil {
		t.Fatalf("NewLocalRepositoryWithPath() error = %v", err)
	}
	return repo, path
}

func TestLocalRepository_CreateAndAddTodo(t *testing.T) {
	repo, _ := newTestRepository(t)

	first, err := repo.CreateAndAddTodo("write tests", false)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	if first.ID != 0 || first.Content != "write tests" || first.Done || first.FinishTime != nil {
		t.Fatalf("unexpected first todo: %+v", first)
	}
	if _, err := time.Parse(time.RFC3339, first.CreateTime); err != nil {
		t.Fatalf("CreateTime is not RFC3339: %v", err)
	}

	second, err := repo.CreateAndAddTodo("ship release", true)
	if err != nil {
		t.Fatalf("CreateAndAddTodo(done) error = %v", err)
	}
	if second.ID != 1 || !second.Done || second.FinishTime == nil {
		t.Fatalf("unexpected done todo: %+v", second)
	}
	if _, err := time.Parse(time.RFC3339, *second.FinishTime); err != nil {
		t.Fatalf("FinishTime is not RFC3339: %v", err)
	}
}

func TestLocalRepository_GetTodosReturnsClones(t *testing.T) {
	repo, _ := newTestRepository(t)
	created, err := repo.CreateAndAddTodo("original", true)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	priority := 3
	created.Priority = &priority
	if err := repo.ModifyTodo(created.ID, created); err != nil {
		t.Fatalf("ModifyTodo() error = %v", err)
	}

	todos, err := repo.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos() error = %v", err)
	}
	todos[0].Content = "changed"
	*todos[0].Priority = 1
	*todos[0].FinishTime = "2000-01-01T00:00:00Z"

	again, err := repo.GetTodoById(created.ID)
	if err != nil {
		t.Fatalf("GetTodoById() error = %v", err)
	}
	if again.Content != "original" || *again.Priority != 3 || *again.FinishTime == "2000-01-01T00:00:00Z" {
		t.Fatalf("repository returned mutable todo data: %+v", again)
	}
}

func TestLocalRepository_GetTodoById(t *testing.T) {
	repo, _ := newTestRepository(t)
	created, err := repo.CreateAndAddTodo("find me", false)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}

	found, err := repo.GetTodoById(created.ID)
	if err != nil {
		t.Fatalf("GetTodoById() error = %v", err)
	}
	if found.Content != "find me" {
		t.Fatalf("found.Content = %q", found.Content)
	}

	if _, err := repo.GetTodoById(999); err == nil {
		t.Fatal("GetTodoById() expected error for missing id")
	}
}

func TestLocalRepository_ModifyTodoPersistsChanges(t *testing.T) {
	repo, path := newTestRepository(t)
	todo, err := repo.CreateAndAddTodo("before", false)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	priority := 5
	todo.Content = "after"
	todo.Done = true
	todo.Priority = &priority
	now := time.Now().Format(time.RFC3339)
	todo.FinishTime = &now
	if err := repo.ModifyTodo(todo.ID, todo); err != nil {
		t.Fatalf("ModifyTodo() error = %v", err)
	}

	reopened, err := NewLocalRepositoryWithPath(path)
	if err != nil {
		t.Fatalf("NewLocalRepositoryWithPath(reopen) error = %v", err)
	}
	got, err := reopened.GetTodoById(todo.ID)
	if err != nil {
		t.Fatalf("GetTodoById() error = %v", err)
	}
	if got.Content != "after" || !got.Done || got.Priority == nil || *got.Priority != 5 || got.FinishTime == nil {
		t.Fatalf("unexpected persisted todo: %+v", got)
	}
}

func TestLocalRepository_ModifyTodoMissing(t *testing.T) {
	repo, _ := newTestRepository(t)
	if err := repo.ModifyTodo(999, &Todo{ID: 999, Content: "missing"}); err == nil {
		t.Fatal("ModifyTodo() expected error for missing id")
	}
}

func TestLocalRepository_RemoveTodosReindexesIDs(t *testing.T) {
	repo, _ := newTestRepository(t)
	for _, content := range []string{"zero", "one", "two"} {
		if _, err := repo.CreateAndAddTodo(content, false); err != nil {
			t.Fatalf("CreateAndAddTodo(%q) error = %v", content, err)
		}
	}

	removed, err := repo.RemoveTodos([]int{1})
	if err != nil {
		t.Fatalf("RemoveTodos() error = %v", err)
	}
	if len(removed) != 1 || removed[0].Content != "one" {
		t.Fatalf("unexpected removed todos: %+v", removed)
	}

	todos, err := repo.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos() error = %v", err)
	}
	if len(todos) != 2 || todos[0].ID != 0 || todos[0].Content != "zero" || todos[1].ID != 1 || todos[1].Content != "two" {
		t.Fatalf("unexpected remaining todos: %+v", todos)
	}
}

func TestLocalRepository_RemoveTodosMissingID(t *testing.T) {
	repo, _ := newTestRepository(t)
	if _, err := repo.CreateAndAddTodo("keep", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}

	removed, err := repo.RemoveTodos([]int{999})
	if err != nil {
		t.Fatalf("RemoveTodos() error = %v", err)
	}
	if len(removed) != 0 || repo.Size() != 1 {
		t.Fatalf("unexpected remove missing result: removed=%+v size=%d", removed, repo.Size())
	}
}

func TestLocalRepository_LoadInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")
	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if _, err := NewLocalRepositoryWithPath(path); err == nil {
		t.Fatal("NewLocalRepositoryWithPath() expected invalid JSON error")
	}
}
