package commands

import (
	"errors"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/checker"
	"strings"
	"testing"
)

type workspaceStopShell struct {
	shell.FakeShell
	commands []string
}

func (s *workspaceStopShell) Exec(command builder.Command, extraArgs ...string) (string, error) {
	if strings.HasPrefix(command.String(), "docker ps --all") {
		return "example-workspace-task-b\nexample-workspace-task-a\n", nil
	}
	if strings.HasPrefix(command.String(), "docker inspect kool-proxy") {
		return "", errors.New("proxy is not running")
	}
	return "", nil
}

func (s *workspaceStopShell) Interactive(command builder.Command, extraArgs ...string) error {
	s.commands = append(s.commands, strings.TrimSpace(command.String()+" "+strings.Join(extraArgs, " ")))
	return nil
}

func newFakeKoolStop() *KoolStop {
	fs := &KoolStop{
		*(newDefaultKoolService().Fake()),
		&KoolStopFlags{false},
		&checker.FakeChecker{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&builder.FakeCommand{},
	}
	fs.shell.(*shell.FakeShell).MockErrStream = io.Discard
	fs.shell.(*shell.FakeShell).MockOutStream = io.Discard
	return fs
}

func TestNewKoolStop(t *testing.T) {
	k := NewKoolStop()

	if _, ok := k.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolStop instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolStop instance")
	} else if k.Flags.Purge {
		t.Errorf("bad default value for Purge flag on default KoolStop instance")
	}

	if _, ok := k.check.(*checker.DefaultChecker); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStop instance")
	}

	if _, ok := k.down.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected compose.DockerCompose on default KoolStop instance")
	}

	if k.down.String() != "docker compose down" {
		t.Errorf("unexpected compose.DockerCompose.Command.String() on default KoolStop instance down")
	}
}

func TestNewStopCommand(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command; error: %v", err)
	}

	if !f.check.(*checker.FakeChecker).CalledCheck {
		t.Errorf("did not call Check")
	}

	if len(f.down.(*builder.FakeCommand).ArgsAppend) > 1 {
		t.Errorf("did not expect to call 2 AppendArgs on KoolStop.down Command")
	}
}

func TestNewStopCommandWithArgument(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)
	cmd.SetArgs([]string{"a", "b"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command; error: %v", err)
	}

	if !f.rm.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("should have called AppendArgs on KoolStop.rm Command")
	}
	appended := f.rm.(*builder.FakeCommand).ArgsAppend
	if len(appended) != 4 {
		t.Errorf("unexpected number of appended args; got %d expected 3", len(appended))
	}
	if appended[0] != "-s" || appended[1] != "-f" {
		t.Error("expected to have set -s -f flags")
	}
	if appended[2] != "a" && appended[3] != "b" {
		t.Error("unexpected arguments on services list")
	}
}

func TestNewStopPurgeCommand(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)

	cmd.SetArgs([]string{"--purge"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command with args; error: %v", err)
	}

	if !f.down.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolStop.down Command")
	}

	argsAppend := f.down.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[1] != "--volumes" || argsAppend[0] != "--remove-orphans" {
		t.Errorf("bad arguments to KoolStop.down Command when passing --purge flag")
	}
}

func TestNewStopPurgeCommandWithServices(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)

	cmd.SetArgs([]string{"--purge", "a", "b"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command with args; error: %v", err)
	}

	if !f.rm.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolStop.rm Command")
	}

	appended := f.rm.(*builder.FakeCommand).ArgsAppend
	if len(appended) != 5 || appended[0] != "-s" || appended[1] != "-f" || appended[2] != "-v" {
		t.Errorf("bad arguments to KoolStop.rm Command when passing --purge flag")
	}
}

func TestStopMainStopsWorkspacesFirst(t *testing.T) {
	f := newFakeKoolStop()
	recorder := &workspaceStopShell{}
	f.shell = recorder
	f.env.Set("KOOL_NAME", "example")
	f.getProjects = builder.NewCommand("docker", "ps", "--all")
	f.down = builder.NewCommand("docker", "compose", "down")

	if err := f.Execute(nil); err != nil {
		t.Fatal(err)
	}
	expected := []string{
		"docker compose --project-name example-workspace-task-a down --remove-orphans",
		"docker compose --project-name example-workspace-task-b down --remove-orphans",
		"docker compose down --remove-orphans",
	}
	if strings.Join(recorder.commands, "\n") != strings.Join(expected, "\n") {
		t.Errorf("expected workspace projects to stop before main:\n%v\ngot:\n%v", expected, recorder.commands)
	}
}

func TestNewFailingDependenciesCheckStopCommand(t *testing.T) {
	f := newFakeKoolStop()

	f.check.(*checker.FakeChecker).MockError = errors.New("check error")
	cmd := NewStopCommand(f)

	assertExecGotError(t, cmd, "check error")
}
