package conf

import (
	"fmt"
	"os"

	cmd "github.com/DoYoungDo/commander-go"
)

func Delete(ctx *cmd.Context) {
	keys := make([]string, 0, len(ctx.Args()))
	for _, arg := range ctx.Args() {
		keys = append(keys, arg.ForceToString())
	}
	if len(keys) == 0 {
		fmt.Fprintln(os.Stdout, "conf del requires key")
		return
	}
	if err := deleteConfig(os.Stdout, keys); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}
