package main

import (
	"fmt"
	"os"
	"todo_list/internal/app"
	"todo_list/internal/processor"
	confprocessor "todo_list/internal/processor/conf"

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
		Options("-a, --append", "在原内容上追加，append优先于insert", nil).
		Options("-i, --insert", "在原内容上头插，insert让步于append", nil).
		Options("-d, --done", "修改为完成状态", nil).
		Options("-p, --priority <priority>", "设置优先级，取值0-5", nil).
		Action(processor.Modify)

	// list
	todo.Command("list", "显示 待办项").
		Options("-d, --done", "只显示完成的", nil).
		Options("-f, --filter <matchs>", "匹配内容", nil).
		Options("-i, --ignoreCase", "忽略大小写", nil).
		Options("-b, --begin <begin>", "开始索引", nil).
		Options("-e, --end <end>", "结束索引", nil).
		Action(processor.List)

	// done
	todo.Command("done", "完成 待办项，等价于：mod <index> -d").
		Arguments("<index...>", "索引序号", nil).
		Action(func(ctx *cmd.Context) {
			for _, arg := range ctx.Args() {
				if err := todo.Parse([]string{"mod", arg.ForceToString(), "-d"}); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
				}
			}
		})

	// clear
	todo.Command("clear", "清空 待办项").
		Options("-f, --force", "不弹出确认，强制清空", nil).
		Action(processor.Clear)

	// conf
	conf := todo.Command("conf", "配置").
		Action(confprocessor.Conf)
	conf.Command("set", "设置配置").
		Arguments("<key>", "配置键", nil).
		Arguments("<value>", "配置值", nil).
		Action(confprocessor.Set)
	conf.Command("list", "列出配置").
		Arguments("[key...]", "配置键", nil).
		Action(confprocessor.List)

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
