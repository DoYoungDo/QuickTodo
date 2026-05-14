package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	TABLE_HEADER_KEY     = "Key"
	TABLE_HEADER_VALUE   = "Value"
	TABLE_HEADER_HISTORY = "History"
)

type ConfigTable struct {
	table.Writer
	showHistory bool
	displayMode string
}

func NewConfigTable() *ConfigTable {
	return NewConfigTableWithHistory(false)
}

func NewConfigTableWithHistory(showHistory bool) *ConfigTable {
	tb := newTable()
	if showHistory {
		tb.AppendHeader(table.Row{TABLE_HEADER_KEY, TABLE_HEADER_VALUE, TABLE_HEADER_HISTORY})
	} else {
		tb.AppendHeader(table.Row{TABLE_HEADER_KEY, TABLE_HEADER_VALUE})
	}
	return &ConfigTable{Writer: tb, showHistory: showHistory}
}

func (t *ConfigTable) AddConfig(key, value string) {
	t.AddConfigWithHistory(key, value, "")
}

func (t *ConfigTable) AddConfigWithHistory(key, value, history string) {
	if t.showHistory {
		t.AppendRow(table.Row{key, value, history})
		return
	}
	t.AppendRow(table.Row{key, value})
}

func (t *ConfigTable) Show() {
	_ = t.ShowTo(os.Stdout)
}

func (t *ConfigTable) SetDisplayMode(mode string) {
	t.displayMode = mode
}

func (t *ConfigTable) ShowTo(w io.Writer) error {
	if t.displayMode == DisplayModeMarkdown {
		_, err := fmt.Fprintln(w, t.RenderMarkdown())
		return err
	}
	_, err := fmt.Fprintln(w, t.Render())
	return err
}
