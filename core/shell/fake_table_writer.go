package shell

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FakeTableWriter mock table writer for testing
type FakeTableWriter struct {
	CalledSetWriter, CalledAppendHeader, CalledAppendRow, CalledRender bool
	Headers, Rows                                                      [][]interface{}
	TableOut                                                           string
}

// SetWriter fake SetWriter behavior
func (f *FakeTableWriter) SetWriter(w io.Writer) {
	f.CalledSetWriter = true
}

// AppendHeader fake AppendHeader behavior
func (f *FakeTableWriter) AppendHeader(columns ...interface{}) {
	f.CalledAppendHeader = true
	f.Headers = append(f.Headers, columns)
}

// AppendRow fake AppendRow behavior
func (f *FakeTableWriter) AppendRow(columns ...interface{}) {
	f.CalledAppendRow = true
	f.Rows = append(f.Rows, columns)
}

// Render fake Render behavior
func (f *FakeTableWriter) Render() {
	f.CalledRender = true

	for _, columns := range f.Headers {
		columnsStr := make([]string, len(columns))

		for i := range columns {
			columnsStr[i] = fmt.Sprintf("%v", columns[i])
		}

		f.TableOut = f.TableOut + fmt.Sprintln(strings.Join(columnsStr, " | "))
	}

	for _, columns := range f.Rows {
		columnsStr := make([]string, len(columns))

		for i := range columns {
			columnsStr[i] = fmt.Sprintf("%v", columns[i])
		}

		f.TableOut = f.TableOut + fmt.Sprintln(strings.Join(columnsStr, " | "))
	}
}

// SortBy fake SortBy behavior
func (f *FakeTableWriter) SortBy(column int) {
	sort.SliceStable(f.Rows, func(i, j int) bool {
		return f.Rows[i][column-1].(string) < f.Rows[j][column-1].(string)
	})
}
