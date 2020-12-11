package shell

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestFakeTableWriter(t *testing.T) {
	f := &FakeTableWriter{}

	f.SetWriter(ioutil.Discard)

	if !f.CalledSetWriter {
		t.Errorf("failed to mock method SetWriter on FakeTableWriter")
	}

	f.AppendHeader("header")

	if !f.CalledAppendHeader || len(f.Headers) != 1 || len(f.Headers[0]) != 1 || f.Headers[0][0] != "header" {
		t.Errorf("failed to mock method AppendHeader on FakeTableWriter")
	}

	f.AppendRow("row")

	if !f.CalledAppendRow || len(f.Rows) != 1 || len(f.Rows[0]) != 1 || f.Rows[0][0] != "row" {
		t.Errorf("failed to mock method AppendRow on FakeTableWriter")
	}

	f.Render()

	expected := `header
row`

	if !f.CalledRender || strings.TrimSpace(expected) != strings.TrimSpace(f.TableOut) {
		t.Errorf("failed to mock method Render on FakeTableWriter")
	}
}

func TestSortByFakeTableWriter(t *testing.T) {
	f := &FakeTableWriter{}

	f.SetWriter(ioutil.Discard)

	f.AppendRow("zRow")
	f.AppendRow("aRow")

	f.SortBy(1)

	if f.Rows[0][0] != "aRow" || f.Rows[1][0] != "zRow" {
		t.Errorf("failed to mock method SortBy on FakeTableWriter")
	}
}
