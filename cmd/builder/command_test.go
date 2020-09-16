package builder

import (
	"os"
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
