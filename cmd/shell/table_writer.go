package shell

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

// DefaultTableWriter holds table output writer
type DefaultTableWriter struct {
	w table.Writer
}

// TableWriter holds table output writer logic
type TableWriter interface {
	SetWriter(io.Writer)
	AppendHeader(...interface{})
	AppendRow(...interface{})
	Render()
}

// NewTableWriter creates a new table writer
func NewTableWriter() TableWriter {
	return &DefaultTableWriter{table.NewWriter()}
}

// SetWriter set table output writer
func (t *DefaultTableWriter) SetWriter(w io.Writer) {
	t.w.SetOutputMirror(w)
}

// AppendHeader append header columns to table
func (t *DefaultTableWriter) AppendHeader(columns ...interface{}) {
	t.w.AppendHeader(columns)
}

// AppendRow append row columns to table
func (t *DefaultTableWriter) AppendRow(columns ...interface{}) {
	t.w.AppendRow(columns)
}

// Render render the table
func (t *DefaultTableWriter) Render() {
	t.w.Render()
}
