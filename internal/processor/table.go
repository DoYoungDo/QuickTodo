package processor

import (
	"quick_todo/internal/setting"
	"quick_todo/internal/ui"
)

func displayTableMode() string {
	st, err := setting.Get()
	if err != nil {
		return ""
	}
	return st.Get(setting.KeyDisplayTableMode)
}

func newTodoTableWithTitle(title string) *ui.ToDoTable {
	tb := ui.NewTodoTableWithTitle(title)
	tb.SetDisplayMode(displayTableMode())
	return tb
}
