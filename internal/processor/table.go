package processor

import (
	"todo_list/internal/setting"
	"todo_list/internal/ui"
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
