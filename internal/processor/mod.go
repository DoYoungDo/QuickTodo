package processor

import (
	"errors"
	"fmt"
	"time"
	"todo_list/internal/data"
	"todo_list/internal/ui"

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
	oAppend, oInsert, oDone, hasPriority, oPriority := func() (bool, bool, bool, bool, int) {
		oA := false
		if opt := ctx.Opt("append"); !opt.IsEmpty() {
			oA = true
		}
		oI := false
		if opt := ctx.Opt("insert"); !opt.IsEmpty() {
			oI = true
		}
		oD := false
		if opt := ctx.Opt("done"); !opt.IsEmpty() {
			oD = true
		}
		hasP := false
		p := 0
		if opt := ctx.Opt("priority"); !opt.IsEmpty() && opt.IsInt() {
			hasP = true
			p = opt.ToInt()
		}

		return oA, oI, oD, hasP, p
	}()

	repository := data.NewLocalRepository()
	todo, err := repository.GetTodoById(index)
	if err != nil {
		fmt.Println(err)
		return
	}
	if oAppend {
		todo.Content += content
	} else if oInsert {
		todo.Content = content + todo.Content
	} else if content != "" {
		todo.Content = content
	}
	if oDone {
		todo.Done = oDone
		timeNow := time.Now().Format(time.RFC3339)
		todo.FinishTime = &timeNow
	}
	if hasPriority {
		todo.Priority = &oPriority
	}

	repository.ModifyTodo(todo.ID, todo)

	tb := ui.NewTodoTableWithTitle("modified")
	tb.AddTodo(todo)
	tb.Show()
}
