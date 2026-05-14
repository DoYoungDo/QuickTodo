package conf

import (
	"io"
	"slices"
	"strings"
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
	return showConfig(out, map[string]string{key: value}, nil, []string{key})
}

func deleteConfig(out io.Writer, keys []string) error {
	st, err := setting.Get()
	if err != nil {
		return err
	}
	if err := st.Delete(keys...); err != nil {
		return err
	}
	return showConfig(out, st.Values(), nil, keys)
}

func listConfig(out io.Writer, keys []string, showHistory bool) error {
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
	if !showHistory {
		return showConfig(out, values, nil, keys)
	}
	history := map[string]string{}
	for _, key := range keys {
		history[key] = strings.Join(st.History(key), "\n")
	}
	return showConfig(out, values, history, keys)
}

func showConfig(out io.Writer, values map[string]string, history map[string]string, keys []string) error {
	tb := ui.NewConfigTableWithHistory(history != nil)
	for _, key := range keys {
		tb.AddConfigWithHistory(key, values[key], history[key])
	}
	return tb.ShowTo(out)
}
