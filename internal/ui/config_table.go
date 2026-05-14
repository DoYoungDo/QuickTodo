package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	TABLE_HEADER_KEY   = "Key"
	TABLE_HEADER_VALUE = "Value"
)

type ConfigTable struct {
	table.Writer
}

func NewConfigTable() *ConfigTable {
	tb := newTable()
	tb.AppendHeader(table.Row{TABLE_HEADER_KEY, TABLE_HEADER_VALUE})
	return &ConfigTable{Writer: tb}
}

func (t *ConfigTable) AddConfig(key, value string) {
	t.AppendRow(table.Row{key, value})
}

func (t *ConfigTable) Show() {
	_ = t.ShowTo(os.Stdout)
}

func (t *ConfigTable) ShowTo(w io.Writer) error {
	_, err := fmt.Fprintln(w, t.Render())
	return err
}
