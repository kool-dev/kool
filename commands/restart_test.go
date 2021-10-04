package commands

import (
	"errors"
	"io"
	"kool-dev/kool/core/shell"
	"testing"
)

// NewRestartCommand initializes new kool start command
func TestRestartCommand(t *testing.T) {
	fakeStop := newFakeKoolService()
	fakeStart := newFakeKoolService()

	cmd := NewRestartCommand(fakeStop, fakeStart)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fakeStop.CalledExecute {
		t.Errorf("restart command did not call Execute on stop service")
	}

	if !fakeStart.CalledExecute {
		t.Errorf("restart command did not call Execute on start service")
	}
}

func TestFailingStartRestartCommand(t *testing.T) {
	fakeStop := newFakeKoolService()
	fakeStart := newFakeKoolService()

	fakeStart.MockExecuteErr = errors.New("start error")

	cmd := NewRestartCommand(fakeStop, fakeStart)

	assertExecGotError(t, cmd, "start error")
}

func TestFailingStopRestartCommand(t *testing.T) {
	fakeStop := newFakeKoolService()
	fakeStart := newFakeKoolService()

	fakeStop.MockExecuteErr = errors.New("stop error")

	cmd := NewRestartCommand(fakeStop, fakeStart)

	assertExecGotError(t, cmd, "stop error")
}

func TestPurgeRestartCommand(t *testing.T) {
	fakeStop := newFakeKoolStop()
	fakeStart := newFakeKoolService()

	cmd := NewRestartCommand(fakeStop, fakeStart)
	cmd.SetArgs([]string{"--purge"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fakeStop.Flags.Purge {
		t.Error("did not set the purge flag to true in the stop service")
	}
}

func TestRebuildRestartCommand(t *testing.T) {
	fakeStop := newFakeKoolService()
	fakeStart := newFakeKoolStart()

	cmd := NewRestartCommand(fakeStop, fakeStart)
	cmd.SetArgs([]string{"--rebuild"})

	fakeStart.rebuilder.(*KoolRebuild).shell.(*shell.FakeShell).MockOutStream = io.Discard

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fakeStart.Flags.Rebuild {
		t.Error("did not set the rebuild flag to true in the start service")
	}
}
