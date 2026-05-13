package main

import (
	"fmt"
	"os"
	"todo_list/internal/app"
	"todo_list/internal/processor"

	cmd "github.com/DoYoungDo/commander-go"
)

func main() {
	todo := cmd.New(app.APP_NAME).
		Version(app.APP_VERSION).
		Description("todo list").
		Arguments("[todo...]", "待办项", nil)

	// add
	todo.Command("add", "添加 待办项").
		Arguments("<todo...>", "待办项", nil).
		Options("-d, --done", "添加时完成", nil).
		Action(processor.Add)

	// rm
	todo.Command("rm", "删除 待办项").
		Arguments("<index...>", "索引序号", nil).
		Action(processor.Remove)

	// mod
	todo.Command("mod", "修改 待办项").
		Arguments("<index>", "索引序号", nil).
		Arguments("[todo]", "待办内容", nil).
		Options("-a, --todoend", "在原内容上追加", false).
		Options("-i, --insert", "在原内容上头插", false).
		Options("-d, --done [done]", "修改为完成", nil).
		Options("-p, --priority <priority>", "设置优先级，取值1-5", nil).
		Action(func(ctx *cmd.Context) {
			fmt.Println("mod todo index:", ctx.Args()[0].ToString())
		})

	// list
	todo.Command("list", "显示 待办项").
		Options("-d, --done", "只显示完成的", nil).
		Options("-f, --filter <matchs>", "匹配内容", nil).
		Options("-i, --ignoreCase", "忽略大小写", nil).
		Options("-b, --begin <begin>", "开始索引", nil).
		Options("-e, --end <end>", "结束索引", nil).
		Action(processor.List)

	// done
	todo.Command("done", "完成 待办项").
		Arguments("<index...>", "索引序号", nil).
		Action(func(ctx *cmd.Context) {
			fmt.Println("done todo:", ctx.Args()[0].ToString())
		})

	// clear
	todo.Command("clear", "清空 待办项").
		Action(func(ctx *cmd.Context) {
			fmt.Println("clear all todos")
		})

	// 根命令 action：todo xxx 等同于 todo add xxx
	todo.Action(func(ctx *cmd.Context) {
		argsSize := len(ctx.Args())
		if argsSize == 0 {
			todo.Parse([]string{"list"})
		} else {
			args := []string{"add"}
			for _, arg := range ctx.Args() {
				args = append(args, arg.ForceToString())
			}
			todo.Parse(args)
		}
	})

	if err := todo.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
