package processor

import (
	"fmt"
	"io"
	"os"
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
	contents := make([]string, 0, len(ctx.Args()))
	for _, ct := range ctx.Args() {
		contents = append(contents, ct.ForceToString())
	}
	if err := addTodos(data.CreateRepository(), os.Stdout, contents, done); err != nil {
		fmt.Println(err)
	}

}

func addTodos(repository data.Repository, out io.Writer, contents []string, done bool) error {
	tb := newTodoTableWithTitle("added")
	failedTodoList := make([]string, 0, len(contents))

	for _, content := range contents {
		todo, err := repository.CreateAndAddTodo(content, done)
		if err != nil {
			fmt.Fprintln(out, err)
			failedTodoList = append(failedTodoList, fmt.Sprintf("add todo:`%v` faild.", content))
			continue
		}

		tb.AddTodo(todo)
	}

	if err := tb.ShowTo(out); err != nil {
		return err
	}
	if len(failedTodoList) > 0 {
		fmt.Fprintln(out, failedTodoList)
	}
	return nil
}
