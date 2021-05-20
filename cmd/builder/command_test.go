package builder

import (
	"errors"
	"os"
	"reflect"
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

func TestCmd(t *testing.T) {
	cmd := NewCommand("echo", "x1", "x2")

	cmdString := cmd.Cmd()
	expected := "echo"

	if cmdString != expected {
		t.Errorf("Failed to get the command executable; expected '%s', got '%s'", expected, cmdString)
	}
}

func TestArgs(t *testing.T) {
	cmd := NewCommand("echo", "x1", "x2")

	cmdArgs := cmd.Args()
	expected := []string{"x1", "x2"}

	if !reflect.DeepEqual(cmdArgs, expected) {
		t.Errorf("Failed to get the command executable; expected '%s', got '%s'", expected, cmdArgs)
	}
}

func TestParse(t *testing.T) {
	line := "echo 'xxx'"

	cmd := NewCommand("")
	err := cmd.Parse(line)

	if err != nil {
		t.Errorf("failed to parse proper command line onto Command; error: %s", err)
		return
	}

	if len(cmd.args) != 1 || cmd.command != "echo" || cmd.args[0] != "xxx" {
		t.Errorf("ParseCommand failed; given %s got %v", line, cmd.String())
	}
}

func TestErrorParseCommand(t *testing.T) {
	originalSplitFn := splitFn

	splitFn = func(s string) ([]string, error) {
		return []string{}, errors.New("split error")
	}

	defer func() {
		splitFn = originalSplitFn
	}()

	_, err := ParseCommand("echo x1")

	if err == nil {
		t.Error("expecting error 'split error', got none")
	} else if err.Error() != "split error" {
		t.Errorf("expecting error 'split error', got %v", err)
	}
}

func TestCopy(t *testing.T) {
	c := NewCommand("cmd", "arg1")

	cp := c.Copy()

	if cp.Cmd() != c.Cmd() || len(cp.Args()) != len(c.Args()) {
		t.Error("bad copy - differente cmd/args")
	}

	countArgs := len(c.Args())
	cp.AppendArgs("arg2", "arg3")

	if len(c.Args()) != countArgs {
		t.Error("unintended change on original command by changing the copy")
	}
}
