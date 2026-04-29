package main

import (
	//"github.com/jroimartin/gocui"
	"fmt"
	"os"
	"todo_list/internal/app"
	"todo_list/internal/processor"

	cmd "github.com/DoYoungDo/commander-go"
)

func main() {
	todo := cmd.New(app.APP_NAME).
		Version(app.APP_VERSION).
		Description("todo list")

	// add
	todo.Command("add", "添加 待办项").
		Arguments("<todo...>", "待办项", nil).
		Options("-d, --done", "添加时完成", false).
		Action(processor.Add)

	// rm
	todo.Command("rm", "删除 待办项").
		Arguments("<index...>", "索引序号", nil).
		Action(func(ctx *cmd.Context) {
			fmt.Println("rm todo:", ctx.Args()[0].ToString())
		})

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
		Arguments("[range]", "显示范围", nil).
		Options("-d, --done [done]", "只显示完成的", nil).
		Options("-c, --count", "显示数量", false).
		Action(func(ctx *cmd.Context) {
			fmt.Println("list todos")
		})

	// done
	todo.Command("done", "完成 待办项").
		Arguments("<index...>", "索引序号", nil).
		Action(func(ctx *cmd.Context) {
			fmt.Println("done todo:", ctx.Args()[0].ToString())
		})

	// mv
	todo.Command("mv", "移动 待办项").
		Arguments("<index>", "待移动的待办项索引序号", nil).
		Arguments("<distindex>", "目标索引序号", nil).
		Action(func(ctx *cmd.Context) {
			fmt.Println("mv todo:", ctx.Args()[0].ToString(), "->", ctx.Args()[0].ToString())
		})

	// find
	todo.Command("find", "查找 待办项").
		Arguments("<todo...>", "查找内容", nil).
		Options("-c, --caseSensitive", "区分大小写", false).
		Options("-s, --single", "匹配单个条件", false).
		Options("-d, --done [done]", "匹配完成待办", nil).
		Action(func(ctx *cmd.Context) {
			fmt.Println("find todo:", ctx.Args()[0].ToString())
		})

	// clear
	todo.Command("clear", "清空 待办项").
		Action(func(ctx *cmd.Context) {
			fmt.Println("clear all todos")
		})

	// 根命令 action：todo xxx 等同于 todo add xxx
	todo.Action(func(ctx *cmd.Context) {
		// todo := ctx.Args()[0]
		// if todo.IsString() && todo.ToString() != "" {
		// fmt.Println("add todo:", todo.ToString())
		// }
	})

	if err := todo.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
