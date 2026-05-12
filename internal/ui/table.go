package ui

import (
	"fmt"
	"strings"
	"time"
	"todo_list/internal/data"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type ToDoTable struct {
	table.Writer

	todoList []*data.Todo
}

func NewTodoTable() *ToDoTable {
	return NewTodoTableWithTitle("")
}
func NewTodoTableWithTitle(title string) *ToDoTable {
	tb := newTable()
	tb.SetTitle(title)
	tb.AppendHeader(table.Row{"ID", "Done", "Level", "Content", "CreateTime", "FinishTime"})
	return &ToDoTable{Writer: tb}
}
func (t *ToDoTable) AddTodo(todo *data.Todo) {
	t.todoList = append(t.todoList, todo)

	t.AppendRow(table.Row{fmt.Sprintf("%v", todo.ID),
		func() string {
			if todo.Done {
				return "✅"
			}
			return "❌"
		}(),
		func() string {
			if todo.Level == nil {
				return ""
			}
			return strings.Repeat("⭐", *todo.Level)
		}(),
		todo.Content,
		func() string {
			tm, _ := time.Parse(time.RFC3339, todo.CreateTime)
			return tm.Format(time.DateTime)
		}(),
		func() string {
			if todo.FinishTime == nil {
				return ""
			}
			tm, _ := time.Parse(time.RFC3339, *todo.FinishTime)
			return tm.Format(time.DateTime)
		}()})
}
func (t *ToDoTable) Show() {
	fmt.Println(t.Render())
}

func newTable() table.Writer {
	tb := table.NewWriter()
	myStyle := table.StyleRounded
	myStyle.Title.Align = text.AlignCenter
	// myStyle.Format.Header = text.FormatTitle
	myStyle.Format.HeaderAlign = text.AlignCenter
	myStyle.Format.HeaderVAlign = text.VAlignMiddle
	myStyle.Options.DrawBorder = true
	tb.SetStyle(myStyle)
	return tb
}
