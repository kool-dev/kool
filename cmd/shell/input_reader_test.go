package shell

import (
	"bytes"
	"testing"
)

func newTestingInputReader() (reader InputReader, buf *bytes.Buffer) {
	reader = NewInputReader()
	buf = bytes.NewBufferString("")
	reader.SetReader(buf)
	return
}

func TestNewInputReader(t *testing.T) {
	i := NewInputReader()

	if _, ok := i.(*DefaultInputReader); !ok {
		t.Errorf("NewInputReader() did not return a *DefaultInputReader")
	}
}

func TestInputReader(t *testing.T) {
	i, b := newTestingInputReader()

	r := i.GetReader()

	if r != b {
		t.Error("failed to get correct reader on GetReader on InputReader")
	}
}
