package processor

import (
	"fmt"
	"io"
	"os"
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
		fmt.Println("rm requires todos")
		fmt.Println("Usage: rm <todo-id> [todo-id...]")
		return
	}
	if err := removeTodos(data.CreateRepository(), os.Stdout, intIds); err != nil {
		fmt.Println(err)
	}
}

func removeTodos(repository data.Repository, out io.Writer, intIds []int) error {
	removedTodos, err := repository.RemoveTodos(intIds)
	if err != nil {
		return err
	}

	if len(removedTodos) > 0 {
		tb := ui.NewTodoTableWithTitle("removed")
		for _, td := range removedTodos {
			tb.AddTodo(td)
		}
		if err := tb.ShowTo(out); err != nil {
			return err
		}
	}

	if repository.Size() > 0 {
		lastTodos, err := repository.GetTodos()
		if err != nil {
			return err
		}
		tb := ui.NewTodoTableWithTitle("last")
		for _, td := range lastTodos {
			tb.AddTodo(td)
		}
		if err := tb.ShowTo(out); err != nil {
			return err
		}
	}
	return nil
}
