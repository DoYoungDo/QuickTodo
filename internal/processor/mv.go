package processor

import (
	"fmt"
	"io"
	"os"
	"todo_list/internal/data"

	cmd "github.com/DoYoungDo/commander-go"
)

func Move(ctx *cmd.Context) {
	index, distIndex, err := func() (int, int, error) {
		args := ctx.Args()
		if !args[0].IsInt() {
			return -1, -1, fmt.Errorf("index must be a int number.")
		}
		if !args[1].IsInt() {
			return -1, -1, fmt.Errorf("distIndex must be a int number.")
		}
		return args[0].ToInt(), args[1].ToInt(), nil
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := MoveTodo(data.CreateRepository(), os.Stdout, index, distIndex); err != nil {
		fmt.Println(err)
		return
	}
}

func MoveTodo(repository data.Repository, out io.Writer, index, distIndex int) error {
	todo, err := repository.MoveTodo(index, distIndex)
	if err != nil {
		return err
	}

	tb := newTodoTableWithTitle("moved")
	tb.AddTodo(todo)

	return tb.ShowTo(out)
}
