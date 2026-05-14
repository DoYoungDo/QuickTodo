package processor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"todo_list/internal/data"

	cmd "github.com/DoYoungDo/commander-go"
)

func Modify(ctx *cmd.Context) {
	index, content, err := func() (int, string, error) {
		args := ctx.Args()
		if !args[0].IsInt() {
			return -1, "", errors.New("index must be int")
		}

		i := args[0].ToInt()
		t := ""
		if len(args) > 1 {
			t = args[1].ForceToString()
		}

		return i, t, nil
	}()
	if err != nil {
		fmt.Println(err)
		return
	}
	oAppend, oInsert, hasDone, oDone, hasPriority, oPriority := func() (bool, bool, bool, bool, bool, int) {
		oA := false
		if opt := ctx.Opt("append"); !opt.IsEmpty() {
			oA = true
		}
		oI := false
		if opt := ctx.Opt("insert"); !opt.IsEmpty() {
			oI = true
		}
		hd := false
		oD := true
		if opt := ctx.Opt("done"); !opt.IsEmpty() {
			hd = true
			if opt.IsBool() {
				oD = opt.ToBool()
			}
		}
		hasP := false
		p := 0
		if opt := ctx.Opt("priority"); !opt.IsEmpty() && opt.IsInt() {
			hasP = true
			p = opt.ToInt()
		}

		return oA, oI, hd, oD, hasP, p
	}()
	if err := modifyTodo(data.CreateRepository(), os.Stdout, modifyOptions{
		index:       index,
		content:     content,
		append:      oAppend,
		insert:      oInsert,
		hasDone:     hasDone,
		done:        oDone,
		hasPriority: hasPriority,
		priority:    oPriority,
	}); err != nil {
		fmt.Println(err)
	}
}

type modifyOptions struct {
	index       int
	content     string
	append      bool
	insert      bool
	hasDone     bool
	done        bool
	hasPriority bool
	priority    int
}

func modifyTodo(repository data.Repository, out io.Writer, opts modifyOptions) error {
	todo, err := repository.GetTodoById(opts.index)
	if err != nil {
		return err
	}
	if opts.append {
		todo.Content += opts.content
	} else if opts.insert {
		todo.Content = opts.content + todo.Content
	} else if opts.content != "" {
		todo.Content = opts.content
	}
	if opts.hasDone {
		todo.Done = opts.done
	}
	if opts.hasPriority {
		todo.Priority = &opts.priority
	}

	if err := repository.ModifyTodo(todo.ID, *todo); err != nil {
		return err
	}

	tb := newTodoTableWithTitle("modified")
	tb.AddTodo(todo)
	return tb.ShowTo(out)
}
