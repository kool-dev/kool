package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gookit/color"
)

func TestNewOutputWriter(t *testing.T) {
	outputWriter := NewOutputWriter()

	if _, assert := outputWriter.(*DefaultOutputWriter); !assert {
		t.Error("NewOutputWriter() did not return a *OutputWriter")
	}
}

func TestSetWriter(t *testing.T) {
	outputWriter := NewOutputWriter()

	outputWriter.SetWriter(bytes.NewBufferString(""))
	writer := outputWriter.GetWriter()

	if _, assert := writer.(*bytes.Buffer); !assert {
		t.Error("SetWriter() did not set writer to *bytes.Buffer")
	}
}

func TestSuccess(t *testing.T) {
	outputWriter, b := newTestingOutputWriter()

	outputWriter.Success("Success")

	expected := color.New(color.Green).Sprint("Success")

	var (
		output string
		err    error
	)

	if output, err = readBufferContent(b); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestWarning(t *testing.T) {
	outputWriter, b := newTestingOutputWriter()

	outputWriter.Warning("Warning")

	expected := color.New(color.Yellow).Sprint("Warning")

	var (
		output string
		err    error
	)

	if output, err = readBufferContent(b); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestError(t *testing.T) {
	outputWriter, b := newTestingOutputWriter()

	outputWriter.Error("Error")

	expected := color.New(color.BgRed, color.FgWhite).Sprint("Error")

	var (
		output string
		err    error
	)

	if output, err = readBufferContent(b); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestExecError(t *testing.T) {
	outputWriter, b := newTestingOutputWriter()

	err := errors.New("Error")
	outputWriter.ExecError("Output", err)

	errorMessage := color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err)
	outputMessage := "Output: Output"

	expected := fmt.Sprintf("%s\n%s", errorMessage, outputMessage)

	var output string

	if output, err = readBufferContent(b); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestExecErrorWithoutMessage(t *testing.T) {
	outputWriter, b := newTestingOutputWriter()

	err := errors.New("Error")
	outputWriter.ExecError("", err)

	expected := color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err)

	var output string

	if output, err = readBufferContent(b); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestExecErrorWithoutError(t *testing.T) {
	outputWriter, b := newTestingOutputWriter()

	outputWriter.ExecError("Output", nil)

	expected := "Output: Output"

	var (
		output string
		err    error
	)

	if output, err = readBufferContent(b); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func newTestingOutputWriter() (outputWriter OutputWriter, b *bytes.Buffer) {
	outputWriter = NewOutputWriter()
	b = bytes.NewBufferString("")
	outputWriter.SetWriter(b)
	return
}

func readBufferContent(b *bytes.Buffer) (content string, err error) {
	var out []byte

	if out, err = ioutil.ReadAll(b); err != nil {
		return
	}

	content = strings.Trim(string(out), "\n")
	return
}
