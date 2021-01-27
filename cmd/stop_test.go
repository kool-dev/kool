package cmd

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/compose"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolStop() *KoolStop {
	return &KoolStop{
		*newFakeKoolService(),
		&KoolStopFlags{false},
		&checker.FakeChecker{},
		&builder.FakeCommand{},
	}
}

func TestNewKoolStop(t *testing.T) {
	k := NewKoolStop()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolStop instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolStop instance")
	}

	if _, ok := k.DefaultKoolService.term.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolStop instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolStop instance")
	} else if k.Flags.Purge {
		t.Errorf("bad default value for Purge flag on default KoolStop instance")
	}

	if _, ok := k.check.(*checker.DefaultChecker); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStop instance")
	}

	if _, ok := k.doStop.(*compose.DockerCompose); !ok {
		t.Errorf("unexpected compose.DockerCompose on default KoolStop instance")
	}

	if k.doStop.(*compose.DockerCompose).Command.String() != "down" {
		t.Errorf("unexpected compose.DockerCompose.Command.String() on default KoolStop instance doStop")
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

	if f.doStop.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not expect to call AppendArgs on KoolStop.doStop Command")
	}
}

func TestNewStopPurgeCommand(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)

	cmd.SetArgs([]string{"--purge"})
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command with args; error: %v", err)
	}

	if !f.doStop.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolStop.doStop Command")
	}

	argsAppend := f.doStop.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--volumes" || argsAppend[1] != "--remove-orphans" {
		t.Errorf("bad arguments to KoolStop.doStop Command when passing --purge flag")
	}
}

func TestNewFailingDependenciesCheckStopCommand(t *testing.T) {
	f := newFakeKoolStop()

	f.check.(*checker.FakeChecker).MockError = errors.New("check error")
	cmd := NewStopCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command with args; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not exit command due to dependencies checking error")
	}

	if err := f.shell.(*shell.FakeShell).Err; err.Error() != "check error" {
		t.Errorf("expecting error 'check error', got '%v'", err)
	}
}
