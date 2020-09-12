package cmd

import (
	"os"
	"testing"
)

func TestParseCustomCommand(t *testing.T) {
	line := "echo 'xxx'"
	cmd := parseCustomCommand(line)

	if len(cmd) != 2 || cmd[0] != "echo" || cmd[1] != "xxx" {
		t.Errorf("parseCustomCommand failed; given %s got %v", line, cmd)
	}
}

func TestParseCustomCommandWithEnvironmentVariable(t *testing.T) {
	env := "arbitraty-value-to-environment-var"
	os.Setenv("ENVVAR", env)
	line := "exec --opt $ENVVAR"
	cmd := parseCustomCommand(line)

	if len(cmd) != 3 || cmd[2] != env {
		t.Errorf("parseCustomCommand failed to parse environment variable; given %s got %v", line, cmd)
	}
}

func TestParseCustomCommandWithUndefinedEnvironmentVariable(t *testing.T) {
	line := "exec --opt VAR=$NON_EXISTING_ENVVAR"
	cmd := parseCustomCommand(line)

	if len(cmd) != 3 || cmd[2] != "VAR=" {
		t.Errorf("parseCustomCommand failed to parse unset environment variable; given %s got %v", line, cmd)
	}
}
