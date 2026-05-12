package ui

import (
	"fmt"
	"strings"
	"time"
	"todo_list/internal/data"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	FLAG_DONE     = "✅"
	FLAG_NOT_DONE = "❌"
	FLAG_LEVEL    = "⭐"

	TABLE_HEADER_ID          = "ID"
	TABLE_HEADER_DONE        = "Done"
	TABLE_HEADER_LEVEL       = "Level"
	TABLE_HEADER_CONTENT     = "Content"
	TABLE_HEADER_CREATE_TIME = "CreateTime"
	TABLE_HEADER_FINISH_TIME = "FinishTime"
)

type ToDoTable struct {
	table.Writer

	todoList []*data.Todo
	filterBy []table.FilterBy
}

func newTable() table.Writer {
	tb := table.NewWriter()
	myStyle := table.StyleRounded
	myStyle.Title.Align = text.AlignCenter
	myStyle.Format.HeaderAlign = text.AlignCenter
	myStyle.Format.HeaderVAlign = text.VAlignMiddle
	myStyle.Options.DrawBorder = true
	tb.SetStyle(myStyle)
	return tb
}

func NewTodoTable() *ToDoTable {
	return NewTodoTableWithTitle("")
}
func NewTodoTableWithTitle(title string) *ToDoTable {
	tb := newTable()
	tb.SetTitle(title)
	tb.AppendHeader(table.Row{TABLE_HEADER_ID, TABLE_HEADER_DONE, TABLE_HEADER_LEVEL, TABLE_HEADER_CONTENT, TABLE_HEADER_CREATE_TIME, TABLE_HEADER_FINISH_TIME})
	return &ToDoTable{Writer: tb}
}
func (t *ToDoTable) AddTodo(todo *data.Todo) {
	t.todoList = append(t.todoList, todo)

	t.AppendRow(table.Row{fmt.Sprintf("%v", todo.ID),
		func() string {
			if todo.Done {
				return FLAG_DONE
			}
			return FLAG_NOT_DONE
		}(),
		func() string {
			if todo.Level == nil {
				return ""
			}
			return strings.Repeat(FLAG_LEVEL, *todo.Level)
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
func (t *ToDoTable) FilterID(begin, end int) {
	t.filterBy = append(t.filterBy, table.FilterBy{Name: TABLE_HEADER_ID, Operator: table.GreaterThanOrEqual, Value: fmt.Sprintf("%v", begin)})
	t.filterBy = append(t.filterBy, table.FilterBy{Name: TABLE_HEADER_ID, Operator: table.LessThanOrEqual, Value: fmt.Sprintf("%v", end)})
}
func (t *ToDoTable) FilterDone(done bool) {
	t.filterBy = append(t.filterBy, table.FilterBy{Name: TABLE_HEADER_DONE, Operator: table.Equal, CustomFilter: func(cellValue string) bool {
		return done && cellValue == FLAG_DONE || !done && cellValue == FLAG_NOT_DONE
	}})
}
func (t *ToDoTable) FilterContent(content string, ignoreCase bool) {
	t.filterBy = append(t.filterBy, table.FilterBy{Name: TABLE_HEADER_CONTENT, Operator: table.Contains, Value: content, IgnoreCase: ignoreCase})
}
func (t *ToDoTable) Show() {
	t.FilterBy(t.filterBy)
	fmt.Println(t.Render())
}
