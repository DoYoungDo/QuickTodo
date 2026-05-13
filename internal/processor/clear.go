package processor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"todo_list/internal/data"
	"todo_list/internal/ui"

	cmd "github.com/DoYoungDo/commander-go"
)

func Clear(ctx *cmd.Context) {
	force := false
	if opt := ctx.Opt("force"); !opt.IsEmpty() {
		force = true
	}
	if err := clearTodos(data.CreateRepository(), os.Stdout, os.Stdin, force); err != nil {
		fmt.Println(err)
	}
}

func clearTodos(repository data.Repository, out io.Writer, in io.Reader, force bool) error {
	if !force {
		fmt.Fprint(out, "This will clear all todos. Continue? [y/N]: ")
		answer, err := bufio.NewReader(in).ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" {
			fmt.Fprintln(out, "clear cancelled")
			return nil
		}
	}

	clearedTodos, err := repository.ClearTodos()
	if err != nil {
		return err
	}
	tb := ui.NewTodoTableWithTitle("cleared")
	for _, todo := range clearedTodos {
		tb.AddTodo(todo)
	}
	return tb.ShowTo(out)
}
