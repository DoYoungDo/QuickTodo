package conf

import (
	"fmt"
	"os"

	cmd "github.com/DoYoungDo/commander-go"
)

func List(ctx *cmd.Context) {
	keys := make([]string, 0, len(ctx.Args()))
	for _, arg := range ctx.Args() {
		keys = append(keys, arg.ForceToString())
	}
	showHistory := false
	if val := ctx.Opt("history"); !val.IsEmpty() {
		showHistory = true
	}
	if err := listConfig(os.Stdout, keys, showHistory); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}
