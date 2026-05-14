package conf

import (
	"fmt"
	"os"

	cmd "github.com/DoYoungDo/commander-go"
)

func Set(ctx *cmd.Context) {
	args := ctx.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stdout, "conf set requires key and value")
		return
	}
	if err := setConfig(os.Stdout, args[0].ForceToString(), args[1].ForceToString()); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}
