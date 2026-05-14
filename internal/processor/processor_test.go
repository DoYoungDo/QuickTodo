package processor

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"todo_list/internal/data"
	"todo_list/internal/ui"
)

func newProcessorTestRepository(t *testing.T) *data.LocalRepository {
	t.Helper()
	repo, err := data.NewLocalRepositoryWithPath(filepath.Join(t.TempDir(), "todos.json"))
	if err != nil {
		t.Fatalf("NewLocalRepositoryWithPath() error = %v", err)
	}
	return repo
}

func TestAddTodosCreatesItemsAndPrintsResult(t *testing.T) {
	repo := newProcessorTestRepository(t)
	var out bytes.Buffer

	if err := addTodos(repo, &out, []string{"write tests", "ship release"}, true); err != nil {
		t.Fatalf("addTodos() error = %v", err)
	}

	todos, err := repo.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos() error = %v", err)
	}
	if len(todos) != 2 || !todos[0].Done || !todos[1].Done {
		t.Fatalf("unexpected todos: %+v", todos)
	}
	output := out.String()
	if !strings.Contains(output, "write tests") || !strings.Contains(output, "ship release") {
		t.Fatalf("output does not include added todos: %s", output)
	}
}

func TestListTodosFilters(t *testing.T) {
	repo := newProcessorTestRepository(t)
	if _, err := repo.CreateAndAddTodo("Write README", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	if _, err := repo.CreateAndAddTodo("ship release", true); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	if _, err := repo.CreateAndAddTodo("write tests", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}

	t.Run("done", func(t *testing.T) {
		var out bytes.Buffer
		if err := listTodos(repo, &out, listOptions{hasDone: true, done: true, begin: 0, end: -1}); err != nil {
			t.Fatalf("listTodos(done) error = %v", err)
		}
		output := out.String()
		if !strings.Contains(output, "ship release") || strings.Contains(output, "Write README") || strings.Contains(output, "write tests") {
			t.Fatalf("unexpected done filter output: %s", output)
		}
	})

	t.Run("content ignore case", func(t *testing.T) {
		var out bytes.Buffer
		if err := listTodos(repo, &out, listOptions{filter: "write", ignoreCase: true, begin: 0, end: -1}); err != nil {
			t.Fatalf("listTodos(filter) error = %v", err)
		}
		output := out.String()
		if !strings.Contains(output, "Write README") || !strings.Contains(output, "write tests") || strings.Contains(output, "ship release") {
			t.Fatalf("unexpected content filter output: %s", output)
		}
	})

	t.Run("id range", func(t *testing.T) {
		var out bytes.Buffer
		if err := listTodos(repo, &out, listOptions{begin: 1, end: 1}); err != nil {
			t.Fatalf("listTodos(range) error = %v", err)
		}
		output := out.String()
		if !strings.Contains(output, "ship release") || strings.Contains(output, "Write README") || strings.Contains(output, "write tests") {
			t.Fatalf("unexpected id range output: %s", output)
		}
	})
}

func TestListTodosShowsEmptyTable(t *testing.T) {
	repo := newProcessorTestRepository(t)
	var out bytes.Buffer

	if err := listTodos(repo, &out, listOptions{begin: 0, end: -1}); err != nil {
		t.Fatalf("listTodos() error = %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "ID") || !strings.Contains(output, "CONTENT") {
		t.Fatalf("empty list should render table headers: %s", output)
	}
}

func TestListTodosFiltersCaseSensitive(t *testing.T) {
	repo := newProcessorTestRepository(t)
	if _, err := repo.CreateAndAddTodo("Write README", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	if _, err := repo.CreateAndAddTodo("write tests", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}

	var out bytes.Buffer
	if err := listTodos(repo, &out, listOptions{filter: "write", ignoreCase: false, begin: 0, end: -1}); err != nil {
		t.Fatalf("listTodos() error = %v", err)
	}
	output := out.String()
	if strings.Contains(output, "Write README") || !strings.Contains(output, "write tests") {
		t.Fatalf("unexpected case-sensitive output: %s", output)
	}
}

func TestModifyTodoUpdatesFields(t *testing.T) {
	repo := newProcessorTestRepository(t)
	created, err := repo.CreateAndAddTodo("base", false)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := modifyTodo(repo, &out, modifyOptions{index: created.ID, content: " updated", append: true}); err != nil {
		t.Fatalf("modifyTodo(append) error = %v", err)
	}
	got, err := repo.GetTodoById(created.ID)
	if err != nil {
		t.Fatalf("GetTodoById() error = %v", err)
	}
	if got.Content != "base updated" {
		t.Fatalf("Content after append = %q", got.Content)
	}

	if err := modifyTodo(repo, &out, modifyOptions{index: created.ID, content: "prefix ", insert: true}); err != nil {
		t.Fatalf("modifyTodo(insert) error = %v", err)
	}
	got, _ = repo.GetTodoById(created.ID)
	if got.Content != "prefix base updated" {
		t.Fatalf("Content after insert = %q", got.Content)
	}

	if err := modifyTodo(repo, &out, modifyOptions{index: created.ID, content: "final", hasDone: true, done: true, hasPriority: true, priority: 4}); err != nil {
		t.Fatalf("modifyTodo(update) error = %v", err)
	}
	got, _ = repo.GetTodoById(created.ID)
	if got.Content != "final" || !got.Done || got.FinishTime == nil || got.Priority == nil || *got.Priority != 4 {
		t.Fatalf("unexpected modified todo: %+v", got)
	}
	if _, err := time.Parse(time.RFC3339, *got.FinishTime); err != nil {
		t.Fatalf("FinishTime is not RFC3339: %v", err)
	}
	if !strings.Contains(out.String(), "modified") || !strings.Contains(out.String(), "final") {
		t.Fatalf("output does not include modified result: %s", out.String())
	}
}

func TestModifyTodoCanMarkUndone(t *testing.T) {
	repo := newProcessorTestRepository(t)
	created, err := repo.CreateAndAddTodo("done", true)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := modifyTodo(repo, &out, modifyOptions{index: created.ID, hasDone: true, done: false}); err != nil {
		t.Fatalf("modifyTodo() error = %v", err)
	}

	got, err := repo.GetTodoById(created.ID)
	if err != nil {
		t.Fatalf("GetTodoById() error = %v", err)
	}
	if got.Done || got.FinishTime != nil {
		t.Fatalf("todo should be undone with no finish time: %+v", got)
	}
	if !strings.Contains(out.String(), ui.FLAG_NOT_DONE) {
		t.Fatalf("output should show undone status: %s", out.String())
	}
}

func TestModifyTodoMissingID(t *testing.T) {
	repo := newProcessorTestRepository(t)
	var out bytes.Buffer
	if err := modifyTodo(repo, &out, modifyOptions{index: 999, content: "missing"}); err == nil {
		t.Fatal("modifyTodo() expected error for missing id")
	}
}

func TestModifyTodoKeepsContentWhenEmpty(t *testing.T) {
	repo := newProcessorTestRepository(t)
	created, err := repo.CreateAndAddTodo("keep", false)
	if err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := modifyTodo(repo, &out, modifyOptions{index: created.ID, hasPriority: true, priority: 2}); err != nil {
		t.Fatalf("modifyTodo() error = %v", err)
	}
	got, err := repo.GetTodoById(created.ID)
	if err != nil {
		t.Fatalf("GetTodoById() error = %v", err)
	}
	if got.Content != "keep" || got.Priority == nil || *got.Priority != 2 {
		t.Fatalf("unexpected todo after empty content modify: %+v", got)
	}
}

func TestRemoveTodosDeletesAndPrintsTables(t *testing.T) {
	repo := newProcessorTestRepository(t)
	for _, content := range []string{"zero", "one", "two"} {
		if _, err := repo.CreateAndAddTodo(content, false); err != nil {
			t.Fatalf("CreateAndAddTodo(%q) error = %v", content, err)
		}
	}
	var out bytes.Buffer

	if err := removeTodos(repo, &out, []int{1}); err != nil {
		t.Fatalf("removeTodos() error = %v", err)
	}

	todos, err := repo.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos() error = %v", err)
	}
	if len(todos) != 2 || todos[0].Content != "zero" || todos[1].ID != 1 || todos[1].Content != "two" {
		t.Fatalf("unexpected remaining todos: %+v", todos)
	}
	output := out.String()
	if !strings.Contains(output, "removed") || !strings.Contains(output, "one") || !strings.Contains(output, "last") || strings.Contains(output, "| 1  | ❌   |          | one") {
		t.Fatalf("unexpected remove output: %s", output)
	}
}

func TestRemoveTodosMissingID(t *testing.T) {
	repo := newProcessorTestRepository(t)
	if _, err := repo.CreateAndAddTodo("keep", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := removeTodos(repo, &out, []int{999}); err != nil {
		t.Fatalf("removeTodos() error = %v", err)
	}
	if repo.Size() != 1 || !strings.Contains(out.String(), "keep") {
		t.Fatalf("unexpected missing remove result: size=%d output=%s", repo.Size(), out.String())
	}
}

func TestRemoveTodosDeletesAll(t *testing.T) {
	repo := newProcessorTestRepository(t)
	if _, err := repo.CreateAndAddTodo("only", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := removeTodos(repo, &out, []int{0}); err != nil {
		t.Fatalf("removeTodos() error = %v", err)
	}
	if repo.Size() != 0 {
		t.Fatalf("Size() = %d, want 0", repo.Size())
	}
	output := out.String()
	if !strings.Contains(output, "removed") || !strings.Contains(output, "only") || strings.Contains(output, "last") {
		t.Fatalf("unexpected delete-all output: %s", output)
	}
}

func TestClearTodosRemovesAllItems(t *testing.T) {
	repo := newProcessorTestRepository(t)
	for _, content := range []string{"first", "second"} {
		if _, err := repo.CreateAndAddTodo(content, false); err != nil {
			t.Fatalf("CreateAndAddTodo(%q) error = %v", content, err)
		}
	}
	var out bytes.Buffer

	if err := clearTodos(repo, &out, strings.NewReader(""), true); err != nil {
		t.Fatalf("clearTodos() error = %v", err)
	}
	if repo.Size() != 0 {
		t.Fatalf("Size() = %d, want 0", repo.Size())
	}
	output := out.String()
	if !strings.Contains(output, "cleared") || !strings.Contains(output, "first") || !strings.Contains(output, "second") {
		t.Fatalf("unexpected clear output: %s", output)
	}
}

func TestClearTodosEmptyRepository(t *testing.T) {
	repo := newProcessorTestRepository(t)
	var out bytes.Buffer

	if err := clearTodos(repo, &out, strings.NewReader(""), true); err != nil {
		t.Fatalf("clearTodos() error = %v", err)
	}
	if repo.Size() != 0 || !strings.Contains(out.String(), "cleared") {
		t.Fatalf("unexpected empty clear result: size=%d output=%s", repo.Size(), out.String())
	}
}

func TestClearTodosCancelsWithoutConfirmation(t *testing.T) {
	repo := newProcessorTestRepository(t)
	if _, err := repo.CreateAndAddTodo("keep", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := clearTodos(repo, &out, strings.NewReader("n\n"), false); err != nil {
		t.Fatalf("clearTodos() error = %v", err)
	}
	if repo.Size() != 1 || !strings.Contains(out.String(), "clear cancelled") {
		t.Fatalf("unexpected cancelled clear result: size=%d output=%s", repo.Size(), out.String())
	}
}

func TestClearTodosConfirmsBeforeClear(t *testing.T) {
	repo := newProcessorTestRepository(t)
	if _, err := repo.CreateAndAddTodo("remove", false); err != nil {
		t.Fatalf("CreateAndAddTodo() error = %v", err)
	}
	var out bytes.Buffer

	if err := clearTodos(repo, &out, strings.NewReader("y\n"), false); err != nil {
		t.Fatalf("clearTodos() error = %v", err)
	}
	if repo.Size() != 0 || !strings.Contains(out.String(), "This will clear all todos") || !strings.Contains(out.String(), "cleared") {
		t.Fatalf("unexpected confirmed clear result: size=%d output=%s", repo.Size(), out.String())
	}
}
