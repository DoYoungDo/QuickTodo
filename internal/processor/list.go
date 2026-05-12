package processor

import (
	"fmt"
	"todo_list/internal/data"
	"todo_list/internal/ui"

	cmd "github.com/DoYoungDo/commander-go"
)

func List(ctx *cmd.Context) {
	hasDone, done, filter, ignoreCase, begin, end := func() (bool, bool, string, bool, int, int) {
		hd, d, f, i, b, e := false, false, "", false, 0, -1
		if val := ctx.Opt("done"); !val.IsEmpty() && val.IsBool() && val.ToBool() {
			hd, d = true, true
		}
		if val := ctx.Opt("filter"); !val.IsEmpty() {
			f = val.ForceToString()
		}
		if val := ctx.Opt("ignoreCase"); !val.IsEmpty() && val.IsBool() && val.ToBool() {
			i = true
		}
		if val := ctx.Opt("begin"); !val.IsEmpty() && val.IsInt() {
			b = val.ToInt()
		}
		if val := ctx.Opt("end"); !val.IsEmpty() && val.IsInt() {
			e = val.ToInt()
		}
		return hd, d, f, i, b, e
	}()
	// fmt.Println(hasDone, done, filter, ignoreCase, begin, end)

	tb := ui.NewTodoTable()

	repository := data.CreateRepository()
	todos, err := repository.GetTodos()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(todos) == 0 {
		tb.Show()
		return
	}

	for _, todo := range todos {
		tb.AddTodo(todo)
	}

	if end == -1 {
		end = todos[len(todos)-1].ID
	}
	tb.FilterID(begin, end)
	if hasDone {
		tb.FilterDone(done)
	}
	if filter != "" {
		tb.FilterContent(filter, ignoreCase)
	}
	tb.Show()
}
