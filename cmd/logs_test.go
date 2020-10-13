package cmd

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolLogs() *KoolLogs {
	return &KoolLogs{
		*newFakeKoolService(),
		&KoolLogsFlags{25, false},
		&builder.FakeCommand{},
	}
}

func newFakeFailedKoolLogs() *KoolLogs {
	return &KoolLogs{
		*newFakeKoolService(),
		&KoolLogsFlags{25, false},
		&builder.FakeCommand{MockError: errors.New("error logs")},
	}
}

func TestNewKoolLogs(t *testing.T) {
	k := NewKoolLogs()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolLogs instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolLogs instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Errorf("unexpected shell.InputReader on default KoolLogs instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolLogs instance")
	} else {
		if k.Flags.Tail != 25 {
			t.Errorf("bad default value for Tail flag on default KoolLogs instance")
		}

		if k.Flags.Follow {
			t.Errorf("bad default value for Follow flag on default KoolLogs instance")
		}
	}

	if _, ok := k.logs.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolLogs instance")
	}

	if k.logs.(*builder.DefaultCommand).String() != "docker-compose logs" {
		t.Errorf("unexpected builder.DefaultCommand.String() on default KoolLogs instance logs")
	}
}

func TestNewLogsCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Errorf("did not call SetWriter")
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--tail" || argsAppend[1] != "25" {
		t.Errorf("bad arguments to KoolLogs.logs Command with default flags")
	}
}

func TestNewLogsTailCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"--tail=10"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--tail" || argsAppend[1] != "10" {
		t.Errorf("bad arguments to KoolLogs.logs Command when passing --tail flag")
	}
}

func TestNewLogsTailAllCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"--tail=0"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--tail" || argsAppend[1] != "all" {
		t.Errorf("bad arguments to KoolLogs.logs Command when passing 0 to --tail flag")
	}
}

func TestNewLogsFollowCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"--follow"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 3 || argsAppend[2] != "--follow" {
		t.Errorf("bad arguments to KoolLogs.logs Command when passing --follow flag")
	}
}

func TestNewLogsServiceCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"app"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	args := f.logs.(*builder.FakeCommand).ArgsInteractive
	if len(args) != 1 || args[0] != "app" {
		t.Errorf("bad arguments to KoolLogs.logs Command when executing it")
	}
}

func TestFailingNewLogsCommand(t *testing.T) {
	f := newFakeFailedKoolLogs()
	cmd := NewLogsCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("expecting command to exit due to an error.")
	}

	if err := f.out.(*shell.FakeOutputWriter).Err; err.Error() != "error logs" {
		t.Errorf("expecting error 'error logs', got '%s'", err.Error())
	}
}
