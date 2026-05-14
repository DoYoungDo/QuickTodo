package conf

import (
	"io"
	"slices"
	"todo_list/internal/setting"
	"todo_list/internal/ui"
)

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
