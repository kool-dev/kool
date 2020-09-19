package builder

import (
	"bytes"
	"io/ioutil"
	"kool-dev/kool/cmd/shell"
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
	cmd := NewCommand("echo", "xxx")

	cmd.AppendArgs("xxxx")

	if len(cmd.args) != 2 {
		t.Errorf("Expected 2 arguments, got %v", len(cmd.args))
	}

	if cmd.args[1] != "xxxx" {
		t.Errorf("Appended 'xxxx', got %s", cmd.args[1])
	}
}

func TestString(t *testing.T) {
	cmd := NewCommand("echo", "xxx")

	cmdStr := cmd.String()

	if cmdStr != "echo xxx" {
		t.Errorf("Expected 'echo xxx', got %s", cmdStr)
	}
}

func TestExec(t *testing.T) {
	cmd := NewCommand("echo", "xxx")

	output, err := cmd.Exec()

	if err != nil {
		t.Fatal(err)
	}

	if output != "xxx" {
		t.Errorf("Expected 'xxx', got %s", output)
	}
}

func TestInteractive(t *testing.T) {
	commander := shell.NewCommander()
	b := bytes.NewBufferString("")
	commanderWriter := &shell.DefaultCommanderWriter{Writer: b}
	commander.SetOut(commanderWriter)

	cmd := &DefaultCommand{"echo", []string{"xxx"}, commander}

	var (
		err error
		out []byte
	)

	if err = cmd.Interactive(); err != nil {
		t.Fatal(err)
	}

	if out, err = ioutil.ReadAll(b); err != nil {
		t.Fatal(err)
	}

	output := strings.Trim(string(out), "\n")

	if output != "xxx" {
		t.Errorf("Expected 'xxx', got %s", output)
	}
}
