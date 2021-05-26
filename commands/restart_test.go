package commands

import (
	"errors"
	"testing"
)

// NewRestartCommand initializes new kool start command
func TestRestartCommand(t *testing.T) {
	fakeStop := &FakeKoolService{}
	fakeStart := &FakeKoolService{}

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
	fakeStop := &FakeKoolService{}
	fakeStart := &FakeKoolService{}

	fakeStart.MockExecError = errors.New("start error")

	cmd := NewRestartCommand(fakeStop, fakeStart)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fakeStart.CalledExit {
		t.Error("did not exit command due to error on start service")
	}

	if !fakeStart.CalledError {
		t.Error("did not call Error due to error on start service")
	}
}

func TestFailingStopRestartCommand(t *testing.T) {
	fakeStop := &FakeKoolService{}
	fakeStart := &FakeKoolService{}

	fakeStop.MockExecError = errors.New("stop error")

	cmd := NewRestartCommand(fakeStop, fakeStart)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fakeStop.CalledExit {
		t.Error("did not exit command due to error on start service")
	}

	if !fakeStop.CalledError {
		t.Error("did not call Error due to error on start service")
	}
}

func TestPurgeRestartCommand(t *testing.T) {
	fakeStop := newFakeKoolStop()
	fakeStart := &FakeKoolService{}

	cmd := NewRestartCommand(fakeStop, fakeStart)
	cmd.SetArgs([]string{"--purge"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fakeStop.Flags.Purge {
		t.Error("did not set the purge flag to true in the stop service")
	}
}
