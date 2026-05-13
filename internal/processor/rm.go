package processor

import (
	"fmt"
	"todo_list/internal/data"
	"todo_list/internal/ui"

	cmd "github.com/DoYoungDo/commander-go"
)

func Remove(ctx *cmd.Context) {
	ids := ctx.Args()
	intIds := make([]int, 0, len(ids))
	for _, id := range ids {
		if id.IsInt() {
			intIds = append(intIds, id.ToInt())
		}
	}
	if len(intIds) == 0 {
		return
	}

	repository := data.CreateRepository()
	removedTodos, err := repository.RemoveTodos(intIds)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(removedTodos) > 0 {
		tb := ui.NewTodoTableWithTitle("removed")
		for _, td := range removedTodos {
			tb.AddTodo(td)
		}
		tb.Show()
	}

	if repository.Size() > 0 {
		lastTodos, err := repository.GetTodos()
		if err != nil {
			fmt.Println(err)
			return
		}
		tb := ui.NewTodoTableWithTitle("last")
		for _, td := range lastTodos {
			tb.AddTodo(td)
		}
		tb.Show()
	}
}
