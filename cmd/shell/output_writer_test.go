package shell

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gookit/color"
)

func readOutput(r io.Reader) (output string, err error) {
	var out []byte

	if out, err = ioutil.ReadAll(r); err != nil {
		return
	}

	output = strings.TrimSpace(string(out))
	return
}

func newTestingOutputWriter() (outputWriter OutputWriter, buf *bytes.Buffer) {
	outputWriter = NewOutputWriter()
	buf = bytes.NewBufferString("")
	outputWriter.SetWriter(buf)

	return
}

func TestNewOutputWriter(t *testing.T) {
	o := NewOutputWriter()

	if _, ok := o.(*DefaultOutputWriter); !ok {
		t.Errorf("NewOutputWriter() did not return a *DefaultOutputWriter")
	}
}

func TestGetWriterOutputWriter(t *testing.T) {
	o, b := newTestingOutputWriter()

	w := o.GetWriter()

	if w != b {
		t.Error("failed to get correct writer on GetWriter on OutputWriter")
	}
}

func TestPrintlnOutputWriter(t *testing.T) {
	o, b := newTestingOutputWriter()

	expected := "testing text"
	o.Println(expected)

	output, err := readOutput(b)

	if err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}

func TestPrintfOutputWriter(t *testing.T) {
	o, b := newTestingOutputWriter()

	expected := "testing text"
	o.Printf("testing %s", "text")

	output, err := readOutput(b)

	if err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}

func TestErrorOutputWriter(t *testing.T) {
	o, b := newTestingOutputWriter()

	o.Error(errors.New("testing error"))

	output, err := readOutput(b)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.BgRed, color.FgWhite).Sprint("error: testing error")

	if output != expected {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}

func TestWarningOutputWriter(t *testing.T) {
	o, b := newTestingOutputWriter()

	o.Warning("testing warning")

	output, err := readOutput(b)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.Yellow).Sprint("testing warning")

	if output != expected {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}

func TestSuccessOutputWriter(t *testing.T) {
	o, b := newTestingOutputWriter()

	o.Success("testing success")

	output, err := readOutput(b)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.Green).Sprint("testing success")

	if output != expected {
		t.Errorf("expecting output '%s', got '%s'", expected, output)
	}
}
