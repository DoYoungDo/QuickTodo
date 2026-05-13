package ui

import (
	"bytes"
	"strings"
	"testing"
	"todo_list/internal/data"
)

func TestToDoTableShowsStatusPriorityAndTime(t *testing.T) {
	finishTime := "2026-05-13T10:11:12Z"
	priority := 3
	tb := NewTodoTable()
	tb.AddTodo(&data.Todo{
		ID:         1,
		Content:    "done task",
		CreateTime: "2026-05-13T09:10:11Z",
		FinishTime: &finishTime,
		Priority:   &priority,
		Done:       true,
	})
	tb.AddTodo(&data.Todo{
		ID:         2,
		Content:    "open task",
		CreateTime: "2026-05-13T09:10:11Z",
		Done:       false,
	})

	var out bytes.Buffer
	if err := tb.ShowTo(&out); err != nil {
		t.Fatalf("ShowTo() error = %v", err)
	}
	output := out.String()
	for _, want := range []string{"done task", "open task", FLAG_DONE, FLAG_NOT_DONE, strings.Repeat(FLAG_PRIORITY, 3), "2026-05-13 09:10:11", "2026-05-13 10:11:12"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output does not contain %q: %s", want, output)
		}
	}
}

func TestToDoTableFiltersDoneAndUndone(t *testing.T) {
	tests := []struct {
		name string
		done bool
		want string
		drop string
	}{
		{name: "done", done: true, want: "done task", drop: "open task"},
		{name: "undone", done: false, want: "open task", drop: "done task"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTodoTable()
			tb.AddTodo(&data.Todo{ID: 0, Content: "done task", CreateTime: "2026-05-13T09:10:11Z", Done: true})
			tb.AddTodo(&data.Todo{ID: 1, Content: "open task", CreateTime: "2026-05-13T09:10:11Z", Done: false})
			tb.FilterDone(tt.done)

			var out bytes.Buffer
			if err := tb.ShowTo(&out); err != nil {
				t.Fatalf("ShowTo() error = %v", err)
			}
			output := out.String()
			if !strings.Contains(output, tt.want) || strings.Contains(output, tt.drop) {
				t.Fatalf("unexpected filter output: %s", output)
			}
		})
	}
}

func TestToDoTableFiltersContentCaseSensitivity(t *testing.T) {
	tests := []struct {
		name       string
		ignoreCase bool
		wantUpper  bool
	}{
		{name: "case sensitive", ignoreCase: false, wantUpper: false},
		{name: "ignore case", ignoreCase: true, wantUpper: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTodoTable()
			tb.AddTodo(&data.Todo{ID: 0, Content: "Write README", CreateTime: "2026-05-13T09:10:11Z"})
			tb.AddTodo(&data.Todo{ID: 1, Content: "write tests", CreateTime: "2026-05-13T09:10:11Z"})
			tb.FilterContent("write", tt.ignoreCase)

			var out bytes.Buffer
			if err := tb.ShowTo(&out); err != nil {
				t.Fatalf("ShowTo() error = %v", err)
			}
			output := out.String()
			if !strings.Contains(output, "write tests") {
				t.Fatalf("expected lowercase match: %s", output)
			}
			if strings.Contains(output, "Write README") != tt.wantUpper {
				t.Fatalf("unexpected uppercase match state: %s", output)
			}
		})
	}
}

func TestToDoTableFiltersID(t *testing.T) {
	tb := NewTodoTable()
	for _, todo := range []*data.Todo{
		{ID: 0, Content: "zero", CreateTime: "2026-05-13T09:10:11Z"},
		{ID: 1, Content: "one", CreateTime: "2026-05-13T09:10:11Z"},
		{ID: 2, Content: "two", CreateTime: "2026-05-13T09:10:11Z"},
	} {
		tb.AddTodo(todo)
	}
	tb.FilterID(1, 1)

	var out bytes.Buffer
	if err := tb.ShowTo(&out); err != nil {
		t.Fatalf("ShowTo() error = %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "one") || strings.Contains(output, "zero") || strings.Contains(output, "two") {
		t.Fatalf("unexpected id filter output: %s", output)
	}
}

func TestToDoTableShowsEmptyTable(t *testing.T) {
	tb := NewTodoTable()
	var out bytes.Buffer
	if err := tb.ShowTo(&out); err != nil {
		t.Fatalf("ShowTo() error = %v", err)
	}
	output := out.String()
	if !strings.Contains(output, strings.ToUpper(TABLE_HEADER_ID)) || !strings.Contains(output, strings.ToUpper(TABLE_HEADER_CONTENT)) {
		t.Fatalf("empty table should still render headers: %s", output)
	}
}
