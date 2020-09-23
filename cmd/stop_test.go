package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolStop() *KoolStop {
	return &KoolStop{
		*newFakeKoolService(),
		&KoolStopFlags{false},
		&FakeStartDependenciesChecker{},
		&builder.FakeCommand{},
	}
}

func TestNewKoolStop(t *testing.T) {
	k := NewKoolStop()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolStop instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolStop instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initliased on default KoolStop instance")
	} else if k.Flags.Purge {
		t.Errorf("bad default value for Purge flag on default KoolStop instance")
	}

	if _, ok := k.check.(*checker.DefaultChecker); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStop instance")
	}

	if _, ok := k.doStop.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStop instance")
	}

	if k.doStop.(*builder.DefaultCommand).String() != "docker-compose down" {
		t.Errorf("unexpected builder.DefaultCommand.String() on default KoolStop instance doStop")
	}
}

func TestNewStopCommand(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing stop command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Errorf("did not call SetWriter")
	}

	if !f.check.(*FakeStartDependenciesChecker).called {
		t.Errorf("did not call Check")
	}

	if f.doStop.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not expect to call AppendArgs on KoolStop.doStop Command")
	}
}

func TestNewStopPurgeCommand(t *testing.T) {
	f := newFakeKoolStop()
	cmd := NewStopCommand(f)

	f.Flags.Purge = true
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
