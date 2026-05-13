package processor

import (
	"fmt"
	"io"
	"os"
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
	if err := listTodos(data.CreateRepository(), os.Stdout, listOptions{
		hasDone:    hasDone,
		done:       done,
		filter:     filter,
		ignoreCase: ignoreCase,
		begin:      begin,
		end:        end,
	}); err != nil {
		fmt.Println(err)
	}

}

type listOptions struct {
	hasDone    bool
	done       bool
	filter     string
	ignoreCase bool
	begin      int
	end        int
}

func listTodos(repository data.Repository, out io.Writer, opts listOptions) error {
	tb := ui.NewTodoTable()

	todos, err := repository.GetTodos()
	if err != nil {
		return err
	}
	if len(todos) == 0 {
		return tb.ShowTo(out)
	}

	for _, todo := range todos {
		tb.AddTodo(todo)
	}

	if opts.end == -1 {
		opts.end = todos[len(todos)-1].ID
	}
	tb.FilterID(opts.begin, opts.end)
	if opts.hasDone {
		tb.FilterDone(opts.done)
	}
	if opts.filter != "" {
		tb.FilterContent(opts.filter, opts.ignoreCase)
	}
	return tb.ShowTo(out)
}
