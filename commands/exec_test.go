package commands

import (
	"bytes"
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/compose"
	"strings"
	"testing"
)

func newFakeKoolExec() *KoolExec {
	return &KoolExec{
		*(newDefaultKoolService().Fake()),
		&KoolExecFlags{[]string{}, false},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "exec"},
	}
}

func newFailedFakeKoolExec() *KoolExec {
	return &KoolExec{
		*(newDefaultKoolService().Fake()),
		&KoolExecFlags{[]string{}, false},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "exec", MockInteractiveError: errors.New("error exec")},
	}
}

func TestNewKoolExec(t *testing.T) {
	k := NewKoolExec()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolExec instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolExec instance")
	} else {
		if len(k.Flags.EnvVariables) > 0 {
			t.Errorf("bad default value for EnvVariables flag on default KoolExec instance")
		}

		if k.Flags.Detach {
			t.Errorf("bad default value for Detach flag on default KoolExec instance")
		}
	}

	if _, ok := k.composeExec.(*compose.DockerCompose); !ok {
		t.Errorf("unexpected compose.DockerCompose on default KoolExec instance")
	}
}

func TestNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if val, ok := f.shell.(*shell.FakeShell).CalledInteractive["exec"]; !ok || !val {
		t.Error("did not call Interactive on KoolExec.composeExec Command")
	}

	interactiveArgs, ok := f.shell.(*shell.FakeShell).ArgsInteractive["exec"]

	if !ok || len(interactiveArgs) != 2 || interactiveArgs[0] != "service" || interactiveArgs[1] != "command" {
		t.Error("bad arguments to Interactive on KoolExec.composeExec Command")
	}
}

func TestNoArgsNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetOut(bytes.NewBufferString(""))

	if err := cmd.Execute(); err == nil {
		t.Error("expecting no arguments error executing exec command")
	}
}

func TestKoolUserEnvNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	f.env.(*environment.FakeEnvStorage).Envs["KOOL_ASUSER"] = "user_testing"
	// mock /etc/passwd return with existing user
	f.composeExec.(*builder.FakeCommand).MockExecOut = "kool:x:user_testing"

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledExec["exec"] {
		t.Error("did not call Exec")
	}

	if !f.composeExec.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs on KoolExec.composeExec Command")
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] != "--user" || argsAppend[1] != "user_testing" {
		t.Errorf("bad arguments to KoolExec.composeExec Command with KOOL_USER environment variable: %v", argsAppend)
	}

	// now check
	f = newFakeKoolExec()
	cmd = NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	f.env.(*environment.FakeEnvStorage).Envs["KOOL_ASUSER"] = "user_testing"
	// mock /etc/passwd return without existing user
	f.composeExec.(*builder.FakeCommand).MockExecOut = ""

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	argsAppend = f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 0 {
		t.Errorf("unexpected args appended: %v", argsAppend)
	}
}

func TestEnvFlagNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"--env=VAR_TEST=1", "service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.composeExec.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs on KoolExec.composeExec Command")
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] != "--env" || argsAppend[1] != "VAR_TEST=1" {
		t.Errorf("bad arguments to KoolExec.composeExec Command with EnvVariables flag")
	}
}

func TestDetachFlagNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"--detach", "service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.composeExec.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs on KoolExec.composeExec Command")
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 1 || argsAppend[0] != "--detach" {
		t.Errorf("bad arguments to KoolExec.composeExec Command with Detach flag")
	}
}

func TestFailingNewExecCommand(t *testing.T) {
	f := newFailedFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	assertExecGotError(t, cmd, "error exec")
}

func TestDockerComposeTerminalAwarness(t *testing.T) {
	f := newFakeKoolExec()
	f.composeExec = compose.NewDockerCompose("cmd")
	f.composeExec.(*compose.DockerCompose).SetShell(&shell.FakeShell{})
	f.composeExec.(*compose.DockerCompose).SetLocalDockerCompose(&builder.FakeCommand{
		MockLookPathError: errors.New("some error"),
	})

	cmd := NewExecCommand(f)
	cmd.SetArgs([]string{"service", "command"})

	f.shell.(*shell.FakeShell).MockIsTerminal = false
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if strings.Contains(f.composeExec.String(), " -t ") {
		t.Errorf("unexpected -t flag when NOT under TTY; %s", f.composeExec.String())
	}

	f.shell.(*shell.FakeShell).MockIsTerminal = true
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !strings.Contains(f.composeExec.String(), " -t ") {
		t.Error("missing -t flag when under TTY")
	}
}

func TestNonTerminalNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()

	f.shell.(*shell.FakeShell).MockIsTerminal = false

	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 1 || argsAppend[0] != "-T" {
		t.Errorf("bad arguments to KoolExec.composeExec Command on non terminal environment")
	}
}
