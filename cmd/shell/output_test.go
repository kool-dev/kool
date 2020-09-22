package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gookit/color"
)

func TestNewOutputWriter(t *testing.T) {
	outputwriter := NewOutputWriter()

	if _, assert := outputwriter.(*DefaultOutputWriter); !assert {
		t.Errorf("NewOutputWriter() did not return a *DefaultOutputWriter")
	}
}

func TestGetWriter(t *testing.T) {
	outputwriter := NewOutputWriter()

	if outputwriter.GetWriter() != os.Stdout {
		t.Errorf("SetWriter() failed; expected %v, got %v", os.Stdout, outputwriter.GetWriter())
	}
}

func TestSetWriter(t *testing.T) {
	outputwriter := NewOutputWriter()
	_, w, err := os.Pipe()

	if err != nil {
		t.Fatal(err)
	}

	outputwriter.SetWriter(w)

	if outputwriter.GetWriter() != w {
		t.Errorf("SetWriter() failed; expected %v, got %v", w, outputwriter.GetWriter())
	}
}

func TestExecError(t *testing.T) {
	outputwriter := NewOutputWriter()

	var buf bytes.Buffer
	outputwriter.SetWriter(&buf)

	err := errors.New("error")
	outputwriter.ExecError("output", errors.New("error"))

	errorMessage := color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err)
	outputMessage := "Output: output"

	expected := fmt.Sprintf("%s\n%s", errorMessage, outputMessage)
	output := strings.TrimSpace(buf.String())

	if output != expected {
		t.Errorf("ExecError() failed; expected '%s', got '%s'", expected, output)
	}
}

func TestExecErrorWithoutError(t *testing.T) {
	outputwriter := NewOutputWriter()

	var buf bytes.Buffer
	outputwriter.SetWriter(&buf)

	outputwriter.ExecError("output", nil)

	expected := "Output: output"
	output := strings.TrimSpace(buf.String())

	if output != expected {
		t.Errorf("ExecError() failed; expected '%s', got '%s'", expected, output)
	}
}

func TestExecErrorWithoutMessage(t *testing.T) {
	outputwriter := NewOutputWriter()

	var buf bytes.Buffer
	outputwriter.SetWriter(&buf)

	err := errors.New("error")
	outputwriter.ExecError("", errors.New("error"))

	expected := color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err)

	output := strings.TrimSpace(buf.String())

	if output != expected {
		t.Errorf("ExecError() failed; expected '%s', got '%s'", expected, output)
	}
}

func TestWarning(t *testing.T) {
	outputwriter := NewOutputWriter()

	var buf bytes.Buffer
	outputwriter.SetWriter(&buf)

	outputwriter.Warning("This is a warning")

	expected := color.New(color.Yellow).Sprint("This is a warning")

	output := strings.TrimSpace(buf.String())

	if output != expected {
		t.Errorf("ExecError() failed; expected '%s', got '%s'", expected, output)
	}
}
