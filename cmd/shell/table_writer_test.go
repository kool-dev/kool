package shell

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewTableWriter(t *testing.T) {
	tableWriter := NewTableWriter()

	if _, ok := tableWriter.(*DefaultTableWriter); !ok {
		t.Errorf("NewTableWriter() did not return a *DefaultTableWriter")
	}
}
func TestTableWriter(t *testing.T) {
	tableWriter := NewTableWriter()

	b := bytes.NewBufferString("")
	tableWriter.SetWriter(b)

	tableWriter.AppendHeader("header")

	tableWriter.AppendRow("row")

	tableWriter.Render()

	var (
		out []byte
		err error
	)

	if out, err = ioutil.ReadAll(b); err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(string(out))
	expected := `
+--------+
| HEADER |
+--------+
| row    |
+--------+
`
	expected = strings.TrimSpace(expected)

	if expected != output {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}

func TestSortByTableWriter(t *testing.T) {
	tableWriter := NewTableWriter()

	b := bytes.NewBufferString("")
	tableWriter.SetWriter(b)

	tableWriter.AppendHeader("header")

	tableWriter.AppendRow("zRow")
	tableWriter.AppendRow("aRow")

	tableWriter.SortBy(1)

	tableWriter.Render()

	var (
		out []byte
		err error
	)

	if out, err = ioutil.ReadAll(b); err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(string(out))
	expected := `
+--------+
| HEADER |
+--------+
| aRow   |
| zRow   |
+--------+
`
	expected = strings.TrimSpace(expected)

	if expected != output {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}
