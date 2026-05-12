package processor

import (
	"fmt"
	"todo_list/internal/data"
	"todo_list/internal/ui"

	cmd "github.com/DoYoungDo/commander-go"
)

func Add(ctx *cmd.Context) {
	done := func() bool {
		if val := ctx.Opt("done"); !val.IsEmpty() && val.IsBool() && val.ToBool() {
			return true
		}
		return false
	}()

	tb := ui.NewTodoTable()

	repository := data.CreateRepository()
	failedTodoList := make([]string, 0, len(ctx.Args()))

	for _, ct := range ctx.Args() {
		todo, err := repository.CreateAndAddTodo(ct.ForceToString(), done)
		if err != nil {
			fmt.Println(err)
			failedTodoList = append(failedTodoList, fmt.Sprintf("add todo:`%v` faild.", ct.ForceToString()))
			continue
		}

		tb.AddTodo(todo)
	}

	tb.Show()
	if len(failedTodoList) > 0 {
		fmt.Println(failedTodoList)
	}
}
