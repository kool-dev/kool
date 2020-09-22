package builder

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	exe, arg := "echo", "xxx"

	cmd := NewCommand(exe, arg)

	if len(cmd.args) != 1 || cmd.command != "echo" || cmd.args[0] != "xxx" {
		t.Errorf("NewCommand failed; given 'echo xxx' got %v", cmd.String())
	}
}

func TestParseCommand(t *testing.T) {
	line := "echo 'xxx'"
	cmd, err := ParseCommand(line)

	if err != nil {
		t.Errorf("failed to parse proper command line onto Command; error: %s", err)
		return
	}

	if len(cmd.args) != 1 || cmd.command != "echo" || cmd.args[0] != "xxx" {
		t.Errorf("ParseCommand failed; given %s got %v", line, cmd.String())
	}
}

func TestParseCommandWithEnvironmentVariable(t *testing.T) {
	env := "arbitraty-value-to-environment-var"
	os.Setenv("ENVVAR", env)
	line := "exec --opt $ENVVAR"
	cmd, err := ParseCommand(line)

	if err != nil {
		t.Errorf("failed parsing proper command with env var; error: %s", err)
		return
	}

	if len(cmd.args) != 2 || cmd.args[1] != env {
		t.Errorf("ParseCommand failed to parse environment variable; given %s got %v", line, cmd.String())
	}
}

func TestParseCommandWithUndefinedEnvironmentVariable(t *testing.T) {
	line := "exec --opt VAR=$NON_EXISTING_ENVVAR"
	cmd, err := ParseCommand(line)

	if err != nil {
		t.Errorf("failed parsing proper command line; error: %s", err)
		return
	}

	if len(cmd.args) != 2 || cmd.args[1] != "VAR=" {
		t.Errorf("ParseCommand failed to parse unset environment variable; given %s got %v", line, cmd.String())
	}
}

func TestAppendArgs(t *testing.T) {
	cmd := NewCommand("echo", "x1")

	cmd.AppendArgs("x2")

	if len(cmd.args) != 2 || cmd.args[1] != "x2" {
		t.Errorf("AppendArgs failed to add a new argument.")
	}
}

func TestString(t *testing.T) {
	cmd := NewCommand("echo", "x1", "x2")

	cmdString := cmd.String()
	expected := "echo x1 x2"

	if cmdString != expected {
		t.Errorf("Failed to convert command to string; expected '%s', got '%s'", expected, cmdString)
	}
}

func TestLookPath(t *testing.T) {
	cmd := NewCommand("go", "version")

	if err := cmd.LookPath(); err != nil {
		t.Errorf("LookPath failed; expected no errors, got '%v'", err)
	}
}

func TestInvalidLookPath(t *testing.T) {
	cmd := NewCommand("fakeCommand", "version")

	if err := cmd.LookPath(); err == nil {
		t.Error("LookPath failed; expected an error, got none.")
	}
}

func TestExec(t *testing.T) {
	cmd := NewCommand("echo", "x")

	output, err := cmd.Exec()

	if err != nil {
		t.Fatal(err)
	}

	output = strings.TrimSpace(output)

	if output != "x" {
		t.Errorf("Exec failed; expected output 'x', got '%s'", output)
	}
}

func TestInteractive(t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Fatal(err)
	}

	originalOutput := os.Stdout
	os.Stdout = w

	defer func(originalOutput *os.File) {
		os.Stdout = originalOutput
	}(originalOutput)

	cmd := NewCommand("echo", "x")
	err = cmd.Interactive()

	w.Close()

	if err != nil {
		t.Errorf("Interactive failed; expected no errors 'x', got '%v'", err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)

	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(buf.String())

	if output != "x" {
		t.Errorf("Interactive failed; expected output 'x', got '%s'", output)
	}
}
