package processor

import (
	"fmt"
	"todo_list/internal/data"

	cmd "github.com/DoYoungDo/commander-go"
)

func Add(ctx *cmd.Context) {
	done := func() bool {
		if val := ctx.Opt("done"); !val.IsEmpty() && val.IsBool() && val.ToBool() {
			return true
		}
		return false
	}()

	repository := data.CreateRepository()

	todoList := make([]*data.Todo, 0, len(ctx.Args()))
	failedTodoList := make([]string, 0, len(ctx.Args()))
	for _, ct := range ctx.Args() {
		todo, err := repository.CreateAndAddTodo(ct.ForceToString(), done)
		if err != nil {
			fmt.Println(err)
			failedTodoList = append(failedTodoList, fmt.Sprintf("add todo:`%v` faild.", ct.ForceToString()))
			continue
		}
		todoList = append(todoList, todo)
	}

	fmt.Println(todoList, failedTodoList)
}
