package shell

import (
	"bytes"
	"io"
	"kool-dev/kool/cmd/builder"
	"os"
	"strings"
	"testing"
)

func TestNewShell(t *testing.T) {
	s := NewShell()

	if _, ok := s.(*DefaultShell); !ok {
		t.Error("NewShell() did not return a *DefaultShell")
	}

	if s.InStream() != os.Stdin {
		t.Error("NewShell() should initialize input stream with os.Stdin")
	}

	if s.OutStream() != os.Stdout {
		t.Error("NewShell() should initialize output stream with os.Stdout")
	}

	if s.ErrStream() != os.Stderr {
		t.Error("NewShell() should initialize error stream with os.Stderr")
	}
}

func TestSetInStreamDefaultShell(t *testing.T) {
	s := NewShell()

	r := bytes.NewReader([]byte{})
	s.SetInStream(r)

	if s.InStream() != r {
		t.Error("failed calling SetInStream() on *DefaultShell")
	}
}

func TestSetOutStreamDefaultShell(t *testing.T) {
	s := NewShell()

	w := bytes.NewBufferString("")
	s.SetOutStream(w)

	if s.OutStream() != w {
		t.Error("failed calling SetOutStream() on *DefaultShell")
	}
}

func TestSetErrStreamDefaultShell(t *testing.T) {
	s := NewShell()

	w := bytes.NewBufferString("")
	s.SetErrStream(w)

	if s.ErrStream() != w {
		t.Error("failed calling SetErrStream() on *DefaultShell")
	}
}

func TestExecDefaultShell(t *testing.T) {
	s := NewShell()
	command := builder.NewCommand("echo", "x")

	output, err := s.Exec(command)

	if err != nil {
		t.Errorf("unexpected error calling Exec on *DefaultShell, err: %v", err)
	}

	output = strings.TrimSpace(output)

	if output != "x" {
		t.Errorf("Exec failed; expected output 'x', got '%s'", output)
	}
}

func TestExec(t *testing.T) {
	output, err := Exec("echo", "x")

	if err != nil {
		t.Errorf("unexpected error calling Exec on *DefaultShell, err: %v", err)
	}

	output = strings.TrimSpace(output)

	if output != "x" {
		t.Errorf("Exec failed; expected output 'x', got '%s'", output)
	}
}

func TestInteractiveDefaultShell(t *testing.T) {
	s := NewShell()
	command := builder.NewCommand("echo", "x")

	r, w, _ := os.Pipe()
	s.SetOutStream(w)

	err := s.Interactive(command)

	w.Close()

	if err != nil {
		t.Errorf("Interactive failed on *DefaultShell; expected no errors 'x', got '%v'", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	output := strings.TrimSpace(buf.String())

	if output != "x" {
		t.Errorf("Interactive failed on *DefaultShell; expected output 'x', got '%s'", output)
	}
}

func TestInteractive(t *testing.T) {
	r, w, _ := os.Pipe()

	originalStdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout
	}()

	err := Interactive("echo", "x")

	w.Close()

	if err != nil {
		t.Errorf("Interactive failed on *DefaultShell; expected no errors 'x', got '%v'", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	output := strings.TrimSpace(buf.String())

	if output != "x" {
		t.Errorf("Interactive failed on *DefaultShell; expected output 'x', got '%s'", output)
	}
}
