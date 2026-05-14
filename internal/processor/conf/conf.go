package conf

import (
	"fmt"
	"io"
	"os"
	"slices"
	"todo_list/internal/setting"
	"todo_list/internal/ui"

	cmd "github.com/DoYoungDo/commander-go"
)

func Conf(ctx *cmd.Context) {
	if err := ctx.Command().Parse([]string{"list"}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

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

func List(ctx *cmd.Context) {
	keys := make([]string, 0, len(ctx.Args()))
	for _, arg := range ctx.Args() {
		keys = append(keys, arg.ForceToString())
	}
	if err := listConfig(os.Stdout, keys); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}

func setConfig(out io.Writer, key, value string) error {
	st, err := setting.Get()
	if err != nil {
		return err
	}
	if err := st.Set(key, value); err != nil {
		return err
	}
	return showConfig(out, map[string]string{key: value}, []string{key})
}

func listConfig(out io.Writer, keys []string) error {
	st, err := setting.Get()
	if err != nil {
		return err
	}
	values := st.Values()
	if len(keys) == 0 {
		keys = make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		slices.Sort(keys)
	}
	return showConfig(out, values, keys)
}

func showConfig(out io.Writer, values map[string]string, keys []string) error {
	tb := ui.NewConfigTable()
	for _, key := range keys {
		tb.AddConfig(key, values[key])
	}
	return tb.ShowTo(out)
}
